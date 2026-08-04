//go:build 386 || amd64
// +build 386 amd64

package cpu

import "github.com/anton2920/gofa/cpu/intel"

func WaitForLoadOperationsToComplete() {
	intel.LFENCE()
}

func WaitForStoreOperationsToComplete() {
	intel.SFENCE()
}

func WaitForLoadAndStoreOperationsToComplete() {
	intel.MFENCE()
}
