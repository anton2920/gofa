package os

type Handle int32

const (
	StandardInputStream = Handle(iota)
	StandardOutputStream
	StandardErrorStream
)
