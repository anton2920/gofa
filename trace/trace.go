//go:build trace

package trace

import (
	"fmt"
	"runtime"
	"slices"
	"unsafe"

	"github.com/anton2920/gofa/intel"
)

type Profiler struct {
	Anchors       [64 * 1024]Anchor
	CurrentParent int32
	LastAnchor    int32

	StartCycles intel.Cycles
	EndCycles   intel.Cycles
}

var GlobalProfiler Profiler

func BeginProfile() {
	clear(GlobalProfiler.Anchors[:])
	GlobalProfiler.StartCycles = intel.RDTSC()
}

func AnchorIndexForPC(pc uintptr) int32 {
	return int32(int(pc) & (len(GlobalProfiler.Anchors) - 1))
}

//go:nosplit
func Start(label string) Block {
	intel.LFENCE()
	pc := *((*uintptr)(unsafe.Add(unsafe.Pointer(&label), -8)))
	return start(pc, label)
}

//go:nosplit
func start(pc uintptr, label string) Block {
	var b Block

	index := AnchorIndexForPC(pc)
	b.ParentIndex = GlobalProfiler.CurrentParent
	GlobalProfiler.CurrentParent = index

	b.AnchorIndex = index
	b.Label = label
	b.PC = pc

	anchor := &GlobalProfiler.Anchors[b.AnchorIndex]
	b.OldElapsedCyclesInclusive = anchor.ElapsedCyclesInclusive

	b.StartCycles = intel.RDTSC()
	return b
}

//go:nosplit
func End(b Block) {
	elapsed := intel.RDTSC() - b.StartCycles
	anchor := &GlobalProfiler.Anchors[b.AnchorIndex]
	parent := &GlobalProfiler.Anchors[b.ParentIndex]
	GlobalProfiler.CurrentParent = b.ParentIndex

	parent.ElapsedCyclesExclusive -= elapsed

	anchor.ElapsedCyclesInclusive = b.OldElapsedCyclesInclusive + elapsed
	anchor.ElapsedCyclesExclusive += elapsed
	anchor.HitCount++

	if (anchor.PC > 0) && (anchor.PC != b.PC) {
		/* NOTE(anton2920): LOL. */
		panic("PC collision, you're f**cked!")
	}
	anchor.PC = b.PC
	anchor.Label = b.Label
}

func CyclesToMsec(c intel.Cycles) float32 {
	return 1000 * float32(c) / float32(intel.CpuHz)
}

func PrintTimeElapsed(label string, totalElapsed, elapsedCyclesExclusive, elapsedCyclesInclusive intel.Cycles, hitCount int) {
	percent := 100 * (float32(elapsedCyclesExclusive) / float32(totalElapsed))

	fmt.Printf("[gofa/trace]: \t %s[%d]: flat: [%.4fms %.2f%%]", label, hitCount, CyclesToMsec(elapsedCyclesExclusive), percent)
	if elapsedCyclesInclusive > elapsedCyclesExclusive {
		percentWidthChildren := 100 * (float32(elapsedCyclesInclusive) / float32(totalElapsed))
		fmt.Printf(", cum [%.4fms %.2f%%]", CyclesToMsec(elapsedCyclesInclusive), percentWidthChildren)
	}
	fmt.Printf("\n")
}

func EndAndPrintProfile() {
	GlobalProfiler.EndCycles = intel.RDTSC()
	totalElapsed := GlobalProfiler.EndCycles - GlobalProfiler.StartCycles

	var totalCycles intel.Cycles
	var totalHits int

	fmt.Printf("[gofa/trace]: Total time: %.4fms\n", CyclesToMsec(totalElapsed))

	/* NOTE(anton2920): sort by flat time in descending order. */
	slices.SortFunc(GlobalProfiler.Anchors[1:], func(a, b Anchor) int {
		if a.ElapsedCyclesExclusive > b.ElapsedCyclesExclusive {
			return -1
		} else {
			return 1
		}
	})

	for i := 1; i < len(GlobalProfiler.Anchors); i++ {
		anchor := &GlobalProfiler.Anchors[i]

		if anchor.ElapsedCyclesExclusive > 0 {
			label := anchor.Label
			if len(label) == 0 {
				label = runtime.FuncForPC(anchor.PC).Name()
			}
			PrintTimeElapsed(label, totalElapsed, anchor.ElapsedCyclesExclusive, anchor.ElapsedCyclesInclusive, anchor.HitCount)
			totalCycles += anchor.ElapsedCyclesExclusive
			totalHits += anchor.HitCount
		}
	}
	PrintTimeElapsed("= Grand total", totalElapsed, totalCycles, 0, totalHits)
}
