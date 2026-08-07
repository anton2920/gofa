package freebsd

/* From <sys/socket.h>. */
type Sockaddr struct {
	Len    uint8
	Family uint8
	Data   [14]byte
}

const SOCK_MAXADDRLEN = 255 /* longest possible addresses */

/*
 * Types
 */
const (
	SOCK_STREAM    = 1 /* stream socket */
	SOCK_DGRAM     = 2 /* datagram socket */
	SOCK_RAW       = 3 /* raw-protocol interface */
	SOCK_RDM       = 4 /* reliably-delivered message */
	SOCK_SEQPACKET = 5 /* sequenced packet stream */
)

/*
 * Option flags per-socket.
 */
const (
	SO_DEBUG        = 0x00000001 /* turn on debugging info recording */
	SO_ACCEPTCONN   = 0x00000002 /* socket has had listen() */
	SO_REUSEADDR    = 0x00000004 /* allow local address reuse */
	SO_KEEPALIVE    = 0x00000008 /* keep connections alive */
	SO_DONTROUTE    = 0x00000010 /* just use interface addresses */
	SO_BROADCAST    = 0x00000020 /* permit sending of broadcast msgs */
	SO_USELOOPBACK  = 0x00000040 /* bypass hardware when possible */
	SO_LINGER       = 0x00000080 /* linger on close if data present */
	SO_OOBINLINE    = 0x00000100 /* leave received OOB data in line */
	SO_REUSEPORT    = 0x00000200 /* allow local address & port reuse */
	SO_TIMESTAMP    = 0x00000400 /* timestamp received dgram traffic */
	SO_NOSIGPIPE    = 0x00000800 /* no SIGPIPE from EPIPE */
	SO_ACCEPTFILTER = 0x00001000 /* there is an accept filter */
	SO_BINTIME      = 0x00002000 /* timestamp received dgram traffic */
	SO_NO_OFFLOAD   = 0x00004000 /* socket cannot be offloaded */
	SO_NO_DDP       = 0x00008000 /* disable direct data placement */
	SO_REUSEPORT_LB = 0x00010000 /* reuse with load balancing */
	SO_RERROR       = 0x00020000 /* keep track of receive errors */
)

/*
 * Additional options, not kept in so_options.
 */
const (
	SO_SNDBUF          = 0x1001      /* send buffer size */
	SO_RCVBUF          = 0x1002      /* receive buffer size */
	SO_SNDLOWAT        = 0x1003      /* send low-water mark */
	SO_RCVLOWAT        = 0x1004      /* receive low-water mark */
	SO_SNDTIMEO        = 0x1005      /* send timeout */
	SO_RCVTIMEO        = 0x1006      /* receive timeout */
	SO_ERROR           = 0x1007      /* get error status and clear */
	SO_TYPE            = 0x1008      /* get socket type */
	SO_LABEL           = 0x1009      /* socket's MAC label */
	SO_PEERLABEL       = 0x1010      /* socket's peer's MAC label */
	SO_LISTENQLIMIT    = 0x1011      /* socket's backlog limit */
	SO_LISTENQLEN      = 0x1012      /* socket's complete queue length */
	SO_LISTENINCQLEN   = 0x1013      /* socket's incomplete queue length */
	SO_FIB             = 0x1014      /* get or set socket FIB */
	SO_SETFIB          = SO_FIB      /* backward compat alias */
	SO_USER_COOKIE     = 0x1015      /* user cookie (dummynet etc.) */
	SO_PROTOCOL        = 0x1016      /* get socket protocol (Linux name) */
	SO_PROTOTYPE       = SO_PROTOCOL /* alias for SO_PROTOCOL (SunOS name) */
	SO_TS_CLOCK        = 0x1017      /* clock type used for SO_TIMESTAMP */
	SO_MAX_PACING_RATE = 0x1018      /* socket's max TX pacing rate (Linux name) */
	SO_DOMAIN          = 0x1019      /* get socket domain */
	SO_SPLICE          = 0x1023      /* splice data to other socket */
)

/*
 * Level number for (get/set)sockopt() to apply to socket itself.
 */
const SOL_SOCKET = 0xffff /* options for socket level */

/*
 * Address families.
 */
const (
	AF_UNSPEC          = 0       /* unspecified */
	AF_LOCAL           = AF_UNIX /* local to host (pipes, portals) */
	AF_UNIX            = 1       /* standardized name for AF_LOCAL */
	AF_INET            = 2       /* internetwork: UDP, TCP, etc. */
	AF_IMPLINK         = 3       /* arpanet imp addresses */
	AF_PUP             = 4       /* pup protocols: e.g. BSP */
	AF_CHAOS           = 5       /* mit CHAOS protocols */
	AF_NETBIOS         = 6       /* SMB protocols */
	AF_ISO             = 7       /* ISO protocols */
	AF_OSI             = AF_ISO
	AF_ECMA            = 8       /* European computer manufacturers */
	AF_DATAKIT         = 9       /* datakit protocols */
	AF_CCITT           = 10      /* CCITT protocols, X.25 etc */
	AF_SNA             = 11      /* IBM SNA */
	AF_DECnet          = 12      /* DECnet */
	AF_DLI             = 13      /* DEC Direct data link interface */
	AF_LAT             = 14      /* LAT */
	AF_HYLINK          = 15      /* NSC Hyperchannel */
	AF_APPLETALK       = 16      /* Apple Talk */
	AF_ROUTE           = 17      /* Internal Routing Protocol */
	AF_LINK            = 18      /* Link layer interface */
	PSEUDO_AF_XTP      = 19      /* eXpress Transfer Protocol (no AF) */
	AF_COIP            = 20      /* connection-oriented IP, aka ST II */
	AF_CNT             = 21      /* Computer Network Technology */
	PSEUDO_AF_RTIP     = 22      /* Help Identify RTIP packets */
	AF_IPX             = 23      /* Novell Internet Protocol */
	AF_SIP             = 24      /* Simple Internet Protocol */
	PSEUDO_AF_PIP      = 25      /* Help Identify PIP packets */
	AF_ISDN            = 26      /* Integrated Services Digital Network*/
	AF_E164            = AF_ISDN /* CCITT E.164 recommendation */
	PSEUDO_AF_KEY      = 27      /* Internal key-management function */
	AF_INET6           = 28      /* IPv6 */
	AF_NATM            = 29      /* native ATM access */
	AF_ATM             = 30      /* ATM */
	PSEUDO_AF_HDRCMPLT = 31      /* Used by BPF to not rewrite headers in interface output routine */
	AF_NETGRAPH        = 32      /* Netgraph sockets */
	AF_SLOW            = 33      /* 802.3ad slow protocol */
	AF_SCLUSTER        = 34      /* Sitara cluster protocol */
	AF_ARP             = 35
	AF_BLUETOOTH       = 36 /* Bluetooth sockets */
	AF_IEEE80211       = 37 /* IEEE 802.11 protocol */
	AF_NETLINK         = 38 /* Netlink protocol */
	AF_INET_SDP        = 40 /* OFED Socket Direct Protocol ipv4 */
	AF_INET6_SDP       = 42 /* OFED Socket Direct Protocol ipv6 */
	AF_HYPERV          = 43 /* HyperV sockets */
	AF_DIVERT          = 44 /* divert(4) */
	AF_MAX             = 44

	AF_VENDOR00 = 39
	AF_VENDOR01 = 41
	AF_VENDOR03 = 45
	AF_VENDOR04 = 47
	AF_VENDOR05 = 49
	AF_VENDOR06 = 51
	AF_VENDOR07 = 53
	AF_VENDOR08 = 55
	AF_VENDOR09 = 57
	AF_VENDOR10 = 59
	AF_VENDOR11 = 61
	AF_VENDOR12 = 63
	AF_VENDOR13 = 65
	AF_VENDOR14 = 67
	AF_VENDOR15 = 69
	AF_VENDOR16 = 71
	AF_VENDOR17 = 73
	AF_VENDOR18 = 75
	AF_VENDOR19 = 77
	AF_VENDOR20 = 79
	AF_VENDOR21 = 81
	AF_VENDOR22 = 83
	AF_VENDOR23 = 85
	AF_VENDOR24 = 87
	AF_VENDOR25 = 89
	AF_VENDOR26 = 91
	AF_VENDOR27 = 93
	AF_VENDOR28 = 95
	AF_VENDOR29 = 97
	AF_VENDOR30 = 99
	AF_VENDOR31 = 101
	AF_VENDOR32 = 103
	AF_VENDOR33 = 105
	AF_VENDOR34 = 107
	AF_VENDOR35 = 109
	AF_VENDOR36 = 111
	AF_VENDOR37 = 113
	AF_VENDOR38 = 115
	AF_VENDOR39 = 117
	AF_VENDOR40 = 119
	AF_VENDOR41 = 121
	AF_VENDOR42 = 123
	AF_VENDOR43 = 125
	AF_VENDOR44 = 127
	AF_VENDOR45 = 129
	AF_VENDOR46 = 131
	AF_VENDOR47 = 133
)

/*
 * Protocol families, same as address families for now.
 */
const (
	PF_UNSPEC    = AF_UNSPEC
	PF_LOCAL     = AF_LOCAL
	PF_UNIX      = PF_LOCAL /* backward compatibility */
	PF_INET      = AF_INET
	PF_IMPLINK   = AF_IMPLINK
	PF_PUP       = AF_PUP
	PF_CHAOS     = AF_CHAOS
	PF_NETBIOS   = AF_NETBIOS
	PF_ISO       = AF_ISO
	PF_OSI       = AF_ISO
	PF_ECMA      = AF_ECMA
	PF_DATAKIT   = AF_DATAKIT
	PF_CCITT     = AF_CCITT
	PF_SNA       = AF_SNA
	PF_DECnet    = AF_DECnet
	PF_DLI       = AF_DLI
	PF_LAT       = AF_LAT
	PF_HYLINK    = AF_HYLINK
	PF_APPLETALK = AF_APPLETALK
	PF_ROUTE     = AF_ROUTE
	PF_LINK      = AF_LINK
	PF_XTP       = PSEUDO_AF_XTP /* really just proto family, no AF */
	PF_COIP      = AF_COIP
	PF_CNT       = AF_CNT
	PF_SIP       = AF_SIP
	PF_IPX       = AF_IPX
	PF_RTIP      = PSEUDO_AF_RTIP /* same format as AF_INET */
	PF_PIP       = PSEUDO_AF_PIP
	PF_ISDN      = AF_ISDN
	PF_KEY       = PSEUDO_AF_KEY
	PF_INET6     = AF_INET6
	PF_NATM      = AF_NATM
	PF_ATM       = AF_ATM
	PF_NETGRAPH  = AF_NETGRAPH
	PF_SLOW      = AF_SLOW
	PF_SCLUSTER  = AF_SCLUSTER
	PF_ARP       = AF_ARP
	PF_BLUETOOTH = AF_BLUETOOTH
	PF_IEEE80211 = AF_IEEE80211
	PF_NETLINK   = AF_NETLINK
	PF_INET_SDP  = AF_INET_SDP
	PF_INET6_SDP = AF_INET6_SDP
	PF_DIVERT    = AF_DIVERT
	PF_MAX       = AF_MAX
)

/*
 * Maximum queue length specifiable by listen.
 */
const SOMAXCONN = 128

/*
 * howto arguments for shutdown(2), specified by Posix.1g.
 */
const (
	SHUT_RD   = 0 /* shut down the reading side */
	SHUT_WR   = 1 /* shut down the writing side */
	SHUT_RDWR = 2 /* shut down both sides */
)
