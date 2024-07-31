//go:build gofaprof

package prof

import (
	"fmt"
	"runtime"
	"slices"
	"unsafe"

	"github.com/anton2920/gofa/intel"
)

type Anchor struct {
	PC    uintptr /* NOTE(anton2920): for lazy function name resolution. */
	Label string  /* NOTE(anton2920): for names of non-function blocks. */

	HitCount int

	ElapsedCyclesExclusive intel.Cycles /* NOTE(anton2920): time for anchor itself. */
	ElapsedCyclesInclusive intel.Cycles /* NOTE(anton2920): time for anchor plus time of its children. */
}

type Block struct {
	PC    uintptr /* NOTE(anton2920): for lazy function name resolution. */
	Label string  /* NOTE(anton2920): for names of non-function blocks. */

	ParentIndex int32
	AnchorIndex int32

	StartCycles               intel.Cycles
	OldElapsedCyclesInclusive intel.Cycles
}

type Profiler struct {
	Anchors       [64 * 1024]Anchor
	CurrentParent int32

	StartCycles intel.Cycles
	EndCycles   intel.Cycles
}

var GlobalProfiler Profiler

func anchorIndexForPC(pc uintptr) int32 {
	return int32(int(pc) & (len(GlobalProfiler.Anchors) - 1))
}

//go:nosplit
func BeginProfile() {
	clear(GlobalProfiler.Anchors[:])
	GlobalProfiler.CurrentParent = 0
	GlobalProfiler.StartCycles = intel.RDTSC()
}

//go:nosplit
func Begin(label string) Block {
	intel.LFENCE()
	pc := *((*uintptr)(unsafe.Add(unsafe.Pointer(&label), -8)))
	return begin(pc, label)
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

func CyclesToNsec(c intel.Cycles) float64 {
	return 1_000_000_000 * float64(c) / float64(intel.CPUHz)
}

func CyclesToMsec(c intel.Cycles) float64 {
	return 1_000 * float64(c) / float64(intel.CPUHz)
}

func PrintTimeElapsed(label string, totalElapsed, elapsedCyclesExclusive, elapsedCyclesInclusive intel.Cycles, hitCount int) {
	percent := 100 * (float64(elapsedCyclesExclusive) / float64(totalElapsed))

	fmt.Printf("[gofa/prof]: \t %s[%d]: flat: [%.4fms %.2f%% %.2fns/op]", label, hitCount, CyclesToMsec(elapsedCyclesExclusive), percent, CyclesToNsec(elapsedCyclesExclusive)/float64(hitCount))
	if elapsedCyclesInclusive > elapsedCyclesExclusive {
		percentWidthChildren := 100 * (float64(elapsedCyclesInclusive) / float64(totalElapsed))
		fmt.Printf(", cum [%.4fms %.2f%% %.2fns/op]", CyclesToMsec(elapsedCyclesInclusive), percentWidthChildren, CyclesToNsec(elapsedCyclesInclusive)/float64(hitCount))
	}
	fmt.Printf("\n")
}

func EndAndPrintProfile() {
	GlobalProfiler.EndCycles = intel.RDTSC()
	totalElapsed := GlobalProfiler.EndCycles - GlobalProfiler.StartCycles

	var totalCycles intel.Cycles
	var totalHits int

	fmt.Printf("[gofa/prof]: Total time: %.4fms\n", CyclesToMsec(totalElapsed))

	slices.SortFunc(GlobalProfiler.Anchors[:], func(a, b Anchor) int {
		if (a.ElapsedCyclesInclusive > 0) && (b.ElapsedCyclesInclusive > 0) {
			if a.ElapsedCyclesInclusive < b.ElapsedCyclesInclusive {
				return 1
			} else {
				return -1
			}
		} else if a.ElapsedCyclesInclusive > 0 {
			if a.ElapsedCyclesInclusive < b.ElapsedCyclesExclusive {
				return 1
			} else {
				return -1
			}
		} else if b.ElapsedCyclesInclusive > 0 {
			if a.ElapsedCyclesExclusive < b.ElapsedCyclesInclusive {
				return 1
			} else {
				return -1
			}
		} else {
			if a.ElapsedCyclesExclusive < b.ElapsedCyclesExclusive {
				return 1
			} else {
				return -1
			}
		}
	})

	for i := 0; i < len(GlobalProfiler.Anchors); i++ {
		anchor := &GlobalProfiler.Anchors[i]

		if anchor.HitCount > 0 {
			label := anchor.Label
			if len(label) == 0 {
				label = runtime.FuncForPC(anchor.PC).Name()
			}
			PrintTimeElapsed(label, totalElapsed, anchor.ElapsedCyclesExclusive, anchor.ElapsedCyclesInclusive, anchor.HitCount)
			totalCycles += anchor.ElapsedCyclesExclusive
			totalHits += anchor.HitCount
		}
	}
	if totalHits > 0 {
		PrintTimeElapsed("= Grand total", totalElapsed, totalCycles, 0, totalHits)
	}
}
