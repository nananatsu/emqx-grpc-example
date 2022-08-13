package adapter

import (
	"io"
)

type BytesReader struct {
	content []byte
}

func (r *BytesReader) Read(p []byte) (n int, err error) {

	readSize := len(p)
	remainSize := len(r.content)

	if remainSize > 0 {
		if remainSize < readSize {
			readSize = remainSize
		}
		for i := 0; i < readSize; i++ {
			p[i] = r.content[i]
		}
	} else {
		return 0, io.EOF
	}
	r.content = r.content[readSize:]
	return readSize, nil
}

func (r *BytesReader) Close() error {
	r.content = nil
	return nil
}
