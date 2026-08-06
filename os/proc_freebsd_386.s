//go:build 386
//+build 386

/* From "textflag.h". */
#define NOSPLIT	4


/* func AsSignalHandler(fn func(Signal)) SignalHandler */
TEXT ·AsSignalHandler(SB), NOSPLIT, $0-8
	MOVL	fn+0(FP), AX
	MOVL	(AX), AX
	MOVL	AX, ret+4(FP)
	RET
