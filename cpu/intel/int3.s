/* From "textflag.h". */
#define NOSPLIT	4

/* func INT3() */
TEXT ·INT3(SB), NOSPLIT, $0
	BYTE $0xCC
	RET
