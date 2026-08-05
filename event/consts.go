package event

import "github.com/anton2920/gofa/bits"

const (
	RequestRead = bits.Flags(1 << iota)
	RequestWrite
)

const (
	TriggerLevel = int(iota)
	TriggerEdge
)
