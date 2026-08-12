//go:build freebsd
// +build freebsd

package freebsd

import (
	"unsafe"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/strings"
)

const (
	/* From <sys/syscall.h>. */
	SYS_accept           = 30
	SYS_access           = 33
	SYS_aio_cancel       = 316
	SYS_aio_bool         = 317
	SYS_aio_read         = 255
	SYS_aio_return       = 314
	SYS_aio_suspend      = 315
	SYS_aio_write        = 256
	SYS_bind             = 104
	SYS_clock_gettime    = 232
	SYS_close            = 6
	SYS_connect          = 98
	SYS_exit             = 1
	SYS_fcntl            = 92
	SYS_fstat            = 551
	SYS_fsync            = 95
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
	SYS_madvise          = 75
	SYS_mkdir            = 136
	SYS_mmap             = 477
	SYS_munmap           = 73
	SYS_mprotect         = 74
	SYS_nanosleep        = 240
	SYS_nmount           = 378
	SYS_open             = 5
	SYS_openat           = 499
	SYS_pread            = 475
	SYS_pwrite           = 476
	SYS_rctl_add_rule    = 528
	SYS_rctl_remove_rule = 529
	SYS_read             = 3
	SYS_rmdir            = 137
	SYS_setsockopt       = 105
	SYS_shm_open2        = 571
	SYS_shutdown         = 134
	SYS_sigaction        = 416
	SYS_socket           = 97
	SYS_stat             = 188
	SYS_unlink           = 10
	SYS_unmount          = 22
	SYS_write            = 4
	SYS_writev           = 121
)

//go:nosplit
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, errno uintptr)

//go:nosplit
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2, errno uintptr)

//go:nosplit
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)

//go:nosplit
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, errno uintptr)

//go:nosplit
func Accept(ctx *context.Context, s int32, addr *Sockaddr, addrlen *uint32) (int32, bool) {
	r1, _, errno := Syscall(SYS_accept, uintptr(s), uintptr(unsafe.Pointer(addr)), uintptr(unsafe.Pointer(addrlen)))
	return int32(r1), ReportPotentialError(ctx, errno)
}

func Access(ctx *context.Context, path string, mode int32) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_access, uintptr(unsafe.Pointer(&buffer[0])), uintptr(mode), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func AioRead(ctx *context.Context, aiocb *Aiocb) bool {
	_, _, errno := RawSyscall(SYS_aio_read, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func AioReturn(ctx *context.Context, aiocb *Aiocb) (int, bool) {
	r1, _, errno := RawSyscall(SYS_aio_return, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return int(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func AioWrite(ctx *context.Context, aiocb *Aiocb) bool {
	_, _, errno := RawSyscall(SYS_aio_write, uintptr(unsafe.Pointer(aiocb)), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Bind(ctx *context.Context, s int32, addr *Sockaddr, addrlen uint32) bool {
	_, _, errno := RawSyscall(SYS_bind, uintptr(s), uintptr(unsafe.Pointer(addr)), uintptr(addrlen))
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func ClockGettime(ctx *context.Context, clockID int32, tp *Timespec) bool {
	_, _, errno := RawSyscall(SYS_clock_gettime, uintptr(clockID), uintptr(unsafe.Pointer(tp)), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Close(ctx *context.Context, fd int32) bool {
	_, _, errno := Syscall(SYS_close, uintptr(fd), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Connect(ctx *context.Context, s int32, name *Sockaddr, namelen uint32) bool {
	_, _, errno := Syscall(SYS_connect, uintptr(s), uintptr(unsafe.Pointer(name)), uintptr(namelen))
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Exit(status int32) {
	RawSyscall(SYS_exit, uintptr(status), 0, 0)
}

//go:nosplit
func Fcntl(ctx *context.Context, fd, cmd int32, arg int32) (int32, bool) {
	r1, _, errno := Syscall(SYS_fcntl, uintptr(fd), uintptr(cmd), uintptr(arg))
	return int32(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Fstat(ctx *context.Context, fd int32, sb *Stat_t) bool {
	_, _, errno := RawSyscall(SYS_fstat, uintptr(fd), uintptr(unsafe.Pointer(sb)), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Fsync(ctx *context.Context, fd int32) bool {
	_, _, errno := RawSyscall(SYS_fsync, uintptr(fd), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Ftruncate(ctx *context.Context, fd int32, length int64) bool {
	_, _, errno := RawSyscall(SYS_ftruncate, uintptr(fd), uintptr(length), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Ioctl(ctx *context.Context, fd int32, request uint, argp unsafe.Pointer) bool {
	_, _, errno := RawSyscall(SYS_ioctl, uintptr(fd), uintptr(request), uintptr(argp))
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Getrandom(ctx *context.Context, buf []byte, flags uint32) (int64, bool) {
	r1, _, errno := Syscall(SYS_getrandom, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(flags))
	return int64(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func JailRemove(ctx *context.Context, jid int32) bool {
	_, _, errno := RawSyscall(SYS_jail_remove, uintptr(jid), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func JailSet(ctx *context.Context, iovs []Iovec, flags int32) (int32, bool) {
	jid, _, errno := RawSyscall(SYS_jail_set, uintptr(unsafe.Pointer(&iovs[0])), uintptr(len(iovs)), uintptr(flags))
	return int32(jid), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Kevent(ctx *context.Context, kq int32, changelist []Kevent_t, eventlist []Kevent_t, timeout *Timespec) (int32, bool) {
	var chptr, evptr unsafe.Pointer
	if len(changelist) > 0 {
		chptr = unsafe.Pointer(&changelist[0])
	}
	if len(eventlist) > 0 {
		evptr = unsafe.Pointer(&eventlist[0])
	}

	r1, _, errno := Syscall6(SYS_kevent, uintptr(kq), uintptr(chptr), uintptr(len(changelist)), uintptr(evptr), uintptr(len(eventlist)), uintptr(unsafe.Pointer(timeout)))
	return int32(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Kill(ctx *context.Context, pid int32, sig Signal) bool {
	_, _, errno := RawSyscall(SYS_kill, uintptr(pid), uintptr(sig), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Kqueue(ctx *context.Context) (int32, bool) {
	r1, _, errno := RawSyscall(SYS_kqueue, 0, 0, 0)
	return int32(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Listen(ctx *context.Context, s int32, backlog int32) bool {
	_, _, errno := RawSyscall(SYS_listen, uintptr(s), uintptr(backlog), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Lseek(ctx *context.Context, fd int32, offset int64, whence int32) (int64, bool) {
	r1, _, errno := RawSyscall(SYS_lseek, uintptr(fd), uintptr(offset), uintptr(whence))
	return int64(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Madvise(ctx *context.Context, addr unsafe.Pointer, len uint, behav int32) bool {
	_, _, errno := RawSyscall(SYS_madvise, uintptr(addr), uintptr(len), uintptr(behav))
	return ReportPotentialError(ctx, errno)
}

func Mkdir(ctx *context.Context, path string, mode int16) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_mkdir, uintptr(unsafe.Pointer(&buffer[0])), uintptr(mode), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Mmap(ctx *context.Context, addr unsafe.Pointer, len uint, prot int32, flags int32, fd int32, offset int64) (unsafe.Pointer, bool) {
	r1, _, errno := RawSyscall6(SYS_mmap, uintptr(addr), uintptr(len), uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset))
	return unsafe.Pointer(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Munmap(ctx *context.Context, addr unsafe.Pointer, len uint) bool {
	_, _, errno := RawSyscall(SYS_munmap, uintptr(addr), uintptr(len), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Mprotect(ctx *context.Context, addr unsafe.Pointer, len uint, prot int32) bool {
	_, _, errno := RawSyscall(SYS_mprotect, uintptr(addr), uintptr(len), uintptr(prot))
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Nanosleep(ctx *context.Context, rqtp, rmtp *Timespec) bool {
	_, _, errno := Syscall(SYS_nanosleep, uintptr(unsafe.Pointer(rqtp)), uintptr(unsafe.Pointer(rmtp)), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Nmount(ctx *context.Context, iovs []Iovec, flags int32) bool {
	_, _, errno := RawSyscall(SYS_nmount, uintptr(unsafe.Pointer(&iovs[0])), uintptr(len(iovs)), uintptr(flags))
	return ReportPotentialError(ctx, errno)
}

func Open(ctx *context.Context, path string, flags int32, mode uint16) (int32, bool) {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	r1, _, errno := Syscall(SYS_open, uintptr(unsafe.Pointer(&buffer[0])), uintptr(flags), uintptr(mode))
	return int32(r1), ReportPotentialError(ctx, errno)
}

func OpenAt(ctx *context.Context, fd int32, path string, flags int32, mode uint16) (int32, bool) {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	r1, _, errno := Syscall6(SYS_openat, uintptr(fd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(flags), uintptr(mode), 0, 0)
	return int32(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Pread(ctx *context.Context, fd int32, buf []byte, offset int64) (int, bool) {
	r1, _, errno := Syscall6(SYS_pread, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(offset), 0, 0)
	return int(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Pwrite(ctx *context.Context, fd int32, buf []byte, offset int64) (int, bool) {
	r1, _, errno := Syscall6(SYS_pwrite, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), uintptr(offset), 0, 0)
	return int(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func RctlAddRule(ctx *context.Context, rule []byte) bool {
	_, _, errno := RawSyscall6(SYS_rctl_add_rule, uintptr(unsafe.Pointer(&rule[0])), uintptr(len(rule)), 0, 0, 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func RctlRemoveRule(ctx *context.Context, filter []byte) bool {
	_, _, errno := RawSyscall6(SYS_rctl_remove_rule, uintptr(unsafe.Pointer(&filter[0])), uintptr(len(filter)), 0, 0, 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Read(ctx *context.Context, fd int32, buf []byte) (int, bool) {
	r1, _, errno := Syscall(SYS_read, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(r1), ReportPotentialError(ctx, errno)
}

func Rmdir(ctx *context.Context, path string) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_rmdir, uintptr(unsafe.Pointer(&buffer[0])), 0, 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Setsockopt(ctx *context.Context, s, level, optname int32, optval unsafe.Pointer, optlen uint32) bool {
	_, _, errno := RawSyscall6(SYS_setsockopt, uintptr(s), uintptr(level), uintptr(optname), uintptr(optval), uintptr(optlen), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Sigaction(ctx *context.Context, sig int32, act *Sigaction_t, oact *Sigaction_t) bool {
	_, _, errno := RawSyscall(SYS_sigaction, uintptr(sig), uintptr(unsafe.Pointer(act)), uintptr(unsafe.Pointer(oact)))
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func ShmOpen2(ctx *context.Context, path string, flags int32, mode uint16, shmflags int32, name string) (int32, bool) {
	r1, _, errno := RawSyscall6(SYS_shm_open2, uintptr(unsafe.Pointer(strings.Data(path))), uintptr(flags), uintptr(mode), uintptr(shmflags), uintptr(unsafe.Pointer(strings.Data(name))), 0)
	return int32(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Shutdown(ctx *context.Context, s int32, how int32) bool {
	_, _, errno := RawSyscall(SYS_shutdown, uintptr(s), uintptr(how), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Socket(ctx *context.Context, domain int32, typ int32, protocol int32) (int32, bool) {
	r1, _, errno := RawSyscall(SYS_socket, uintptr(domain), uintptr(typ), uintptr(protocol))
	return int32(r1), ReportPotentialError(ctx, errno)
}

func Stat(ctx *context.Context, path string, sb *Stat_t) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_stat, uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(sb)), 0)
	return ReportPotentialError(ctx, errno)
}

func Unlink(ctx *context.Context, path string) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_unlink, uintptr(unsafe.Pointer(&buffer[0])), 0, 0)
	return ReportPotentialError(ctx, errno)
}

func Unmount(ctx *context.Context, path string, flags int32) bool {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_unmount, uintptr(unsafe.Pointer(&buffer[0])), uintptr(flags), 0)
	return ReportPotentialError(ctx, errno)
}

//go:nosplit
func Write(ctx *context.Context, fd int32, buf []byte) (int, bool) {
	r1, _, errno := Syscall(SYS_write, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(r1), ReportPotentialError(ctx, errno)
}

//go:nosplit
func Writev(ctx *context.Context, fd int32, iov []Iovec) (int, bool) {
	r1, _, errno := Syscall(SYS_writev, uintptr(fd), uintptr(unsafe.Pointer(&iov[0])), uintptr(len(iov)))
	return int(r1), ReportPotentialError(ctx, errno)
}
