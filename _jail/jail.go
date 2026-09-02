//go:build freebsd
// +build freebsd

package jail

import (
	"fmt"
	"os/exec"
	"sync/atomic"
	stdsyscall "syscall"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
)

type Jail struct {
	ID               int32
	Index            uint32
	WorkingDirectory string
}

const (
	JailNamePrefix = "gofa-"
	MaxJailNameLen = len(JailNamePrefix) + ints.Bufsize

	MaxJailRctlPrefixLen = len("jail:") + MaxJailNameLen + len(":")
	MaxJailRctlRuleLen   = MaxJailRctlPrefixLen + ints.Bufsize
)

var (
	JailsRootDir = "./jails"

	JailLastIndex uint32
)

func PutName(buffer []byte, j Jail) int {
	var n int

	n += copy(buffer[n:], JailNamePrefix)
	n += slices.PutInt(buffer[n:], int(j.Index))

	return n
}

func PutPath(buffer []byte, j Jail) int {
	var n int

	n += copy(buffer[n:], j.WorkingDirectory)

	buffer[n] = '/'
	n++

	n += copy(buffer[n:], JailsRootDir)

	buffer[n] = '/'
	n++

	n += copy(buffer[n:], "containers/")
	n += slices.PutInt(buffer[n:], int(j.Index))

	return n
}

func PutTmp(buffer []byte, j Jail) int {
	var n int

	n += PutPath(buffer[n:], j)
	n += copy(buffer[n:], "/tmp")

	return n
}

func PutEnv(buffer []byte, j Jail) int {
	var n int

	n += copy(buffer[n:], j.WorkingDirectory)

	buffer[n] = '/'
	n++

	n += copy(buffer[n:], JailsRootDir)

	buffer[n] = '/'
	n++

	n += copy(buffer[n:], "envs/")
	n += slices.PutInt(buffer[n:], int(j.Index))

	return n

}

func PutRctlPrefix(buffer []byte, name []byte) int {
	var n int

	n += copy(buffer[n:], "jail:")
	n += copy(buffer[n:], name)

	buffer[n] = ':'
	n++

	return n
}

func PutRctlRule(buffer []byte, prefix []byte, rule string) int {
	var n int

	n += copy(buffer[n:], prefix)
	n += copy(buffer[n:], rule)
	buffer[n] = '\x00'

	return n
}

func New(template string, wd string) (Jail, error) {
	var j Jail
	var err error

	j.WorkingDirectory = wd
	j.Index = atomic.AddUint32(&JailLastIndex, 1)

	tmpl := make([]byte, syscall.PATH_MAX)
	n := copy(tmpl, template)
	tmpl = tmpl[:n+1]

	name := make([]byte, MaxJailNameLen)
	n = PutName(name, j)
	name = name[:n+1]

	path := make([]byte, syscall.PATH_MAX)
	n = PutPath(path, j)
	path = path[:n+1]

	tmp := make([]byte, syscall.PATH_MAX)
	n = PutTmp(tmp, j)
	tmp = tmp[:n+1]

	env := make([]byte, syscall.PATH_MAX)
	n = PutEnv(env, j)
	env = env[:n+1]

	if err := syscall.Access(bytes.AsString(tmpl), 0); err != nil {
		return Jail{}, fmt.Errorf("provided template does not exist: %v", err)
	}

	if err := syscall.Mkdir(bytes.AsString(path), 0755); err != nil {
		if err.(syscall.Error).Errno != syscall.EEXIST {
			return Jail{}, fmt.Errorf("failed to create path: %v", err)
		}
	}

	if err := syscall.Mkdir(bytes.AsString(env), 0755); err != nil {
		if err.(syscall.Error).Errno != syscall.EEXIST {
			return Jail{}, fmt.Errorf("failed to create environment directory: %v", err)
		}
	}

	if err := syscall.Nmount([]syscall.Iovec{
		syscall.Iovec("target\x00"), syscall.IovecForByteSlice(tmpl),
		syscall.Iovec("fspath\x00"), syscall.IovecForByteSlice(path),
		syscall.Iovec("fstype\x00"), syscall.Iovec("nullfs\x00"),
		syscall.Iovec("ro\x00"), syscall.IovecZ,
	}, 0); err != nil {
		return Jail{}, fmt.Errorf("failed to mount container directory: %v", err)
	}

	if err := syscall.Nmount([]syscall.Iovec{
		syscall.Iovec("target\x00"), syscall.IovecForByteSlice(env),
		syscall.Iovec("fspath\x00"), syscall.IovecForByteSlice(tmp),
		syscall.Iovec("fstype\x00"), syscall.Iovec("nullfs\x00"),
		syscall.Iovec("rw\x00"), syscall.IovecZ,
	}, 0); err != nil {
		syscall.Unmount(bytes.AsString(path), 0)
		return Jail{}, fmt.Errorf("failed to mount environment directory: %v", err)
	}

	j.ID, err = syscall.JailSet([]syscall.Iovec{
		syscall.Iovec("host.hostname\x00"), syscall.Iovec("gofa-jail\x00"),
		syscall.Iovec("name\x00"), syscall.IovecForByteSlice(name),
		syscall.Iovec("path\x00"), syscall.IovecForByteSlice(path),
		syscall.Iovec("persist\x00"), syscall.IovecZ,
	}, syscall.JAIL_CREATE)
	if err != nil {
		syscall.Unmount(bytes.AsString(tmp), 0)
		syscall.Unmount(bytes.AsString(path), 0)
		return Jail{}, err
	}

	/* TODO(anton2920): replace with something that works in jail. */
	/*
		prefix := make([]byte, MaxJailRctlPrefixLen)
		n = PutRctlPrefix(prefix, name[:len(name)-1])
		prefix = prefix[:n+1]

		rule := make([]byte, MaxJailRctlRuleLen)
		rules := [...]string{
			"maxproc:deny=16",
			"vmemoryuse:deny=2684354560",
			"memoryuse:deny=536870912",
			"swapuse:deny=536870912",
		}
		for i := 0; i < len(rules); i++ {
			n := PutRctlRule(rule, prefix[:len(prefix)-1], rules[i])

			if err := syscall.RctlAddRule(rule[:n+1]); err != nil {
				syscall.RctlRemoveRule(prefix)
				syscall.JailRemove(j.ID)
				syscall.Unmount(bytes.AsString(tmp), 0)
				syscall.Unmount(bytes.AsString(path), 0)
				return Jail{}, fmt.Errorf("failed to add rule %d for jail: %v", i, err)
			}
		}
	*/

	return j, nil
}

func Protect(j Jail) error {
	tmp := make([]byte, syscall.PATH_MAX)
	n := PutTmp(tmp, j)
	tmp = tmp[:n+1]

	env := make([]byte, syscall.PATH_MAX)
	n = PutEnv(env, j)
	env = env[:n+1]

	if err := syscall.Unmount(bytes.AsString(tmp), 0); err != nil {
		return fmt.Errorf("failed to unmount environment: %v", err)
	}

	if err := syscall.Nmount([]syscall.Iovec{
		syscall.Iovec("target\x00"), syscall.IovecForByteSlice(env),
		syscall.Iovec("fspath\x00"), syscall.IovecForByteSlice(tmp),
		syscall.Iovec("fstype\x00"), syscall.Iovec("nullfs\x00"),
		syscall.Iovec("ro\x00"), syscall.IovecZ,
	}, 0); err != nil {
		return fmt.Errorf("failed to mount environment directory: %v", err)
	}

	return nil
}

func Remove(j Jail) error {
	var err error

	name := make([]byte, MaxJailNameLen)
	n := PutName(name, j)
	name = name[:n+1]

	path := make([]byte, syscall.PATH_MAX)
	n = PutPath(path, j)
	path = path[:n+1]

	tmp := make([]byte, syscall.PATH_MAX)
	n = PutTmp(tmp, j)
	tmp = tmp[:n+1]

	env := make([]byte, syscall.PATH_MAX)
	n = PutEnv(env, j)
	env = env[:n+1]

	prefix := make([]byte, MaxJailRctlPrefixLen)
	n = PutRctlPrefix(prefix, name[:len(name)-1])
	prefix = prefix[:n+1]

	if err1 := syscall.RctlRemoveRule(prefix); err1 != nil {
		err = fmt.Errorf("failed to remove jail rules: %v", err1)
	}

	if err1 := syscall.JailRemove(j.ID); err1 != nil {
		err = fmt.Errorf("failed to remove jail: %v", err1)
	}

	if err1 := syscall.Unmount(bytes.AsString(tmp), 0); err1 != nil {
		err = fmt.Errorf("failed to unmount environment: %v", err1)
	}

	if err1 := syscall.Unmount(bytes.AsString(path), 0); err1 != nil {
		err = fmt.Errorf("failed to unmount container: %v", err1)
	}

	if err1 := syscall.Rmdir(bytes.AsString(env)); err1 != nil {
		err = fmt.Errorf("failed to remove environment directory: %v", err1)
	}

	if err1 := syscall.Rmdir(bytes.AsString(path)); err1 != nil {
		err = fmt.Errorf("failed to remove container directory: %v", err1)
	}

	return err
}

func Command(j Jail, exe string, args ...string) *exec.Cmd {
	cmd := exec.Command(exe, args...)
	cmd.SysProcAttr = &stdsyscall.SysProcAttr{Setsid: true, Jail: int(j.ID)}
	cmd.Dir = "/tmp"
	return cmd
}
