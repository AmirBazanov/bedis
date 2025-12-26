package resp

type Value struct {
	Type    Type
	Integer int
	Bytes   []byte
	Array   []Value
	IsNil   bool
}
