package collector

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strconv"
	"unicode/utf16"

	"github.com/tdakkota/win32metadata/types"
)

func readValue(cst types.Constant) (interface{}, error) {
	order := binary.LittleEndian
	switch cst.Type {
	case types.ELEMENT_TYPE_I1:
		if len(cst.Value) < 1 {
			return nil, io.ErrUnexpectedEOF
		}
		return int8(cst.Value[0]), nil
	case types.ELEMENT_TYPE_U1:
		if len(cst.Value) < 1 {
			return nil, io.ErrUnexpectedEOF
		}
		return uint8(cst.Value[0]), nil
	case types.ELEMENT_TYPE_I2:
		if len(cst.Value) < 2 {
			return nil, io.ErrUnexpectedEOF
		}
		return int16(order.Uint16(cst.Value)), nil
	case types.ELEMENT_TYPE_U2:
		if len(cst.Value) < 2 {
			return nil, io.ErrUnexpectedEOF
		}
		return uint16(order.Uint16(cst.Value)), nil
	case types.ELEMENT_TYPE_I4:
		if len(cst.Value) < 4 {
			return nil, io.ErrUnexpectedEOF
		}
		return int32(order.Uint32(cst.Value)), nil
	case types.ELEMENT_TYPE_U4:
		if len(cst.Value) < 4 {
			return nil, io.ErrUnexpectedEOF
		}
		return uint32(order.Uint32(cst.Value)), nil
	case types.ELEMENT_TYPE_I8:
		if len(cst.Value) < 8 {
			return nil, io.ErrUnexpectedEOF
		}
		return int64(order.Uint64(cst.Value)), nil
	case types.ELEMENT_TYPE_U8:
		if len(cst.Value) < 8 {
			return nil, io.ErrUnexpectedEOF
		}
		return uint64(order.Uint64(cst.Value)), nil
	case types.ELEMENT_TYPE_R4:
		if len(cst.Value) < 4 {
			return nil, io.ErrUnexpectedEOF
		}
		return math.Float32frombits(order.Uint32(cst.Value)), nil
	case types.ELEMENT_TYPE_R8:
		if len(cst.Value) < 8 {
			return nil, io.ErrUnexpectedEOF
		}
		return math.Float64frombits(order.Uint64(cst.Value)), nil
	case types.ELEMENT_TYPE_STRING:
		remainder := len(cst.Value) % 2
		length := len(cst.Value) - remainder

		buf, n := make([]uint16, length), 0
		for i := 0; i < length; i += 2 {
			buf[n] = order.Uint16(cst.Value[i:])
			n++
		}

		return strconv.Quote(string(utf16.Decode(buf))), nil
	default:
		return nil, fmt.Errorf("unexpected constant type %v", cst.Type)
	}
}
