package resp

type Value struct {
	Type    Type
	Integer int64
	Bytes   []byte
	Array   []*Value
	IsNil   bool
}
