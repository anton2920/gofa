package trace_

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/trace"
)

func EndAndPrintProfile() {
	trace.EndProfile()
	profile := make([]byte, 64*1024)
	n := trace.DumpProfile(profile)
	print(bytes.AsString(profile[:n]))
	//os.WriteToFile(os.StandardErrorStream, profile[:n])
}
