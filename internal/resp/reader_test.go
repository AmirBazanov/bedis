package resp

import (
	"bytes"
	"testing"
)

func TestReaderBulkString(t *testing.T) {
	respCmd := "$3\r\nASD\r\n"
	buf := bytes.NewBufferString(respCmd)
	reader := New(buf, nil)
	value, err := reader.ReadValue()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if string(value.Bytes) != "ASD" {
		t.Fatalf("strings don't match: %s", value.Bytes)
	}
}
func TestReaderArrayOfBulk(t *testing.T) {
	respCmd := "*3\r\n$3\r\nSET\r\n$1\r\nA\r\n$1\r\nB\r\n"
	buf := bytes.NewBufferString(respCmd)
	reader := New(buf, nil)
	value, err := reader.ReadValue()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if len(value.Array) != 3 {
		t.Fatalf("array of values is wrong size %d", len(value.Array))
	}
	if value.Array[0].Type != BulkString || string(value.Array[0].Bytes) != "SET" {
		t.Fatalf("data don't match: %s", value.Array[0].Bytes)
	}
	if value.Array[1].Type != BulkString || string(value.Array[1].Bytes) != "A" {
		t.Fatalf("data don't match: %s", value.Array[1].Bytes)
	}

	if value.Array[2].Type != BulkString || string(value.Array[2].Bytes) != "B" {
		t.Fatalf("data don't match: %s", value.Array[2].Bytes)
	}
}

func TestReaderInteger(t *testing.T) {
	respCmd := ":-100\r\n"
	buf := bytes.NewBufferString(respCmd)
	reader := New(buf, nil)
	value, err := reader.ReadValue()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if value.Integer != -100 {
		t.Fatalf("data don't match: %d", value.Integer)
	}

}
