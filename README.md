# vcpgmon

## Proof of Concept/Demo
Demonstrates inferring network latency from client to PostgreSQL server written in golang.  Utilizes *pcap*
but *pfring* is preferable due to its peroformance.  *pfring* is not available on all systems.

## Flow

###Given:
1. Client sends query to Server
2. Server Returns back Row description, Data Row, Command completion, and
   Ready for query message
3. Then we have this packet back from the Client
```
We want the time delta between 1 and 2 for measuring:
   Ts = Time to send packet to server
   Te = Time for server to execute Query
   Tr1 =Time for server to send first packet back with results

Between 2 and 3 we have:
   Tc1 = Time for server to send first packet to client
   
```


   
