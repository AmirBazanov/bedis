package resp

type Type byte

const (
	SimpleString Type = '+'
	SimpleError  Type = '-'
	BulkString   Type = '$'
	Integer      Type = ':'
	Array        Type = '*'
	// TODO: Other Types
)
