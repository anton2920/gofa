package syscall

import (
	"syscall"
	"unsafe"

	"github.com/anton2920/gofa/strings"
)

const (
	/* From <sys/syscall.h>. */
	SYS_accept           = 30
	SYS_access           = 33
	SYS_aio_cancel       = 316
	SYS_aio_error        = 317
	SYS_aio_read         = 255
	SYS_aio_return       = 314
	SYS_aio_suspend      = 315
	SYS_aio_write        = 256
	SYS_bind             = 104
	SYS_clock_gettime    = 232
	SYS_close            = 6
	SYS_exit             = 1
	SYS_fcntl            = 92
	SYS_fstat            = 551
	SYS_ftruncate        = 480
	SYS_getrandom        = 563
	SYS_ioctl            = 54
	SYS_jail_remove      = 508
	SYS_jail_set         = 507
	SYS_kevent           = 560
	SYS_kill             = 37
	SYS_kqueue           = 362
	SYS_listen           = 106
	SYS_lseek            = 478
	SYS_mkdir            = 136
	SYS_mmap             = 477
	SYS_munmap           = 73
	SYS_nanosleep        = 240
	SYS_nmount           = 378
	SYS_open             = 5
	SYS_pread            = 475
	SYS_pwrite           = 476
	SYS_rctl_add_rule    = 528
	SYS_rctl_remove_rule = 529
	SYS_read             = 3
	SYS_rmdir            = 137
	SYS_setsockopt       = 105
	SYS_sigaction        = 416
	SYS_shm_open2        = 571
	SYS_shutdown         = 134
	SYS_socket           = 97
	SYS_stat             = 188
	SYS_unlink           = 10
	SYS_unmount          = 22
	SYS_write            = 4
	SYS_writev           = 121
)

func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, errno uintptr)

func Syscall(trap, a1, a2, a3 uintptr) (uintptr, uintptr, uintptr) {
	r1, r2, errno := syscall.Syscall(trap, a1, a2, a3)
	return r1, r2, uintptr(errno)
}

func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)

func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (uintptr, uintptr, uintptr) {
	r1, r2, errno := syscall.Syscall6(trap, a1, a2, a3, a4, a5, a6)
	return r1, r2, uintptr(errno)
}

func Accept(s int32, addr *Sockaddr, addrlen *uint32) (int32, error) {
	r1, _, errno := Syscall(SYS_accept, uintptr(s), uintptr(unsafe.Pointer(addr)), uintptr(unsafe.Pointer(addrlen)))
	return int32(r1), NewError("accept", errno)
}

func Access(path string, mode int32) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_access, uintptr(unsafe.Pointer(&buffer[0])), uintptr(mode), 0)
	return NewError("access", errno)
}

func AioRead(aiocb *Aiocb) error {
	_, _, errno := RawSyscall(SYS_aio_read, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return NewError("aio_read", errno)
}

func AioReturn(aiocb *Aiocb) (int, error) {
	r1, _, errno := RawSyscall(SYS_aio_return, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return int(r1), NewError("aio_return", errno)
}

func AioWrite(aiocb *Aiocb) error {
	_, _, errno := RawSyscall(SYS_aio_write, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return NewError("aio_write", errno)
}

func Bind(s int32, addr *Sockaddr, addrlen uint32) error {
	_, _, errno := RawSyscall(SYS_bind, uintptr(s), uintptr(unsafe.Pointer(addr)), uintptr(addrlen))
	return NewError("bind", errno)
}

func ClockGettime(clockID int32, tp *Timespec) error {
	_, _, errno := RawSyscall(SYS_clock_gettime, uintptr(clockID), uintptr(unsafe.Pointer(tp)), 0)
	return NewError("clock_gettime", errno)
}

func Close(fd int32) error {
	_, _, errno := Syscall(SYS_close, uintptr(fd), 0, 0)
	return NewError("close", errno)
}

func Exit(status int32) {
	RawSyscall(SYS_exit, uintptr(status), 0, 0)
}

func Fcntl(fd, cmd int32, arg int32) (int32, error) {
	r1, _, errno := Syscall(SYS_fcntl, uintptr(fd), uintptr(cmd), uintptr(arg))
	return int32(r1), NewError("fcntl", errno)
}

func Fstat(fd int32, sb *Stat_t) error {
	_, _, errno := RawSyscall(SYS_fstat, uintptr(fd), uintptr(unsafe.Pointer(sb)), 0)
	return NewError("fstat", errno)
}

func Ftruncate(fd int32, length int64) error {
	_, _, errno := RawSyscall(SYS_ftruncate, uintptr(fd), uintptr(length), 0)
	return NewError("ftruncate", errno)
}

func Ioctl(fd int32, request uint, argp unsafe.Pointer) error {
	_, _, errno := RawSyscall(SYS_ioctl, uintptr(fd), uintptr(request), uintptr(argp))
	return NewError("ioctl", errno)
}

func Getrandom(buf []byte, flags uint32) (int64, error) {
	r1, _, errno := Syscall(SYS_getrandom, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(flags))
	return int64(r1), NewError("getrandom", errno)
}

func JailRemove(jid int32) error {
	_, _, errno := RawSyscall(SYS_jail_remove, uintptr(jid), 0, 0)
	return NewError("jail_remove", errno)
}

func JailSet(iovs []Iovec, flags int32) (int32, error) {
	jid, _, errno := RawSyscall(SYS_jail_set, uintptr(unsafe.Pointer(&iovs[0])), uintptr(len(iovs)), uintptr(flags))
	return int32(jid), NewError("jail_set", errno)
}

func Kevent(kq int32, changelist []Kevent_t, eventlist []Kevent_t, timeout *Timespec) (int, error) {
	var chptr, evptr unsafe.Pointer
	if len(changelist) > 0 {
		chptr = unsafe.Pointer(&changelist[0])
	}
	if len(eventlist) > 0 {
		evptr = unsafe.Pointer(&eventlist[0])
	}

	r1, _, errno := Syscall6(SYS_kevent, uintptr(kq), uintptr(chptr), uintptr(len(changelist)), uintptr(evptr), uintptr(len(eventlist)), uintptr(unsafe.Pointer(timeout)))
	return int(r1), NewError("kevent", errno)
}

func Kill(pid int32, sig Signal) error {
	_, _, errno := RawSyscall(SYS_kill, uintptr(pid), uintptr(sig), 0)
	return NewError("kill", errno)
}

func Kqueue() (int32, error) {
	r1, _, errno := RawSyscall(SYS_kqueue, 0, 0, 0)
	return int32(r1), NewError("kqueue", errno)
}

func Listen(s int32, backlog int32) error {
	_, _, errno := RawSyscall(SYS_listen, uintptr(s), uintptr(backlog), 0)
	return NewError("listen", errno)
}

func Lseek(fd int32, offset int64, whence int32) (int64, error) {
	r1, _, errno := RawSyscall(SYS_lseek, uintptr(fd), uintptr(offset), uintptr(whence))
	return int64(r1), NewError("lseek", errno)
}

func Mkdir(path string, mode int16) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_mkdir, uintptr(unsafe.Pointer(&buffer[0])), uintptr(mode), 0)
	return NewError("mkdir", errno)
}

func Mmap(addr unsafe.Pointer, len uint64, prot, flags, fd int32, offset int64) (unsafe.Pointer, error) {
	r1, _, errno := RawSyscall6(SYS_mmap, uintptr(addr), uintptr(len), uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset))
	return unsafe.Pointer(r1), NewError("mmap", errno)
}

func Munmap(addr unsafe.Pointer, len uint64) error {
	_, _, errno := RawSyscall(SYS_munmap, uintptr(addr), uintptr(len), 0)
	return NewError("munmap", errno)
}

func Nanosleep(rqtp, rmtp *Timespec) error {
	_, _, errno := Syscall(SYS_nanosleep, uintptr(unsafe.Pointer(rqtp)), uintptr(unsafe.Pointer(rmtp)), 0)
	return NewError("nanosleep", errno)
}

func Nmount(iovs []Iovec, flags int32) error {
	_, _, errno := RawSyscall(SYS_nmount, uintptr(unsafe.Pointer(&iovs[0])), uintptr(len(iovs)), uintptr(flags))
	return NewError("nmount", errno)
}

func Open(path string, flags int32, mode uint16) (int32, error) {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	r1, _, errno := Syscall(SYS_open, uintptr(unsafe.Pointer(&buffer[0])), uintptr(flags), uintptr(mode))
	return int32(r1), NewError("open", errno)
}

func Pread(fd int32, buf []byte, offset int64) (int, error) {
	r1, _, errno := Syscall6(SYS_pread, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(offset), 0, 0)
	return int(r1), NewError("pread", errno)
}

func Pwrite(fd int32, buf []byte, offset int64) (int, error) {
	r1, _, errno := Syscall6(SYS_pwrite, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(offset), 0, 0)
	return int(r1), NewError("pwrite", errno)
}

func RctlAddRule(rule []byte) error {
	_, _, errno := RawSyscall6(SYS_rctl_add_rule, uintptr(unsafe.Pointer(&rule[0])), uintptr(len(rule)), 0, 0, 0, 0)
	return NewError("rctl_add_rule", errno)
}

func RctlRemoveRule(filter []byte) error {
	_, _, errno := RawSyscall6(SYS_rctl_remove_rule, uintptr(unsafe.Pointer(&filter[0])), uintptr(len(filter)), 0, 0, 0, 0)
	return NewError("rctl_remove_rule", errno)
}

func Read(fd int32, buf []byte) (int64, error) {
	r1, _, errno := Syscall(SYS_read, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int64(r1), NewError("read", errno)
}

func Rmdir(path string) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_rmdir, uintptr(unsafe.Pointer(&buffer[0])), 0, 0)
	return NewError("rmdir", errno)
}

func Setsockopt(s, level, optname int32, optval unsafe.Pointer, optlen uint32) error {
	_, _, errno := RawSyscall6(SYS_setsockopt, uintptr(s), uintptr(level), uintptr(optname), uintptr(optval), uintptr(optlen), 0)
	return NewError("setsockopt", errno)
}

func Sigaction(sig int32, act *Sigaction_t, oact *Sigaction_t) error {
	_, _, errno := RawSyscall(SYS_sigaction, uintptr(sig), uintptr(unsafe.Pointer(act)), uintptr(unsafe.Pointer(oact)))
	return NewError("sigaction", errno)
}

func ShmOpen2(path string, flags int32, mode uint16, shmflags int32, name string) (int32, error) {
	r1, _, errno := RawSyscall6(SYS_shm_open2, uintptr(unsafe.Pointer(strings.Data(path))), uintptr(flags), uintptr(mode), uintptr(shmflags), uintptr(unsafe.Pointer(strings.Data(name))), 0)
	return int32(r1), NewError("shm_open2", errno)
}

func Shutdown(s int32, how int32) error {
	_, _, errno := RawSyscall(SYS_shutdown, uintptr(s), uintptr(how), 0)
	return NewError("shutdown", errno)
}

func Socket(domain, typ, protocol int32) (int32, error) {
	r1, _, errno := RawSyscall(SYS_socket, uintptr(domain), uintptr(typ), uintptr(protocol))
	return int32(r1), NewError("socket", errno)
}

func Stat(path string, sb *Stat_t) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_stat, uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(sb)), 0)
	return NewError("stat", errno)
}

func Unlink(path string) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_unlink, uintptr(unsafe.Pointer(&buffer[0])), 0, 0)
	return NewError("unlink", errno)
}

func Unmount(path string, flags int32) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_unmount, uintptr(unsafe.Pointer(&buffer[0])), uintptr(flags), 0)
	return NewError("unmount", errno)
}

func Write(fd int32, buf []byte) (int64, error) {
	r1, _, errno := Syscall(SYS_write, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int64(r1), NewError("write", errno)
}

func Writev(fd int32, iov []Iovec) (int64, error) {
	r1, _, errno := Syscall(SYS_writev, uintptr(fd), uintptr(unsafe.Pointer(&iov[0])), uintptr(len(iov)))
	return int64(r1), NewError("writev", errno)
}
