//go:build gofatrace
// +build gofatrace

package trace_

import (
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/time/time_"
	"github.com/anton2920/gofa/trace"
)

var prof trace.Profiler

func init() {
	/* NOTE(anton2920): len must be a power of two for fast modulus calculation. */
	prof.Anchors = make([]trace.Anchor, 8192)
	prof.Prefix = "[trace_]: "
}

func BeginProfile() {
	prof.BeginProfile()
}

//go:nosplit
func Begin(label string) trace.Block {
	return prof.Begin(label)
}

//go:nosplit
func End(t trace.Block) {
	t.End()
}

func EndAndPrintProfile() {
	prof.EndProfile()

	/* NOTE(anton2920): before accessing anchors, wait for possible background work to stop. */
	time_.Sleep(200 * time.Millisecond)

	profile := make([]byte, 16*1024)
	n := prof.DumpProfile(profile)
	os.WriteToFile(os.StandardErrorStream, profile[:n])
}
