package syscall

const (
	/* From <aio.h>. */
	/*
	 * Returned by aio_cancel:
	 */
	AIO_CANCELED    = 0x1
	AIO_NOTCANCELED = 0x2
	AIO_ALLDONE     = 0x3

	/*
	 * LIO opcodes
	 */
	LIO_NOP      = 0x0
	LIO_WRITE    = 0x1
	LIO_READ     = 0x2
	LIO_VECTORED = 0x4
	LIO_WRITEV   = (LIO_WRITE | LIO_VECTORED)
	LIO_READV    = (LIO_READ | LIO_VECTORED)
	LIO_SYNC     = 0x8
	LIO_DSYNC    = (0x10 | LIO_SYNC)
	LIO_MLOCK    = 0x20

	/*
	 * LIO modes
	 */
	LIO_NOWAIT = 0x0
	LIO_WAIT   = 0x1
)

/* From <aio.h>. */
type Aiocb struct {
	Fildes     int32  /* File descriptor */
	Offset     int64  /* File offset for I/O */
	Buf        []byte /* I/O buffer in process space + Number of bytes for I/O  + int __spare__[2] */
	_          uintptr
	LioOpcode  int32      /* LIO opcode */
	AioReqprio int32      /* Request priority -- ignored */
	_          [3]uintptr /* Private members for aiocb */
	Sigevent   Sigevent   /* Signal to deliver */
}
