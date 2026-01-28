package resp

import (
	"bufio"
	"bytes"
	"testing"
)

func TestWriteSimpleString(t *testing.T) {
	var buf bytes.Buffer
	v := Value{
		Type:  SimpleString,
		Bytes: []byte("HELLO"),
	}
	w := bufio.NewWriter(&buf)
	writer := NewWriter(w, nil)
	err := writer.WriteValue(&v)
	if err != nil {
		t.Error(err)
	}
}
