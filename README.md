#  Golang LVS ![Build Status Image](https://travis-ci.org/nanobox-io/golang-lvs.svg)

A small wrapper around ipvsadm to support go interacting with the Linux Virtual Server.


### Data Types:

#### Ipvs
Data:
 - MulticastInterface: String with the name of the interface broadcast the multicast state information on.
 - Syncid: Id to use when broadcasting state.
 - Tcp: Timeout for TCP connections.
 - Tcpfin: Timeout for TCP-FIN packets.
 - Udp: Timeout for UDP connections.
 - Services: Slice of Services.

Methods:
 - FindService
 - AddService
 - EditService
 - RemoveService
 - SetTimeouts
 - Restore
 - Save
 - StartDaemon
 - StopDaemon
 - Zero

#### Service
Data:
 - Host: IP associated to the service.
 - Port: Port that the service listens to.
 - Type: Type of service (tcp, udp, fwmark).
 - Scheduler: Method of assigning connections to downstream servers (rr, wrr, lc, wlc, lblc, lblcr, dh, sh, sed, nq).
 - Persistence: Persistent connection timeout.
 - Netmask: Netmask to use to group connections together.
 - Servers: Slice of Servers.

Methods:
 - FindServer
 - AddServer
 - EditServer
 - RemoveServer
 - Zero
 - ToJson
 - FromJson
 - String

#### Server
Data:
 - Host: IP associated with the server.
 - Port: Port the downstream server is listening on.
 - Forwarder: Method to forward to the downstream server (g=gatewaying, i=ipip, m=masquerading).
 - Weight: Relative weight of this server to the others. 0 means no new connections.
 - UpperThreshold: Stop sending connections when this limit is reached. 0 means no limit.
 - LowerThreshold: Restart sending connections when connections drop to this number. 0 means not set.

Methods:
 - ToJson
 - FromJson
 - String