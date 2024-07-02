#include "textflag.h"

TEXT ·RDTSC(SB), NOSPLIT, $0
	RDTSC
	MOVL AX, ret+0(FP)
	MOVL DX, ret+4(FP)
	RET
