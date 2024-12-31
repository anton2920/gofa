//go:build gofatrace
// +build gofatrace

package trace

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"unsafe"

	"github.com/anton2920/gofa/intel"
	"github.com/anton2920/gofa/util"
)

type Anchor struct {
	PC    uintptr /* for lazy function name resolution. */
	Label string  /* for names of non-function blocks. */

	HitCount int

	ElapsedCyclesExclusive intel.Cycles /* time for anchor itself. */
	ElapsedCyclesInclusive intel.Cycles /* time for anchor plus time of its children. */
}

type Block struct {
	PC    uintptr /* for lazy function name resolution. */
	Label string  /* for names of non-function blocks. */

	ParentIndex int32
	AnchorIndex int32

	StartCycles               intel.Cycles
	OldElapsedCyclesInclusive intel.Cycles
}

type Anchors []Anchor

type Profiler struct {
	Anchors       Anchors
	CurrentParent int32

	StartCycles intel.Cycles
	EndCycles   intel.Cycles
}

var GlobalProfiler Profiler

func init() {
	/* NOTE(anton2920): len must be a power of two for fast modulus calculation. */
	GlobalProfiler.Anchors = make(Anchors, 64*1024)
}

func (as Anchors) Len() int {
	return len(as)
}

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

func (as Anchors) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func anchorIndexForPC(pc uintptr) int32 {
	idx := int32(int(pc) & (len(GlobalProfiler.Anchors) - 1))
	return idx + int32(util.Bool2Int(idx == 0))
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
func Begin(label string) Block {
	intel.LFENCE()
	return begin(util.GetCallerPC(unsafe.Pointer(&label)), label)
}

func BeginProfile() {
	for i := 0; i < len(GlobalProfiler.Anchors); i++ {
		GlobalProfiler.Anchors[i] = Anchor{}
	}
	GlobalProfiler.CurrentParent = 0
	GlobalProfiler.StartCycles = intel.RDTSC()
}

func CyclesToNsec(c intel.Cycles) float64 {
	return 1000000000 * float64(c) / float64(intel.CPUHz)
}

func CyclesToMsec(c intel.Cycles) float64 {
	return 1000 * float64(c) / float64(intel.CPUHz)
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
		println("anchor.PC", anchor.PC, "b.PC", b.PC)
		panic("PC collision, you're f**cked!")
	}
	anchor.PC = b.PC
	anchor.Label = b.Label
}

func EndAndPrintProfile() {
	GlobalProfiler.EndCycles = intel.RDTSC()
	totalElapsed := GlobalProfiler.EndCycles - GlobalProfiler.StartCycles

	var totalCycles intel.Cycles
	var totalHits int

	fmt.Fprintf(os.Stderr, "[gofa/trace]: Total time: %.4fms\n", CyclesToMsec(totalElapsed))

	sort.Sort(GlobalProfiler.Anchors)

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

func PrintTimeElapsed(label string, totalElapsed, elapsedCyclesExclusive, elapsedCyclesInclusive intel.Cycles, hitCount int) {
	percent := 100 * (float64(elapsedCyclesExclusive) / float64(totalElapsed))

	fmt.Fprintf(os.Stderr, "[gofa/trace]: \t %s[%d]: flat: [%.4fms %.2f%% %.2fns/op]", label, hitCount, CyclesToMsec(elapsedCyclesExclusive), percent, CyclesToNsec(elapsedCyclesExclusive)/float64(hitCount))
	if elapsedCyclesInclusive > elapsedCyclesExclusive {
		percentWidthChildren := 100 * (float64(elapsedCyclesInclusive) / float64(totalElapsed))
		fmt.Fprintf(os.Stderr, ", cum [%.4fms %.2f%% %.2fns/op]", CyclesToMsec(elapsedCyclesInclusive), percentWidthChildren, CyclesToNsec(elapsedCyclesInclusive)/float64(hitCount))
	}
	fmt.Fprintf(os.Stderr, "\n")
}
