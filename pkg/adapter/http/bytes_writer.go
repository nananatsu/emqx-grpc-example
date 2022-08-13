package adapter

import "bytes"

type BytesWriter struct {
	content bytes.Buffer
}

func (w *BytesWriter) Write(p []byte) (n int, err error) {
	return w.content.Write(p)
}
