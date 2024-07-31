#include "textflag.h"

/* func LFENCE() */
TEXT ·LFENCE(SB), NOSPLIT, $0
	LFENCE
	RET

/* func MFENCE() */
TEXT ·MFENCE(SB), NOSPLIT, $0
	MFENCE
	RET

/* func SFENCE() */
TEXT ·SFENCE(SB), NOSPLIT, $0
	SFENCE
	RET
