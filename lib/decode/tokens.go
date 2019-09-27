package decode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

type token interface {
	getN(b []byte) (int, error)
}

func getToken(b []byte) (token, int) {
	if len(b) == 0 {
		return binCtl{ctlEOF}, 0
	}
	t := binary.BigEndian.Uint16(b[0:2])
	if tt, ok := binMap[t]; ok {
		return tt, 2
	}
	return binID{fmt.Sprintf("%X", t)}, 2
}

type binCtlType int

const (
	ctlBad binCtlType = iota
	ctlEOF
	ctlEquals
	ctlOpenGroup
	ctlCloseGroup
	ctlMagicNumber
)

type binCtl struct {
	typ binCtlType
}

func (x binCtl) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binCtl) getValue(b []byte) (binCtlType, int, error) {
	switch x.typ {
	case ctlBad:
		return ctlBad, 0, errors.New("unrecognized control structure")
	case ctlMagicNumber:
		if bytes.Compare(b[0:4], []byte{0x34, 0x62, 0x69, 0x6e}) == 0 {
			return ctlMagicNumber, 4, nil
		}
		return ctlBad, 0, errors.New("not EU4bin")
	default:
		return x.typ, 0, nil
	}
}

type binString struct{}

func (x binString) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binString) getValue(b []byte) (string, int, error) {
	// string = <len=little_endian_u16> <win-1252_bytes>
	// string = <len=little_endian_u16> <win-1252_bytes>
	l := binary.LittleEndian.Uint16(b[0:2])
	str, err := charmap.Windows1252.NewDecoder().Bytes(b[2 : 2+l])
	if err != nil {
		return "", 0, fmt.Errorf("error decoding string: %v", err)
	}
	return string(str), 2 + int(l), nil
}

type binInteger struct{}

func (x binInteger) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binInteger) getValue(b []byte) (int, int, error) {
	// int = <little_endian_i32>
	return int(int32(binary.LittleEndian.Uint32(b[0:4]))), 4, nil
}

type binBoolean struct{}

func (x binBoolean) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binBoolean) getValue(b []byte) (bool, int, error) {
	// bool = 0x00 | 0x01
	return b[0] != 0, 1, nil
}

type binFloat struct{}

func (x binFloat) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binFloat) getValue(b []byte) (float64, int, error) {
	// float = <little_endian_i32>, divide by 1000
	return float64(int32(binary.LittleEndian.Uint32(b[0:4]))) / 1000.0, 4, nil
}

type binFloatFive struct{}

func (x binFloatFive) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binFloatFive) getValue(b []byte) (float64, int, error) {
	// floatFive = <little_endian_i32, q16.16, mul by 2> <little_endian_i32, unknown>
	return (float64(binary.LittleEndian.Uint16(b[0:2]))/0x10000f + float64(int32(binary.LittleEndian.Uint16(b[2:4])))) * 2.0, 8, nil
}

type binID struct{ id string }

func (x binID) getN(b []byte) (int, error) {
	_, n, err := x.getValue(b)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (x binID) getValue(b []byte) (string, int, error) {
	return x.id, 0, nil
}

var binMap = map[uint16]token{
	0x0100: binCtl{ctlEquals},
	0x0300: binCtl{ctlOpenGroup},
	0x0400: binCtl{ctlCloseGroup},
	0x4555: binCtl{ctlMagicNumber},
	0x0F00: binString{},
	0x1700: binString{},
	0x0C00: binInteger{},
	0x1400: binInteger{},
	0x0E00: binBoolean{},
	0x0D00: binFloat{},
	0x6701: binFloatFive{},
	0x9001: binFloatFive{},
	0x4D28: binID{`date`},
}
