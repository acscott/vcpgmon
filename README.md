# vcpgmon

## Proof of Concept/Demo
Demonstrates inferring network latency from client to PostgreSQL server written in golang.  Utilizes *pcap*
but *pfring* is preferable due to its peroformance.  *pfring* is not available on all systems.


###Given:
1. PostgreSQL Client sends query to Server
2. Server Returns back Row description, Data Row, Command completion, and
   Ready for query message
3. Then we receive a tcp ACK packet back from the Client                                .
```
The time delta between 1 and 2 can be used for measuring:
   Ts = Time to send packet to server
   Te = Time for server to execute Query
   Tr1 =Time for server to send first packet back with results

Between 2 and 3 we have:
   Tc1 = Time for server to send 1 or packets to client
   Ta = Time for client to send back an AC
```
If we use a query with a small result set say
`
    Select 1;
`
Then it can be a probing query to discover network latency.  Also, you can use this tool to **compare network latency between other systems.**

It might be noted that **Tc1** is included in the latency measurement.  This is why you should use the same query when comparing across samples.

##Example


