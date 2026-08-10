//go:build gofatrace
// +build gofatrace

package trace

import (
	"runtime" /* TODO(anton2920): replace with my own soring routine. */
	"unsafe"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/cpu"
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/funcs"
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
	Profiler *Profiler

	PC    uintptr /* for lazy function name resolution. */
	Label string  /* for names of non-function blocks. */

	ParentIndex int32
	AnchorIndex int32

	StartCycles               cpu.Cycles
	OldElapsedCyclesInclusive cpu.Cycles
}

type Profiler struct {
	Anchors       []Anchor
	CurrentParent int32

	StartCycles cpu.Cycles
	EndCycles   cpu.Cycles
}

const prefix = "[trace]: "

func AnchorLess(a *Anchor, b *Anchor) bool {
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

func InsertionSortAnchors(as []Anchor) {
	for i := 1; i < len(as); i++ {
		for j := i; (j > 0) && (AnchorLess(&as[j], &as[j-1])); j-- {
			as[j], as[j-1] = as[j-1], as[j]
		}
	}
}

func SortAnchors(as []Anchor) {
	InsertionSortAnchors(as)
}

//go:nosplit
func (p *Profiler) anchorIndexForPC(pc uintptr) int32 {
	mask := len(p.Anchors)/2 - 1

	start := int(pc & uintptr(mask))
	start += bools.ToInt(start == 0)
	if (pc == p.Anchors[start].PC) || (p.Anchors[start].PC == 0) {
		return int32(start)
	}

	var idx int
	for idx = start + 1; (pc != p.Anchors[idx].PC) && (p.Anchors[idx].PC != 0) && (idx != start); {
		idx = (idx + 1) & mask
		idx += bools.ToInt(idx == 0)
	}
	if idx == start {
		panic("not enough space for a new anchor")
	}

	return int32(idx)
}

//go:nosplit
func (p *Profiler) begin(pc uintptr, label string) Block {
	var b Block

	b.Profiler = p

	index := p.anchorIndexForPC(pc)
	b.ParentIndex = p.CurrentParent
	p.CurrentParent = index

	b.AnchorIndex = index
	b.Label = label
	b.PC = pc

	anchor := &p.Anchors[b.AnchorIndex]
	b.OldElapsedCyclesInclusive = anchor.ElapsedCyclesInclusive

	b.StartCycles = cpu.ReadPerformanceCounter()
	return b
}

//go:nosplit
func (p *Profiler) Begin(label string) Block {
	cpu.WaitForLoadOperationsToComplete()
	return p.begin(funcs.GetCallerPC(unsafe.Pointer(&label)), label)
}

func (p *Profiler) BeginProfile() {
	for i := 0; i < len(p.Anchors)/2; i++ {
		p.Anchors[i] = Anchor{}
	}
	p.CurrentParent = 0
	p.StartCycles = cpu.ReadPerformanceCounter()
}

func CyclesToNanoseconds(c cpu.Cycles) float64 {
	return 1000000000 * float64(c) / float64(cpu.CPUHz)
}

func CyclesToMilliseconds(c cpu.Cycles) float64 {
	return 1000 * float64(c) / float64(cpu.CPUHz)
}

//go:nosplit
func (b *Block) End() {
	elapsed := cpu.ReadPerformanceCounter() - b.StartCycles
	anchor := &b.Profiler.Anchors[b.AnchorIndex]
	parent := &b.Profiler.Anchors[b.ParentIndex]
	b.Profiler.CurrentParent = b.ParentIndex

	parent.ElapsedCyclesExclusive -= elapsed

	anchor.ElapsedCyclesInclusive = b.OldElapsedCyclesInclusive + elapsed
	anchor.ElapsedCyclesExclusive += elapsed
	anchor.HitCount++

	anchor.PC = b.PC
	anchor.Label = b.Label
	anchor.ParentIndex = b.ParentIndex
}

func (p *Profiler) dumpTimeElapsed(f *fmt.Formatter, label string, totalElapsed cpu.Cycles, curr *Anchor, parent *Anchor) {
	percentTotal := 100 * (float64(curr.ElapsedCyclesExclusive) / float64(totalElapsed))
	percentParent := 100 * (float64(curr.ElapsedCyclesExclusive) / float64(parent.ElapsedCyclesInclusive))
	f.S(prefix).S("\t ").S(label).S("[").D(curr.HitCount).S("]: flat [").Prec(4).F(curr.ElapsedCyclesExclusive.ToMilliseconds()).S("ms ").Prec(2).F(percentTotal).S("%/").Prec(2).F(percentParent).S("% ").Prec(2).F(curr.ElapsedCyclesExclusive.ToNanoseconds() / float64(curr.HitCount)).S("ns/op")

	if curr.ElapsedCyclesInclusive > curr.ElapsedCyclesExclusive {
		percentWithChildrenTotal := 100 * (float64(curr.ElapsedCyclesInclusive) / float64(totalElapsed))
		percentWithChildrenParent := 100 * (float64(curr.ElapsedCyclesInclusive) / float64(parent.ElapsedCyclesInclusive))
		f.S(", cum [").Prec(4).F(curr.ElapsedCyclesInclusive.ToMilliseconds()).S("ms ").Prec(2).F(percentWithChildrenTotal).S("%/").Prec(2).F(percentWithChildrenParent).S("% ").Prec(2).F(curr.ElapsedCyclesInclusive.ToNanoseconds() / float64(curr.HitCount)).S("ns/op")
	}

	f.Ln()
}

func (p *Profiler) EndProfile() {
	p.EndCycles = cpu.ReadPerformanceCounter()
}

func (p *Profiler) DumpProfile(buffer []byte) int {
	var f fmt.Formatter
	f.InitWithByteSlice(buffer)

	totalElapsed := p.EndCycles - p.StartCycles

	var totalCycles cpu.Cycles
	var totalHits int

	f.S(prefix).S("Total time: ").Prec(4).F(totalElapsed.ToMilliseconds()).S("ms").Ln()

	/* NOTE(anton2920): Anchor.ParentIndex uses original order, so we need to preserve it after Sort. Copy filled half to the latter half to create a backup. */
	half := len(p.Anchors) / 2
	copy(p.Anchors[half:], p.Anchors[:half])

	for i := 0; i < half; i++ {
		anchor := &p.Anchors[i]
		parent := &p.Anchors[half+int(anchor.ParentIndex)]

		if anchor.HitCount > 0 {
			label := anchor.Label
			if len(label) == 0 {
				label = runtime.FuncForPC(anchor.PC).Name()
			}
			p.dumpTimeElapsed(&f, label, totalElapsed, anchor, parent)
			totalCycles += anchor.ElapsedCyclesExclusive
			totalHits += anchor.HitCount
		}
	}
	if totalHits > 0 {
		var curr, parent Anchor

		curr.ElapsedCyclesExclusive = totalCycles
		curr.HitCount = totalHits

		p.dumpTimeElapsed(&f, "= Grand total", totalElapsed, &curr, &parent)
	}

	return len(f.String())
}
