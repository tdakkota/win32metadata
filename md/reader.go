package md

import (
	"encoding/binary"
	"io"
)

type reader struct {
	r   io.Reader
	err error
}

func (r reader) Err() error {
	return r.err
}

func (r *reader) Read(dst interface{}) bool {
	r.err = binary.Read(r.r, binary.LittleEndian, dst)
	return r.err == nil
}
