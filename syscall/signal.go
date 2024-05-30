package syscall

import (
	"fmt"
	"unsafe"
)

type Signal int

/* From <sys/_sigset.h>. */
const _SIG_WORDS = 4

type Sigset struct {
	Bits [_SIG_WORDS]uint32
}

/* From <sys/signal.h>. */
type Sigaction_t struct {
	Handler unsafe.Pointer
	Flags   int32
	Mask    Sigset
}

const (
	/* From <sys/signal.h>. */
	SIGHUP    = Signal(1)  /*  hangup  */
	SIGINT    = Signal(2)  /*  interrupt  */
	SIGQUIT   = Signal(3)  /*  quit  */
	SIGILL    = Signal(4)  /*  illegal instr. (not reset when caught)  */
	SIGTRAP   = Signal(5)  /*  trace trap (not reset when caught)  */
	SIGABRT   = Signal(6)  /*  abort()  */
	SIGIOT    = SIGABRT    /*  compatibility  */
	SIGEMT    = Signal(7)  /*  EMT instruction  */
	SIGFPE    = Signal(8)  /*  floating point exception  */
	SIGKILL   = Signal(9)  /*  kill (cannot be caught or ignored)  */
	SIGBUS    = Signal(10) /*  bus error  */
	SIGSEGV   = Signal(11) /*  segmentation violation  */
	SIGSYS    = Signal(12) /*  non-existent system call invoked  */
	SIGPIPE   = Signal(13) /*  write on a pipe with no one to read it  */
	SIGALRM   = Signal(14) /*  alarm clock  */
	SIGTERM   = Signal(15) /*  software termination signal from kill  */
	SIGURG    = Signal(16) /*  urgent condition on IO channel  */
	SIGSTOP   = Signal(17) /*  sendable stop signal not from tty  */
	SIGTSTP   = Signal(18) /*  stop signal from tty  */
	SIGCONT   = Signal(19) /*  continue a stopped process  */
	SIGCHLD   = Signal(20) /*  to parent on child stop or exit  */
	SIGTTIN   = Signal(21) /*  to readers pgrp upon background tty read  */
	SIGTTOU   = Signal(22) /*  like TTIN if (tp->t_local&LTOSTOP)  */
	SIGIO     = Signal(23) /*  input/output possible signal  */
	SIGXCPU   = Signal(24) /*  exceeded CPU time limit  */
	SIGXFSZ   = Signal(25) /*  exceeded file size limit  */
	SIGVTALRM = Signal(26) /*  virtual time alarm  */
	SIGPROF   = Signal(27) /*  profiling time alarm  */
	SIGWINCH  = Signal(28) /*  window size changes  */
	SIGINFO   = Signal(29) /*  information request  */
	SIGUSR1   = Signal(30) /*  user defined signal 1  */
	SIGUSR2   = Signal(31) /*  user defined signal 2  */
	SIGTHR    = Signal(32) /*  reserved by thread library.  */
	SIGLWP    = SIGTHR
	SIGLIBRT  = Signal(33) /*  reserved by real-time library.  */
)

const (
	SA_ONSTACK   = 0x0001 /* take signal on signal stack */
	SA_RESTART   = 0x0002 /* restart system call on signal return */
	SA_RESETHAND = 0x0004 /* reset to SIG_DFL when taking signal */
	SA_NODEFER   = 0x0010 /* don't mask the signal we're delivering */
	SA_NOCLDWAIT = 0x0020 /* don't keep zombies around */
	SA_SIGINFO   = 0x0040 /* signal handler with SA_SIGINFO args */
)

/* From <sys/signal.h>. */
var SIG_IGN = unsafe.Pointer(uintptr(1))

var signals = [...]string{
	SIGHUP:    "terminal line hangup",
	SIGINT:    "interrupt program",
	SIGQUIT:   "quit program",
	SIGILL:    "illegal instruction",
	SIGTRAP:   "trace trap",
	SIGABRT:   "abort(3) call",
	SIGEMT:    "emulate instruction executed",
	SIGFPE:    "floating-point exception",
	SIGKILL:   "kill program",
	SIGBUS:    "bus error",
	SIGSEGV:   "segmentation violation",
	SIGSYS:    "non-existent system call invoked",
	SIGPIPE:   "write on a pipe with no reader",
	SIGALRM:   "real-time timer expired",
	SIGTERM:   "software termination signal",
	SIGURG:    "urgent condition present on socket",
	SIGSTOP:   "stop",
	SIGTSTP:   "stop signal generated from keyboard",
	SIGCONT:   "continue after stop",
	SIGCHLD:   "child status has changed",
	SIGTTIN:   "background read attempted from control terminal",
	SIGTTOU:   "background write attempted to control terminal",
	SIGIO:     "I/O is possible on a descriptor",
	SIGXCPU:   "cpu time limit exceeded",
	SIGXFSZ:   "file size limit exceeded",
	SIGVTALRM: "virtual time alarm",
	SIGPROF:   "profiling time alarm",
	SIGWINCH:  "window size change",
	SIGINFO:   "status request from keyboard",
	SIGUSR1:   "user defined signal 1",
	SIGUSR2:   "user defined signal 2",
	SIGTHR:    "reserved by thread library",
	SIGLIBRT:  "reserved by real-time library",
}

func (s Signal) Signal() {}

func (s Signal) String() string {
	if (s >= 0) && (int(s) <= len(signals)) {
		return signals[s]
	}
	return "<UNKNOWN SIGNAL>"
}

func IgnoreSignals(signals ...Signal) error {
	for i := 0; i < len(signals); i++ {
		if err := Sigaction(int32(signals[i]), &Sigaction_t{Handler: SIG_IGN, Flags: SA_ONSTACK | SA_RESTART | SA_SIGINFO}, nil); err != nil {
			return fmt.Errorf("failed to set ignore handler to signal: %w", err)
		}
	}
	return nil
}
