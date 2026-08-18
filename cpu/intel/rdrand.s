#include "textflag.h"


/* func RDRANDW() (uint16, bool) */
TEXT ·RDRANDW(SB), NOSPLIT, $0-3
	BYTE	$0x66;	BYTE	$0x0F;	BYTE	$0xC7;	BYTE	$0xF2	// RDRANDW	DX
	JCC	NotOKW
	MOVW	DX, n+0(FP)
	MOVB	$1, ok+2(FP)
	JMP	RetW
NotOKW:
	MOVW	$0, n+0(FP)
	MOVB	$0, ok+2(FP)
RetW:
	RET


/* func RDRANDL() (uint32, bool) */
TEXT ·RDRANDL(SB), NOSPLIT, $0-5
	BYTE	$0x0F;	BYTE	$0xC7;	BYTE	$0xF2 			// RDRANDL	DX
	JCC	NotOKL
	MOVL	DX, n+0(FP)
	MOVB	$1, ok+4(FP)
	JMP	RetL
NotOKL:
	MOVL	$0, n+0(FP)
	MOVB	$0, ok+4(FP)
RetL:
	RET


/* func RDRANDQ() (uint64, bool) */
TEXT ·RDRANDQ(SB), NOSPLIT, $0-9
	BYTE	$0x48;	BYTE	$0x0F;	BYTE	$0xC7;	BYTE	$0xF2	// RDRANDQ	DX
	JCC	NotOKQ
	MOVQ	DX, n+0(FP)
	MOVB	$1, ok+8(FP)
	JMP	RetQ
NotOKQ:
	MOVQ	$0, n+0(FP)
	MOVB	$0, ok+8(FP)
RetQ:
	RET


