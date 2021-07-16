package md

import (
	"encoding/binary"
	"io"
)

// reader is a simple helper for easy decoding
//
// See https://blog.golang.org/errors-are-values.
type reader struct {
	r   io.Reader
	err error
}

// Err returns the error, if any, that was encountered during last Read.
func (r reader) Err() error {
	return r.err
}

// Read reads value from stream to dst. It returns true on success, or false
// if error encountered during decoding.
func (r *reader) Read(dst interface{}) bool {
	r.err = binary.Read(r.r, binary.LittleEndian, dst)
	return r.err == nil
}
