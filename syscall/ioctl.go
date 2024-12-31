package syscall

/* From <sys/ioccom.h>. */
const (
	IOCPARM_SHIFT = 13                         /* number of bits for ioctl size */
	IOCPARM_MASK  = ((1 << IOCPARM_SHIFT) - 1) /* parameter length mask */

	IOC_OUT   = 0x40000000         /* copy out parameters */
	IOC_IN    = 0x80000000         /* copy in parameters */
	IOC_INOUT = (IOC_IN | IOC_OUT) /* copy parameters in and out */
)

func IOC(inout uint, group uint, num uint, len uint) uint {
	return ((inout) | (((len) & IOCPARM_MASK) << 16) | ((group) << 8) | (num))
}

func IOR(group uint, num uint, len uint) uint {
	return IOC(IOC_OUT, group, num, len)
}

func IOW(group uint, num uint, len uint) uint {
	return IOC(IOC_IN, group, num, len)
}

func IOWR(group uint, num uint, len uint) uint {
	return IOC(IOC_INOUT, group, num, len)
}
