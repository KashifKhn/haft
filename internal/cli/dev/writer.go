package dev

import (
	"bytes"
	"io"
)

type crlfWriter struct {
	w io.Writer
}

func newCRLFWriter(w io.Writer) *crlfWriter {
	return &crlfWriter{w: w}
}

func (c *crlfWriter) Write(p []byte) (n int, err error) {
	written := 0
	for len(p) > 0 {
		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			n, err := c.w.Write(p)
			written += n
			return written, err
		}

		if idx > 0 && p[idx-1] == '\r' {
			n, err := c.w.Write(p[:idx+1])
			written += n
			if err != nil {
				return written, err
			}
			p = p[idx+1:]
			continue
		}

		if idx > 0 {
			n, err := c.w.Write(p[:idx])
			written += n
			if err != nil {
				return written, err
			}
		}

		_, err := c.w.Write([]byte("\r\n"))
		if err != nil {
			return written, err
		}
		written += 1
		p = p[idx+1:]
	}
	return written, nil
}
