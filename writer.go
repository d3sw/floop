package floop

// bufferedWriter is a writer that buffers until a new line is received.
type BufferedWriter struct {
	buf      []byte
	callback func([]byte)
}

func NewBufferedWriter(cb func([]byte)) *BufferedWriter {
	return &BufferedWriter{callback: cb, buf: make([]byte, 0)}
}

// Write writes bytes to be parsed/analyzed.  The buffer is flushed only if the byte slice ends
// in a new line.
func (wr *BufferedWriter) Write(b []byte) (int, error) {
	// Flush and issue callback
	if b[len(b)-1] == '\n' {
		wr.callback(append(wr.buf, b...))
		wr.buf = []byte{}
	} else {
		// Append to internal buffer as there is no newline.
		wr.buf = append(wr.buf, b...)
	}

	return len(b), nil
}
