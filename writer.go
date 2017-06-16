package floop

import (
	"bytes"
	"io"
)

// BufferedWriter is a writer that buffers all the data on top of a callback buffer writer
type BufferedWriter struct {
	buffer *bytes.Buffer // complete buffer if enabled
	wr     io.Writer
}

// NewBufferedWriter instantiates a new BufferedWriter.  The cb is called each time a wrie ending in
// a new line is found.  If buffer is true a data copy is kept internally which can be used later.
func NewBufferedWriter(cb func([]byte), buffer bool) *BufferedWriter {

	cbw := newCallbackWriter(cb, '\n')
	bw := &BufferedWriter{}

	// Set write function based on requested buffering strategy
	if buffer {
		bw.buffer = bytes.NewBuffer(nil)
		// buffering enabled
		bw.wr = io.MultiWriter(cbw, bw.buffer)
	} else {
		bw.wr = cbw
	}

	return bw
}

// Bytes returns all bytes written till now.  If the internal buffer is not enabled nil is
// returned
func (wr *BufferedWriter) Bytes() []byte {
	if wr.buffer == nil {
		return nil
	}

	return wr.buffer.Bytes()
}

// Write writes the byte slice using the configured writer function
func (wr *BufferedWriter) Write(b []byte) (int, error) {
	return wr.wr.Write(b)
}

func newCallbackWriter(cb func([]byte), delim byte) *callbackWriter {
	cbw := &callbackWriter{
		delim:    delim,
		callback: cb,
		buf:      make([]byte, 0),
	}
	if cbw.callback == nil {
		cbw.callback = func([]byte) {}
	}

	return cbw
}

type callbackWriter struct {
	delim    byte
	buf      []byte
	callback func([]byte)
}

// write writes bytes to be parsed/analyzed.  The buffer is flushed only if the byte slice ends
// in a new line.
func (wr *callbackWriter) Write(b []byte) (int, error) {
	// Flush and issue callback
	if b[len(b)-1] == wr.delim {
		wr.callback(append(wr.buf, b...))
		wr.buf = []byte{}
	} else {
		// Append to internal buffer as there is no newline.
		wr.buf = append(wr.buf, b...)
	}

	return len(b), nil
}
