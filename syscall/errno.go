package syscall

type Errno int

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

var errors = [...]string{
	"",
	EPERM:   "Operation not permitted",
	ENOENT:  "No such file or directory",
	ESRCH:   "No such process",
	EINTR:   "Interrupted system call",
	EIO:     "Input/output error",
	ENXIO:   "Device not configured",
	E2BIG:   "Argument list too long",
	ENOEXEC: "Exec format error",
	EBADF:   "Bad file descriptor",
	ECHILD:  "No child processes",
	EDEADLK: "Resource deadlock avoided",
	/* 11 was EAGAIN */
	ENOMEM:  "Cannot allocate memory",
	EACCES:  "Permission denied",
	EFAULT:  "Bad address",
	ENOTBLK: "Block device required",
	EBUSY:   "Device busy",
	EEXIST:  "File exists",
	EXDEV:   "Cross-device link",
	ENODEV:  "Operation not supported by device",
	ENOTDIR: "Not a directory",
	EISDIR:  "Is a directory",
	EINVAL:  "Invalid argument",
	ENFILE:  "Too many open files in system",
	EMFILE:  "Too many open files",
	ENOTTY:  "Inappropriate ioctl for device",
	ETXTBSY: "Text file busy",
	EFBIG:   "File too large",
	ENOSPC:  "No space left on device",
	ESPIPE:  "Illegal seek",
	EROFS:   "Read-only filesystem",
	EMLINK:  "Too many links",
	EPIPE:   "Broken pipe",

	/* math software */
	EDOM:   "Numerical argument out of domain",
	ERANGE: "Result too large",

	/* non-blocking and interrupt i/o */
	EAGAIN:      "Resource temporarily unavailable",
	EINPROGRESS: "Operation now in progress",
	EALREADY:    "Operation already in progress",

	/* ipc/network software -- argument errors */
	ENOTSOCK:        "Socket operation on non-socket",
	EDESTADDRREQ:    "Destination address required",
	EMSGSIZE:        "Message too long",
	EPROTOTYPE:      "Protocol wrong type for socket",
	ENOPROTOOPT:     "Protocol not available",
	EPROTONOSUPPORT: "Protocol not supported",
	ESOCKTNOSUPPORT: "Socket type not supported",
	EOPNOTSUPP:      "Operation not supported",
	EPFNOSUPPORT:    "Protocol family not supported",
	EAFNOSUPPORT:    "Address family not supported by protocol family",
	EADDRINUSE:      "Address already in use",
	EADDRNOTAVAIL:   "Can't assign requested address",

	/* ipc/network software -- operational errors */
	ENETDOWN:     "Network is down",
	ENETUNREACH:  "Network is unreachable",
	ENETRESET:    "Network dropped connection on reset",
	ECONNABORTED: "Software caused connection abort",
	ECONNRESET:   "Connection reset by peer",
	ENOBUFS:      "No buffer space available",
	EISCONN:      "Socket is already connected",
	ENOTCONN:     "Socket is not connected",
	ESHUTDOWN:    "Can't send after socket shutdown",
	ETOOMANYREFS: "Too many references: can't splice",
	ETIMEDOUT:    "Operation timed out",
	ECONNREFUSED: "Connection refused",

	ELOOP:        "Too many levels of symbolic links",
	ENAMETOOLONG: "File name too long",

	/* should be rearranged */
	EHOSTDOWN:    "Host is down",
	EHOSTUNREACH: "No route to host",
	ENOTEMPTY:    "Directory not empty",

	/* quotas & mush */
	EPROCLIM: "Too many processes",
	EUSERS:   "Too many users",
	EDQUOT:   "Disc quota exceeded",

	/* Network File System */
	ESTALE:        "Stale NFS file handle",
	EREMOTE:       "Too many levels of remote in path",
	EBADRPC:       "RPC struct is bad",
	ERPCMISMATCH:  "RPC version wrong",
	EPROGUNAVAIL:  "RPC prog. not avail",
	EPROGMISMATCH: "Program version wrong",
	EPROCUNAVAIL:  "Bad procedure for program",

	ENOLCK: "No locks available",
	ENOSYS: "Function not implemented",

	EFTYPE:    "Inappropriate file type or format",
	EAUTH:     "Authentication error",
	ENEEDAUTH: "Need authenticator",
	EIDRM:     "Identifier removed",
	ENOMSG:    "No message of desired type",
	EOVERFLOW: "Value too large to be stored in data type",
	ECANCELED: "Operation canceled",
	EILSEQ:    "Illegal byte sequence",
	ENOATTR:   "Attribute not found",

	EDOOFUS: "Programming error",

	EBADMSG:   "Bad message",
	EMULTIHOP: "Multihop attempted",
	ENOLINK:   "Link has been severed",
	EPROTO:    "Protocol error",

	ENOTCAPABLE:     "Capabilities insufficient",
	ECAPMODE:        "Not permitted in capability mode",
	ENOTRECOVERABLE: "State not recoverable",
	EOWNERDEAD:      "Previous owner died",
	EINTEGRITY:      "Integrity check failed",
}

func (errno Errno) String() string {
	if (errno >= 0) && (errno <= ELAST) {
		return errors[errno]
	}
	return "<UNKNOWN ERROR>"
}
