package syscall

import (
	"syscall"
	"unsafe"

	"github.com/anton2920/gofa/util"
)

const (
	SYS_accept        = 43
	SYS_bind          = 49
	SYS_clock_gettime = 228
	SYS_close         = 3
	SYS_exit          = 60
	SYS_ftruncate     = 77
	SYS_kill          = 62
	SYS_listen        = 50
	SYS_memfd_create  = 319
	SYS_mkdir         = 83
	SYS_mmap          = 9
	SYS_munmap        = 11
	SYS_open          = 2
	SYS_pread         = 17
	SYS_pwrite        = 18
	SYS_read          = 0
	SYS_rmdir         = 84
	SYS_setsockopt    = 54
	SYS_sigaction     = 13
	SYS_socket        = 41
	SYS_unlink        = 87
	SYS_write         = 1
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

func Ftruncate(fd int32, length int64) error {
	_, _, errno := RawSyscall(SYS_ftruncate, uintptr(fd), uintptr(length), 0)
	return NewError("ftruncate", errno)
}

func Kill(pid int32, sig Signal) error {
	_, _, errno := RawSyscall(SYS_kill, uintptr(pid), uintptr(sig), 0)
	return NewError("kill", errno)
}

func Listen(s int32, backlog int32) error {
	_, _, errno := RawSyscall(SYS_listen, uintptr(s), uintptr(backlog), 0)
	return NewError("listen", errno)
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

func Read(fd int32, buf []byte) (int, error) {
	r1, _, errno := Syscall(SYS_read, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(r1), NewError("read", errno)
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
	r1, _, errno := RawSyscall(SYS_memfd_create, uintptr(unsafe.Pointer(util.StringData(path))), uintptr(flags), 0)
	return int32(r1), NewError("shm_open2", errno)
}

func Socket(domain, typ, protocol int32) (int32, error) {
	r1, _, errno := RawSyscall(SYS_socket, uintptr(domain), uintptr(typ), uintptr(protocol))
	return int32(r1), NewError("socket", errno)
}

func Unlink(path string) error {
	buffer := make([]byte, PATH_MAX+1)
	copy(buffer[:PATH_MAX], path)

	_, _, errno := RawSyscall(SYS_unlink, uintptr(unsafe.Pointer(&buffer[0])), 0, 0)
	return NewError("unlink", errno)
}

func Write(fd int32, buf []byte) (int, error) {
	r1, _, errno := Syscall(SYS_write, uintptr(fd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(r1), NewError("write", errno)
}
