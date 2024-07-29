package trace

import "github.com/anton2920/gofa/intel"

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
