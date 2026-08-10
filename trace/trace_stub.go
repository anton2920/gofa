//go:build !gofatrace
// +build !gofatrace

package trace

type (
	Block    struct{}
	Profiler struct{}
)

//go:nosplit
func (p Profiler) BeginBody(_ uintptr, _ string) int { return 0 }

//go:nosplit
func (p Profiler) Begin(_ string) int { return 0 }

func (p Profiler) BeginProfile() {}

//go:nosplit
func (b Block) End() {}

//go:nosplit
func (p Profiler) EndProfile() {}

func (p Profiler) DumpProfile(_ []byte) int { return 0 }
