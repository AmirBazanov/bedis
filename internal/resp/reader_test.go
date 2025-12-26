package resp

import (
	"bytes"
	"testing"
)

// TODO: Is it right
func TestArrayReader(t *testing.T) {
	respCmd := "*3\r\n$3\r\nSET\r\n$1\r\nA\r\n$1\r\nB\r\n"
	buf := bytes.NewBufferString(respCmd)

	value, err := ReadValues(buf)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if len(value.Array) != 3 {
		t.Fatalf("array of values is wrong size %d", len(value.Array))
	}
	if value.Array[0].Type != BulkString && string(value.Array[0].Bytes) != "SET" {
		t.Fatalf("wrong format or value.Array of first element")
	}
	if value.Array[1].Type != BulkString && string(value.Array[1].Bytes) != "A" {
		t.Fatalf("wrong format or value.Array of first element")
	}

	if value.Array[2].Type != BulkString && string(value.Array[2].Bytes) != "B" {
		t.Fatalf("wrong format or value of first element")
	}
}
