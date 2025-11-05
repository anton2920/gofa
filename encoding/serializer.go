package encoding

type Serializer interface {
	Begin()
	End()

	ObjectBegin()
	Key(string)
	ObjectEnd()

	ArrayBegin()
	ArrayEnd()

	Bool(bool)
	Int32(int32)
	Uint32(uint32)
	Int64(int64)
	String(string)
}
