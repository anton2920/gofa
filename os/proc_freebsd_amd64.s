//go:build amd64
//+build amd64

/* From "textflag.h". */
#define NOSPLIT	4


/* func AsSignalHandler(fn func(Signal)) SignalHandler */
TEXT ·AsSignalHandler(SB), NOSPLIT, $0-16
	MOVQ	fn+0(FP), AX
	MOVQ	(AX), AX
	MOVQ	AX, ret+8(FP)
	RET
