package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrToRead = errors.New("error while reading")
)

// ReadValues TODO: Rewrite this shit, should be able to read all types
// and save in in struct (probably recursive read for arrays)
func ReadValues(reader io.Reader) (Value, error) {
	buf := make([]byte, 4096)
	read := bufio.NewReader(reader)
	buf, _, err := read.ReadLine()
	if err != nil {
		return Value{}, ErrToRead
	}
	arrLen, err := strconv.Atoi(string(buf[1]))
	if arrLen == -1 {
		return Value{Type: Array, IsNil: true}, nil
	}
	values := make([]Value, 0, arrLen)
	for i := 0; i < arrLen; i++ {
		val, err := ReadValues(reader)
		if err != nil {
			fmt.Print(err)
		}
		values = append(values, val)
	}
	return Value{Type: Array, Array: values}, nil
}
