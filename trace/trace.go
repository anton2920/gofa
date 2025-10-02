//go:build gofatrace
// +build gofatrace

package trace

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"unsafe"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/pointers"
	_ "github.com/anton2920/gofa/time"
)

type Anchor struct {
	ParentIndex int32

	PC    uintptr /* for lazy function name resolution. */
	Label string  /* for names of non-function blocks. */

	HitCount int

	ElapsedCyclesExclusive cpu.Cycles /* time for anchor itself. */
	ElapsedCyclesInclusive cpu.Cycles /* time for anchor plus time of its children. */
}

type Block struct {
	PC    uintptr /* for lazy function name resolution. */
	Label string  /* for names of non-function blocks. */

	ParentIndex int32
	AnchorIndex int32

	StartCycles               cpu.Cycles
	OldElapsedCyclesInclusive cpu.Cycles
}

type Anchors []Anchor

type Profiler struct {
	Anchors       Anchors
	CurrentParent int32

	StartCycles cpu.Cycles
	EndCycles   cpu.Cycles
}

var GlobalProfiler Profiler

func init() {
	/* NOTE(anton2920): len must be a power of two for fast modulus calculation. */
	GlobalProfiler.Anchors = make(Anchors, 1024)
}

func (as Anchors) Len() int { return len(as) }

func (as Anchors) Less(i, j int) bool {
	a := &as[i]
	b := &as[j]

	if (a.ElapsedCyclesInclusive > 0) && (b.ElapsedCyclesInclusive > 0) {
		if a.ElapsedCyclesInclusive < b.ElapsedCyclesInclusive {
			return false
		} else {
			return true
		}
	} else if a.ElapsedCyclesInclusive > 0 {
		if a.ElapsedCyclesInclusive < b.ElapsedCyclesExclusive {
			return false
		} else {
			return true
		}
	} else if b.ElapsedCyclesInclusive > 0 {
		if a.ElapsedCyclesExclusive < b.ElapsedCyclesInclusive {
			return false
		} else {
			return true
		}
	} else {
		if a.ElapsedCyclesExclusive < b.ElapsedCyclesExclusive {
			return false
		} else {
			return true
		}
	}
}

func (as Anchors) Swap(i, j int) { as[i], as[j] = as[j], as[i] }

/* GetCallerPC returns a value of %IP register that is going to be used by RET instruction. arg0 is the address of the first agrument function of interest accepts. */
//go:nosplit
func GetCallerPC(arg0 unsafe.Pointer) uintptr {
	return *(*uintptr)(pointers.Add(arg0, -int(unsafe.Sizeof(arg0))))
}

//go:nosplit
func anchorIndexForPC(pc uintptr) int32 {
	anchors := GlobalProfiler.Anchors

	start := int(pc) & (len(GlobalProfiler.Anchors) - 1)
	start += bools.ToInt(start == 0)
	if (pc == anchors[start].PC) || (anchors[start].PC == 0) {
		return int32(start)
	}

	var idx int
	for idx = start + 1; (pc != anchors[idx].PC) && (anchors[idx].PC != 0) && (idx != start); {
		idx = (idx + 1) & (len(GlobalProfiler.Anchors) - 1)
		idx += bools.ToInt(idx == 0)
	}
	if idx == start {
		panic("not enough space for new anchor")
	}

	return int32(idx)
}

//go:nosplit
func begin(pc uintptr, label string) Block {
	var b Block

	index := anchorIndexForPC(pc)
	b.ParentIndex = GlobalProfiler.CurrentParent
	GlobalProfiler.CurrentParent = index

	b.AnchorIndex = index
	b.Label = label
	b.PC = pc

	anchor := &GlobalProfiler.Anchors[b.AnchorIndex]
	b.OldElapsedCyclesInclusive = anchor.ElapsedCyclesInclusive

	b.StartCycles = cpu.ReadPerformanceCounter()
	return b
}

//go:nosplit
func Begin(label string) Block {
	cpu.LoadFence()
	return begin(GetCallerPC(unsafe.Pointer(&label)), label)
}

func BeginProfile() {
	for i := 0; i < len(GlobalProfiler.Anchors); i++ {
		GlobalProfiler.Anchors[i] = Anchor{}
	}
	GlobalProfiler.CurrentParent = 0
	GlobalProfiler.StartCycles = cpu.ReadPerformanceCounter()
}

func CyclesToNanoseconds(c cpu.Cycles) float64 {
	return 1000000000 * float64(c) / float64(cpu.CPUHz)
}

func CyclesToMilliseconds(c cpu.Cycles) float64 {
	return 1000 * float64(c) / float64(cpu.CPUHz)
}

//go:nosplit
func End(b Block) {
	elapsed := cpu.ReadPerformanceCounter() - b.StartCycles
	anchor := &GlobalProfiler.Anchors[b.AnchorIndex]
	parent := &GlobalProfiler.Anchors[b.ParentIndex]
	GlobalProfiler.CurrentParent = b.ParentIndex

	parent.ElapsedCyclesExclusive -= elapsed

	anchor.ElapsedCyclesInclusive = b.OldElapsedCyclesInclusive + elapsed
	anchor.ElapsedCyclesExclusive += elapsed
	anchor.HitCount++

	anchor.PC = b.PC
	anchor.Label = b.Label
	anchor.ParentIndex = b.ParentIndex
}

func PrintTimeElapsed(label string, totalElapsed cpu.Cycles, curr *Anchor, parent *Anchor) {
	percentTotal := 100 * (float64(curr.ElapsedCyclesExclusive) / float64(totalElapsed))
	percentParent := 100 * (float64(curr.ElapsedCyclesExclusive) / float64(parent.ElapsedCyclesInclusive))
	fmt.Fprintf(os.Stderr, "[trace]: \t %s[%d]: flat [%.4fms %.2f%%/%.2f%% %.2fns/op]", label, curr.HitCount, CyclesToMilliseconds(curr.ElapsedCyclesExclusive), percentTotal, percentParent, CyclesToNanoseconds(curr.ElapsedCyclesExclusive)/float64(curr.HitCount))

	if curr.ElapsedCyclesInclusive > curr.ElapsedCyclesExclusive {
		percentWithChildrenTotal := 100 * (float64(curr.ElapsedCyclesInclusive) / float64(totalElapsed))
		percentWithChildrenParent := 100 * (float64(curr.ElapsedCyclesInclusive) / float64(parent.ElapsedCyclesInclusive))
		fmt.Fprintf(os.Stderr, ", cum [%.4fms %.2f%%/%.2f%% %.2fns/op]", CyclesToMilliseconds(curr.ElapsedCyclesInclusive), percentWithChildrenTotal, percentWithChildrenParent, CyclesToNanoseconds(curr.ElapsedCyclesInclusive)/float64(curr.HitCount))
	}

	fmt.Fprintf(os.Stderr, "\n")
}

func EndAndPrintProfile() {
	GlobalProfiler.EndCycles = cpu.ReadPerformanceCounter()
	totalElapsed := GlobalProfiler.EndCycles - GlobalProfiler.StartCycles

	var totalCycles cpu.Cycles
	var totalHits int

	fmt.Fprintf(os.Stderr, "[trace]: Total time: %.4fms\n", CyclesToMilliseconds(totalElapsed))

	/* NOTE(anton2920): Anchor.ParentIndex uses original order, so we need to preserve it after Sort. */
	backup := make(Anchors, len(GlobalProfiler.Anchors))
	copy(backup, GlobalProfiler.Anchors)

	sort.Sort(GlobalProfiler.Anchors)

	for i := 0; i < len(GlobalProfiler.Anchors); i++ {
		anchor := &GlobalProfiler.Anchors[i]
		parent := &backup[anchor.ParentIndex]

		if anchor.HitCount > 0 {
			label := anchor.Label
			if len(label) == 0 {
				label = runtime.FuncForPC(anchor.PC).Name()
			}
			PrintTimeElapsed(label, totalElapsed, anchor, parent)
			totalCycles += anchor.ElapsedCyclesExclusive
			totalHits += anchor.HitCount
		}
	}
	if totalHits > 0 {
		var curr, parent Anchor

		curr.ElapsedCyclesExclusive = totalCycles
		curr.HitCount = totalHits

		PrintTimeElapsed("= Grand total", totalElapsed, &curr, &parent)
	}
}
