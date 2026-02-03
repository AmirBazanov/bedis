package writersticky

import "bufio"

type WriterSticky struct {
	W   *bufio.Writer
	Err error
}

func (rw *WriterSticky) Write(p []byte) {
	if rw.Err != nil {
		return
	}
	_, rw.Err = rw.W.Write(p)
}

func (rw *WriterSticky) WriteByte(b byte) {
	if rw.Err != nil {
		return
	}
	rw.Err = rw.W.WriteByte(b)
}

func (rw *WriterSticky) WriteString(s string) {
	if rw.Err != nil {
		return
	}
	_, rw.Err = rw.W.WriteString(s)
}
