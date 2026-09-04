package freebsd

/* From <sys/unistd.h>. */
const (
	RFNAMEG     = (1 << 0)  /* UNIMPL new plan9 `name space' */
	RFENVG      = (1 << 1)  /* UNIMPL copy plan9 `env space' */
	RFFDG       = (1 << 2)  /* copy fd table */
	RFNOTEG     = (1 << 3)  /* UNIMPL create new plan9 `note group' */
	RFPROC      = (1 << 4)  /* change child (else changes curproc) */
	RFMEM       = (1 << 5)  /* share `address space' */
	RFNOWAIT    = (1 << 6)  /* give child to init */
	RFCNAMEG    = (1 << 10) /* UNIMPL zero plan9 `name space' */
	RFCENVG     = (1 << 11) /* UNIMPL zero plan9 `env space' */
	RFCFDG      = (1 << 12) /* close all fds, zero fd table */
	RFTHREAD    = (1 << 13) /* enable kernel thread support */
	RFSIGSHARE  = (1 << 14) /* share signal handlers */
	RFLINUXTHPN = (1 << 16) /* do linux clone exit parent notification */
	RFSTOPPED   = (1 << 17) /* leave child in a stopped state */
	RFHIGHPID   = (1 << 18) /* use a pid higher than 10 (idleproc) */
	RFTSIGZMB   = (1 << 19) /* select signal for exit parent notification */
	RFTSIGSHIFT = 20        /* selected signal number is in bits 20-27  */
	RFTSIGMASK  = 0xFF
	RFPROCDESC  = (1 << 28) /* return a process descriptor */
	/* kernel: parent sleeps until child exits (vfork) */
	RFPPWAIT = (1 << 31)
	/* user: vfork(2) semantics, clear signals */
	RFSPAWN      = (1 << 31)
	RFFLAGS      = (RFFDG | RFPROC | RFMEM | RFNOWAIT | RFCFDG | RFTHREAD | RFSIGSHARE | RFLINUXTHPN | RFSTOPPED | RFHIGHPID | RFTSIGZMB | RFPROCDESC | RFSPAWN | RFPPWAIT)
	RFKERNELONLY = (RFSTOPPED | RFHIGHPID | RFPROCDESC)
)
