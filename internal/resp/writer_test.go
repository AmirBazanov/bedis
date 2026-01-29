package resp

import (
	"bufio"
	"bytes"
	"testing"
)

func testWrite(t *testing.T, v *Value, cmd string) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	writer := NewWriter(w, nil)
	err := writer.Value(v)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	if cmd != buf.String() {
		t.Fatalf("Read and Write dont match: %s != %s", cmd, buf.String())
	}

	t.Logf("OK: %s == %s", cmd, buf.String())
}

func TestWriteSimpleString(t *testing.T) {
	cmd := "+HELLO\r\n"
	v := Value{
		Type:  SimpleString,
		Bytes: []byte("HELLO"),
	}

	testWrite(t, &v, cmd)
}

func TestWriteInteger(t *testing.T) {
	cmd := ":10\r\n"
	v := Value{
		Type:    Integer,
		Integer: 10,
	}
	testWrite(t, &v, cmd)
}

func TestWriteSimpleError(t *testing.T) {
	cmd := "-ERROR\r\n"
	v := Value{
		Type:  SimpleError,
		Bytes: []byte("ERROR"),
	}

	testWrite(t, &v, cmd)
}

func TestWriteBulkString(t *testing.T) {
	cmd := "$5\r\nABCDE\r\n"
	v := Value{
		Type:  BulkString,
		Bytes: []byte("ABCDE"),
	}
	testWrite(t, &v, cmd)
}

func TestWriteArray(t *testing.T) {
	cmd := "*2\r\n-err\r\n:100\r\n"
	v := Value{
		Type: Array,
		Array: []*Value{
			{
				Type:  SimpleError,
				Bytes: []byte("err"),
			},
			{
				Type:    Integer,
				Integer: 100,
			},
		},
	}

	testWrite(
		t, &v, cmd)
}
