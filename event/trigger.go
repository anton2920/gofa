package event

import "github.com/anton2920/gofa/bits"

const (
	TriggerLevel = bits.Flags(1 << iota)
	TriggerEdge
	TriggerOnce
)
