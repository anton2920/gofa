//go:build !freebsd
// +build !freebsd

package jail

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync/atomic"
	stdsyscall "syscall"

	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/util"
)

type Jail struct {
	Index            uint32
	WorkingDirectory string
}

var (
	JailsRootDir = "./jails"

	JailLastIndex uint32
)

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

func New(template string, wd string) (Jail, error) {
	var j Jail

	j.WorkingDirectory = wd
	j.Index = atomic.AddUint32(&JailLastIndex, 1)

	env := make([]byte, syscall.PATH_MAX)
	n := PutEnv(env, j)
	env = env[:n+1]

	if err := syscall.Mkdir(util.Slice2String(env), 0755); err != nil {
		if err.(syscall.Error).Errno != syscall.EEXIST {
			return Jail{}, fmt.Errorf("failed to create environment directory: %v", err)
		}
	}

	return j, nil
}

func Protect(j Jail) error {
	log.Warnf("Failed to protect working directory: unsupported on %q!", runtime.GOOS)
	return nil
}

func Remove(j Jail) error {
	env := make([]byte, syscall.PATH_MAX)
	n := PutEnv(env, j)
	env = env[:n+1]

	if err := syscall.Rmdir(util.Slice2String(env)); err != nil {
		return fmt.Errorf("failed to remove environment directory: %v", err)
	}

	return nil
}

func Command(j Jail, exe string, args ...string) *exec.Cmd {
	cmd := exec.Command(exe, args...)
	cmd.SysProcAttr = &stdsyscall.SysProcAttr{Setsid: true}
	cmd.Dir = "/tmp"
	return cmd
}
