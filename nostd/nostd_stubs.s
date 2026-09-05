//go:build gofanostd && (gofanostd13 || gofanostd14 || gofanostd15 || gofanostd16 || gofanostdxx)
// +build gofanostd
// +build gofanostd13 gofanostd14 gofanostd15 gofanostd16 gofanostdxx

TEXT runtime·memcopy(SB), 6, $-0
	JMP	·memcopy(SB)

TEXT runtime·memcopy0(SB), 6, $-0
	RET

TEXT runtime·memcopy8(SB), 6, $-0
	JMP	·memcopy8(SB)

TEXT runtime·memcopy16(SB), 6, $-0
	JMP	·memcopy16(SB)

TEXT runtime·memcopy32(SB), 6, $-0
	JMP	·memcopy32(SB)

TEXT runtime·memcopy64(SB), 6, $-0
	JMP	·memcopy64(SB)

TEXT runtime·memcopy128(SB), 6, $-0
	JMP	·memcopy128(SB)

TEXT runtime·typedmemmove(SB), 6, $-0
	JMP	·typedmemmove(SB)

TEXT runtime·typedslicecopy(SB), 6, $-0
	JMP	·typedslicecopy(SB)

TEXT runtime·memmove(SB), 6, $-0
	JMP	·memmove(SB)

TEXT runtime·memequal64(SB), 6, $-0
	JMP	·memequal64(SB)

TEXT runtime·cmpstring(SB), 6, $-0
	JMP	·cmpstring(SB)

TEXT runtime·cmpbytes(SB), 6, $-0
	JMP	·cmpbytes(SB)

TEXT runtime·growslice(SB), 6, $-0
	JMP	runtime·gopanic(SB)

TEXT runtime·memhash(SB), 6, $-0
	JMP	·memhash(SB)

TEXT runtime·strhash(SB), 6, $-0
	JMP	·strhash(SB)

TEXT runtime·memequal(SB), 6, $-0
	JMP	·memequal(SB)

TEXT runtime·eqstring(SB), 6, $-0
	JMP	·eqstring(SB)

TEXT runtime·writebarrierfat(SB), 6, $-0
	JMP	·writebarrierfat(SB)

TEXT runtime·writebarrierptr(SB), 6, $-0
	JMP	·writebarrierptr(SB)

TEXT runtime·writebarrierslice(SB), 6, $-0
	JMP	·writebarrierslice(SB)


TEXT runtime·memprint(SB), 6, $-0
	RET

TEXT runtime·morestack_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack00_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack01_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack10_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack11_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack8_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack16_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack24_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack32_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack40_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack48_noctxt(SB), 6, $-0
	RET

TEXT runtime·morestack(SB), 6, $-0
	RET

TEXT runtime·morestack00(SB), 6, $-0
	RET

TEXT runtime·morestack01(SB), 6, $-0
	RET

TEXT runtime·morestack10(SB), 6, $-0
	RET

TEXT runtime·morestack11(SB), 6, $-0
	RET

TEXT runtime·morestack8(SB), 6, $-0
	RET

TEXT runtime·morestack16(SB), 6, $-0
	RET

TEXT runtime·morestack24(SB), 6, $-0
	RET

TEXT runtime·morestack32(SB), 6, $-0
	RET

TEXT runtime·morestack40(SB), 6, $-0
	RET

TEXT runtime·morestack48(SB), 6, $-0
	RET


GLOBL runtime·algarray(SB), 8, $192
GLOBL runtime·writeBarrier(SB), 8, $8
GLOBL runtime·writeBarrierEnabled(SB), 24, $8
GLOBL runtime·firstmoduledata(SB), 8, $1024
