package os

type Handle int32

const (
	Stdin = Handle(iota)
	Stdout
	Stderr
)
