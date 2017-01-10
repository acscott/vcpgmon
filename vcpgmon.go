package main

// Ideally we only monitor things the version of pg has
// this way we do not have extra conditionals to process every loop
// for maximum performance.  We want our monitor to have as low impact as
// possible.
// So we prefer to use the pfring collector for tcp capture rather than pcap
// for instance.
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"os"

	//DEBUG:
	//"encoding/hex"
	"math"
)

var (
	/*	DB_USER = "postgres"
		DB_PASSWORD = "postgres"
		DB_NAME = "tpcc"*/
	IF = "enp0s3"
	PORT = "5432"
)

type ProductVersion struct {
	number string
	text   string
}
type vcDb struct {
	productName      string
	version          ProductVersion
	connectionString string
	db               *sql.DB
}

type Metrics struct {
	netLatency time.Duration
}

type Stat struct {
	min                     time.Duration //ns
	max                     time.Duration
	sum                     time.Duration
	avg                     float64
	svar                    float64
	n                       int64

	oldmM, newM, oldS, newS float64
}

func (s *Stat) calc(sample time.Duration) {
	fsample := float64(sample)
	if sample < s.min {
		s.min = sample
	}
	if sample > s.max {
		s.max = sample
	}
	s.n++
	s.sum += sample
	s.avg = float64(int64(s.sum) / s.n) / 1000 // divide by 1000 to microsecond

	//thanks to http://www.johndcook.com/blog/standard_deviation/
	if s.n == 1 {
		s.oldmM = fsample
		s.newM = fsample
		s.oldS = 0.0
		s.svar = 0.0
	} else {
		s.newM = s.oldmM + (fsample - s.oldmM) / float64(s.n)
		s.newS = s.oldS + (fsample - s.oldmM) * (fsample - s.newM)

		s.oldmM = s.newM
		s.oldS = s.newS

		s.svar = s.newS / (float64(s.n) - 1)
	}
}

/*var (
	srv vcDb
)*/

func main() {
	var stat Stat
	var microsecond = string("µs")
	nethandle := getPacketSource()

	c := pollNet(nethandle)
	for {
		select {
		case s := <-c:
			stat.calc(s.netLatency)
			fmt.Printf(" %v\tLatency: %-v\tAvg: %-.0f%v \tSDev: %-.2f%v", stat.n, s.netLatency, stat.avg,
				microsecond, math.Sqrt(stat.svar) / 1000, microsecond)
		}
	}
}

func pollNet(packetSource *gopacket.PacketSource) chan Metrics {
	var priorTS, TS time.Time
	c := make(chan Metrics)
	m := Metrics{netLatency: 100}

	//see https://godoc.org/github.com/google/gopacket on Fast Decoding with DecodingLayerParser
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("poll")
			//TODO: What do we do with layer type gopacket.Payload?  Examples also pass &payload here
			//Then DecodeLayers will not emit an err "No decoder for layer type Payload" which
			//may be okay to ignore, but throwing an error is often slightly slower than the
			//alternative in many languages.  So to speed this up we may want to pass in &payload
			//Time the difference to see which is more performant.
			parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp)
			decoded := []gopacket.LayerType{}

			//To get network latency between server and client we look for a packet sent by the server
			//followed by a packet from the client that does not include a PostgreSQL message.  This
			//last packet is a TCP ACK to the server from the client to acknowledge receipt.  It is
			//because of this ACK that we can infer the network latency.
			for packet := range packetSource.Packets() {

				//Is there a way to Decode only some of the layers and not all of the layers
				//for speed/efficiency?

				err := parser.DecodeLayers(packet.Data(), &decoded)
				//DEBUG
				//fmt.Println("================================================================")
				//fmt.Println("Timestamp for packet:",packet.Metadata().CaptureInfo.Timestamp)

				//
				if err != nil {
					//fmt.Println("Error at parser.DecodeLayers(packet.Data(), &decoded:")
					//fmt.Println(err.)
				}
				for _, layerType := range decoded {
					switch layerType {
					case layers.LayerTypeIPv6:
						fmt.Print("\n    IP6 ", ip6.SrcIP, ip6.DstIP)

					case layers.LayerTypeIPv4:
						fmt.Print("\n    IP4 ", ip4.SrcIP, ip4.DstIP)
						if (packet.ApplicationLayer() != nil) {

							//DEBUG
							//fmt.Println("Message:\n", hex.Dump(packet.ApplicationLayer().LayerContents()))

							priorTS = packet.Metadata().CaptureInfo.Timestamp
						} else if tcp.ACK {
							TS = packet.Metadata().CaptureInfo.Timestamp
							if !priorTS.IsZero() {
								//subtract prior timestamp from TS to get latency
								m.netLatency = TS.Sub(priorTS)
								c <- m

								//DEBUG:
								//fmt.Print("  Latency:",TS.Sub(priorTS))
								zeroTime(&priorTS)
								zeroTime(&TS)
							}
						}
					}
				}
			}
		}
	}()

	return c
}

func getLatency() float32 {
	return 100
}

func getPacketSource() *gopacket.PacketSource {

	if handle, err := pcap.OpenLive(IF, 65535, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp and port " + PORT); err != nil {
		// optional
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		return packetSource
	}

}

func zeroTime(p *time.Time) {
	var zt time.Time

	*p = zt
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	cnt := len(os.Args)
	if cnt > 1 {
		IF = os.Args[1]
	}
	if cnt > 2 {
		PORT = os.Args[2]
	}
	fmt.Printf("Monitoring network interface %v on port %v\n", IF, PORT)
}
