package wire

type ValueType byte

const (
	ValueTypeNone = ValueType(iota)
	ValueTypeByte
	ValueTypeInt32
	ValueTypeString
	ValueTypeSlice
)

var SerialType2String = [...]string{
	ValueTypeByte:   "byte",
	ValueTypeInt32:  "int32",
	ValueTypeString: "string",
	ValueTypeSlice:  "[]T",
}
