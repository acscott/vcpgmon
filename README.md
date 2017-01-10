# vcpgmon

## Proof of Concept/Demo
Demonstrates inferring network latency from client to PostgreSQL server written in golang.  Utilizes *pcap*
but *pfring* is preferable due to its peroformance.  *pfring* is not available on all systems.


###Given:
1. PostgreSQL Client sends query to Server
2. Server Returns back Row description, Data Row, Command completion, and
   Ready for query message
3. Then we have this example ACK packet back from the Client

```  2017-01-09 12:44:47.349851 b4:6d:83:8e:e2:16 (oui Unknown) > 08:00:27:91:9e:51 (oui Unknown), ethertype IPv4 (0x0800), length 66: (tos 0x0, ttl 64, id 13372, offset 0, flags [DF], proto TCP (6), length 52)
    192.168.1.106.47878 > lab1.postgres: Flags [.], cksum 0xe785 (correct), ack 463, win 237, options [nop,nop,TS val 81030353 ecr 523807427], length 0
        0x0000:  0800 2791 9e51 b46d 838e e216 0800 4500  ..'..Q.m......E.
        0x0010:  0034 343c 4000 4006 8260 c0a8 016a c0a8  .44<@.@..`...j..
        0x0020:  016d bb06 1538 4318 d27c 8aac 5e01 8010  .m...8C..|..^...
        0x0030:  00ed e785 0000 0101 080a 04d4 6cd1 1f38  ............l..8
        0x0040:  aac3                                     .
```

The time delta between 1 and 2 can be used for measuring:
   Ts = Time to send packet to server
   Te = Time for server to execute Query
   Tr1 =Time for server to send first packet back with results

Between 2 and 3 we have:
   Tc1 = Time for server to send a firtst packet to client
   Ta = Time for client to send back an AC
   
   


   
