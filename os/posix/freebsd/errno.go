package freebsd

import "github.com/anton2920/gofa/context"

type Errno uintptr

const (
	/* From <errno.h>. */
	EPERM   = Errno(1)  /* Operation not permitted */
	ENOENT  = Errno(2)  /* No such file or directory */
	ESRCH   = Errno(3)  /* No such process */
	EINTR   = Errno(4)  /* Interrupted system call */
	EIO     = Errno(5)  /* Input/output error */
	ENXIO   = Errno(6)  /* Device not configured */
	E2BIG   = Errno(7)  /* Argument list too long */
	ENOEXEC = Errno(8)  /* Exec format error */
	EBADF   = Errno(9)  /* Bad file descriptor */
	ECHILD  = Errno(10) /* No child processes */
	EDEADLK = Errno(11) /* Resource deadlock avoided */
	/* 11 was EAGAIN */
	ENOMEM  = Errno(12) /* Cannot allocate memory */
	EACCES  = Errno(13) /* Permission denied */
	EFAULT  = Errno(14) /* Bad address */
	ENOTBLK = Errno(15) /* Block device required */
	EBUSY   = Errno(16) /* Device busy */
	EEXIST  = Errno(17) /* File exists */
	EXDEV   = Errno(18) /* Cross-device link */
	ENODEV  = Errno(19) /* Operation not supported by device */
	ENOTDIR = Errno(20) /* Not a directory */
	EISDIR  = Errno(21) /* Is a directory */
	EINVAL  = Errno(22) /* Invalid argument */
	ENFILE  = Errno(23) /* Too many open files in system */
	EMFILE  = Errno(24) /* Too many open files */
	ENOTTY  = Errno(25) /* Inappropriate ioctl for device */
	ETXTBSY = Errno(26) /* Text file busy */
	EFBIG   = Errno(27) /* File too large */
	ENOSPC  = Errno(28) /* No space left on device */
	ESPIPE  = Errno(29) /* Illegal seek */
	EROFS   = Errno(30) /* Read-only filesystem */
	EMLINK  = Errno(31) /* Too many links */
	EPIPE   = Errno(32) /* Broken pipe */

	/* math software */
	EDOM   = Errno(33) /* Numerical argument out of domain */
	ERANGE = Errno(34) /* Result too large */

	/* non-blocking and interrupt i/o */
	EAGAIN      = Errno(35) /* Resource temporarily unavailable */
	EINPROGRESS = Errno(36) /* Operation now in progress */
	EALREADY    = Errno(37) /* Operation already in progress */

	/* ipc/network software -- argument errors */
	ENOTSOCK        = Errno(38) /* Socket operation on non-socket */
	EDESTADDRREQ    = Errno(39) /* Destination address required */
	EMSGSIZE        = Errno(40) /* Message too long */
	EPROTOTYPE      = Errno(41) /* Protocol wrong type for socket */
	ENOPROTOOPT     = Errno(42) /* Protocol not available */
	EPROTONOSUPPORT = Errno(43) /* Protocol not supported */
	ESOCKTNOSUPPORT = Errno(44) /* Socket type not supported */
	EOPNOTSUPP      = Errno(45) /* Operation not supported */
	EPFNOSUPPORT    = Errno(46) /* Protocol family not supported */
	EAFNOSUPPORT    = Errno(47) /* Address family not supported by protocol family */
	EADDRINUSE      = Errno(48) /* Address already in use */
	EADDRNOTAVAIL   = Errno(49) /* Can't assign requested address */

	/* ipc/network software -- operational errors */
	ENETDOWN     = Errno(50) /* Network is down */
	ENETUNREACH  = Errno(51) /* Network is unreachable */
	ENETRESET    = Errno(52) /* Network dropped connection on reset */
	ECONNABORTED = Errno(53) /* Software caused connection abort */
	ECONNRESET   = Errno(54) /* Connection reset by peer */
	ENOBUFS      = Errno(55) /* No buffer space available */
	EISCONN      = Errno(56) /* Socket is already connected */
	ENOTCONN     = Errno(57) /* Socket is not connected */
	ESHUTDOWN    = Errno(58) /* Can't send after socket shutdown */
	ETOOMANYREFS = Errno(59) /* Too many references: can't splice */
	ETIMEDOUT    = Errno(60) /* Operation timed out */
	ECONNREFUSED = Errno(61) /* Connection refused */

	ELOOP        = Errno(62) /* Too many levels of symbolic links */
	ENAMETOOLONG = Errno(63) /* File name too long */

	/* should be rearranged */
	EHOSTDOWN    = Errno(64) /* Host is down */
	EHOSTUNREACH = Errno(65) /* No route to host */
	ENOTEMPTY    = Errno(66) /* Directory not empty */

	/* quotas & mush */
	EPROCLIM = Errno(67) /* Too many processes */
	EUSERS   = Errno(68) /* Too many users */
	EDQUOT   = Errno(69) /* Disc quota exceeded */

	/* Network File System */
	ESTALE        = Errno(70) /* Stale NFS file handle */
	EREMOTE       = Errno(71) /* Too many levels of remote in path */
	EBADRPC       = Errno(72) /* RPC struct is bad */
	ERPCMISMATCH  = Errno(73) /* RPC version wrong */
	EPROGUNAVAIL  = Errno(74) /* RPC prog. not avail */
	EPROGMISMATCH = Errno(75) /* Program version wrong */
	EPROCUNAVAIL  = Errno(76) /* Bad procedure for program */

	ENOLCK = Errno(77) /* No locks available */
	ENOSYS = Errno(78) /* Function not implemented */

	EFTYPE    = Errno(79) /* Inappropriate file type or format */
	EAUTH     = Errno(80) /* Authentication error */
	ENEEDAUTH = Errno(81) /* Need authenticator */
	EIDRM     = Errno(82) /* Identifier removed */
	ENOMSG    = Errno(83) /* No message of desired type */
	EOVERFLOW = Errno(84) /* Value too large to be stored in data type */
	ECANCELED = Errno(85) /* Operation canceled */
	EILSEQ    = Errno(86) /* Illegal byte sequence */
	ENOATTR   = Errno(87) /* Attribute not found */

	EDOOFUS = Errno(88) /* Programming error */

	EBADMSG   = Errno(89) /* Bad message */
	EMULTIHOP = Errno(90) /* Multihop attempted */
	ENOLINK   = Errno(91) /* Link has been severed */
	EPROTO    = Errno(92) /* Protocol error */

	ENOTCAPABLE     = Errno(93) /* Capabilities insufficient */
	ECAPMODE        = Errno(94) /* Not permitted in capability mode */
	ENOTRECOVERABLE = Errno(95) /* State not recoverable */
	EOWNERDEAD      = Errno(96) /* Previous owner died */
	EINTEGRITY      = Errno(97) /* Integrity check failed */

	ELAST = Errno(97) /* Must be equal largest errno */
)

var strerror = [...]string{
	"",
	EPERM:   "operation not permitted",
	ENOENT:  "no such file or directory",
	ESRCH:   "no such process",
	EINTR:   "interrupted system call",
	EIO:     "input/output error",
	ENXIO:   "device not configured",
	E2BIG:   "argument list too long",
	ENOEXEC: "exec format error",
	EBADF:   "bad file descriptor",
	ECHILD:  "no child processes",
	EDEADLK: "resource deadlock avoided",
	/* 11 was EAGAIN */
	ENOMEM:  "cannot allocate memory",
	EACCES:  "permission denied",
	EFAULT:  "bad address",
	ENOTBLK: "block device required",
	EBUSY:   "device busy",
	EEXIST:  "file exists",
	EXDEV:   "cross-device link",
	ENODEV:  "operation not supported by device",
	ENOTDIR: "not a directory",
	EISDIR:  "is a directory",
	EINVAL:  "invalid argument",
	ENFILE:  "too many open files in system",
	EMFILE:  "too many open files",
	ENOTTY:  "inappropriate ioctl for device",
	ETXTBSY: "text file busy",
	EFBIG:   "file too large",
	ENOSPC:  "no space left on device",
	ESPIPE:  "illegal seek",
	EROFS:   "read-only filesystem",
	EMLINK:  "too many links",
	EPIPE:   "broken pipe",

	/* math software */
	EDOM:   "numerical argument out of domain",
	ERANGE: "result too large",

	/* non-blocking and interrupt i/o */
	EAGAIN:      "resource temporarily unavailable",
	EINPROGRESS: "operation now in progress",
	EALREADY:    "operation already in progress",

	/* ipc/network software -- argument errors */
	ENOTSOCK:        "socket operation on non-socket",
	EDESTADDRREQ:    "destination address required",
	EMSGSIZE:        "message too long",
	EPROTOTYPE:      "protocol wrong type for socket",
	ENOPROTOOPT:     "protocol not available",
	EPROTONOSUPPORT: "protocol not supported",
	ESOCKTNOSUPPORT: "socket type not supported",
	EOPNOTSUPP:      "operation not supported",
	EPFNOSUPPORT:    "protocol family not supported",
	EAFNOSUPPORT:    "address family not supported by protocol family",
	EADDRINUSE:      "address already in use",
	EADDRNOTAVAIL:   "can't assign requested address",

	/* ipc/network software -- operational errors */
	ENETDOWN:     "network is down",
	ENETUNREACH:  "network is unreachable",
	ENETRESET:    "network dropped connection on reset",
	ECONNABORTED: "software caused connection abort",
	ECONNRESET:   "connection reset by peer",
	ENOBUFS:      "no buffer space available",
	EISCONN:      "socket is already connected",
	ENOTCONN:     "socket is not connected",
	ESHUTDOWN:    "can't send after socket shutdown",
	ETOOMANYREFS: "too many references: can't splice",
	ETIMEDOUT:    "operation timed out",
	ECONNREFUSED: "connection refused",

	ELOOP:        "too many levels of symbolic links",
	ENAMETOOLONG: "file name too long",

	/* should be rearranged */
	EHOSTDOWN:    "host is down",
	EHOSTUNREACH: "no route to host",
	ENOTEMPTY:    "directory not empty",

	/* quotas & mush */
	EPROCLIM: "too many processes",
	EUSERS:   "too many users",
	EDQUOT:   "disc quota exceeded",

	/* Network File System */
	ESTALE:        "stale NFS file handle",
	EREMOTE:       "too many levels of remote in path",
	EBADRPC:       "RPC struct is bad",
	ERPCMISMATCH:  "RPC version wrong",
	EPROGUNAVAIL:  "RPC prog. not avail",
	EPROGMISMATCH: "program version wrong",
	EPROCUNAVAIL:  "bad procedure for program",

	ENOLCK: "no locks available",
	ENOSYS: "function not implemented",

	EFTYPE:    "inappropriate file type or format",
	EAUTH:     "authentication error",
	ENEEDAUTH: "need authenticator",
	EIDRM:     "identifier removed",
	ENOMSG:    "no message of desired type",
	EOVERFLOW: "value too large to be stored in data type",
	ECANCELED: "operation canceled",
	EILSEQ:    "illegal byte sequence",
	ENOATTR:   "attribute not found",

	EDOOFUS: "programming error",

	EBADMSG:   "bad message",
	EMULTIHOP: "multihop attempted",
	ENOLINK:   "link has been severed",
	EPROTO:    "protocol error",

	ENOTCAPABLE:     "capabilities insufficient",
	ECAPMODE:        "not permitted in capability mode",
	ENOTRECOVERABLE: "state not recoverable",
	EOWNERDEAD:      "previous owner died",
	EINTEGRITY:      "integrity check failed",
}

func ReportPotentialError(ctx *context.Context, errno uintptr) bool {
	if errno == 0 {
		return true
	} else {
		ctx.NewErrorWithCode(int(errno)).S(Errno(errno).String())
		return false
	}
}

func (errno Errno) String() string {
	if (errno >= 0) && (errno <= ELAST) {
		return strerror[errno]
	}
	return "<UNKNOWN ERROR>"
}
