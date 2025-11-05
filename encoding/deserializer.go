package encoding

type Deserializer interface {
	Begin() bool
	End() bool

	ObjectBegin() bool
	Key(*string) bool
	ObjectEnd() bool

	ArrayBegin() bool
	Next() bool
	ArrayEnd() bool

	Bool(*bool) bool
	Int32(*int32) bool
	Uint32(*uint32) bool
	Int64(*int64) bool
	String(*string) bool
}
