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
	//"golang.org/x/tools/go/gcimporter15/testdata"
	"github.com/google/gopacket/layers"
)

const (
	DB_USER = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME = "tpcc"
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
	netLatency float32
}

var (
	srv vcDb
)

func main() {

	nethandle := getPacketSource()

	c := pollNet(nethandle)
	for {
		select {
		case s := <-c:
			fmt.Println(s)
		}
	}
}

func pollNet(packetSource *gopacket.PacketSource) chan Metrics {
	c := make(chan Metrics)
	m := Metrics{netLatency: 100}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("poll")
			for packet := range packetSource.Packets() {
				fmt.Println(packet.Layer(layers.LayerTypeIPv4))
			}
			m.netLatency = getLatency()
			c <- m
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	//1.  Load configuration file; fall back to getting last config file
	//from collector (using hard-coded default collector).
	//2.  Check OS, hardware, database system config and compare to prior
	// configuration. Save the configuration  including date/time.
	// As part of the state be sure to include timezone, localization,
	// ntp settings.
	//3. Connect to database.
	//4. Connect to collector.  If collector not available attempt to
	// diagnose why.  Firewall? Internet connectivity? SELinux?
	// Permission changes?  Host file, gateway changes?  etc.  If connector
	// connects record time to establish connection.
	//

	fmt.Println("stub: func init()")
	srv.connectionString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	//srv.db = new(sql.DB)
	var err error
	srv.db, err = sql.Open("postgres", srv.connectionString)
	//defer db.Close()  sql package doc indicates closing is rarely needed
	checkErr(err)

	err = srv.db.Ping()
	checkErr(err)

	fmt.Println("connected.")
}
