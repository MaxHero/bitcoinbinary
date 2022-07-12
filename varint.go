package bitcoinbinary

import (
	"encoding/binary"
	"io"
)

// This file implements "varint" encoding of 64-bit integers.
// The encoding corresponds to bitcoin varint format
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer:
//    Value     	Storage length	Format
// <  0xFD      	1             	uint8
// <= 0xFFFF    	3             	0xFD followed by the length as uint16
// <= 0xFFFFFFFF	5             	0xFE followed by the length as uint32
//              	9             	0xFF followed by the length as uint64

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	MaxVarintLen8  = 1
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 9
)

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) int {
	if x < 0xFD {
		buf[0] = byte(x)
		return MaxVarintLen8
	}
	if x <= 0xFFFF {
		buf[0] = 0xFD
		binary.LittleEndian.PutUint16(buf[1:], uint16(x))
		return MaxVarintLen16
	}
	if x <= 0xFFFFFFFF {
		buf[0] = 0xFE
		binary.LittleEndian.PutUint32(buf[1:], uint32(x))
		return MaxVarintLen32
	}
	buf[0] = 0xFF
	binary.LittleEndian.PutUint64(buf[1:], x)
	return MaxVarintLen64
}

// Uvarint decodes a uint64 from buf and returns that value and the
// number of bytes read (> 0). If buf too small, the value and the number
// of bytes n is 0
//
func Uvarint(buf []byte) (uint64, int) {
	if len(buf) == 0 {
		return 0, 0
	}

	switch buf[0] {
	case 0xFD:
		if len(buf) < MaxVarintLen16 {
			return 0, 0
		}
		return uint64(binary.LittleEndian.Uint16(buf[1:])), MaxVarintLen16
	case 0xFE:
		if len(buf) < MaxVarintLen32 {
			return 0, 0
		}
		return uint64(binary.LittleEndian.Uint32(buf[1:])), MaxVarintLen32
	case 0xFF:
		if len(buf) < MaxVarintLen64 {
			return 0, 0
		}
		return binary.LittleEndian.Uint64(buf[1:]), MaxVarintLen64
	default:
		return uint64(buf[0]), MaxVarintLen8
	}
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.ByteReader) (uint64, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	var l int

	switch b {
	case 0xFD:
		l = MaxVarintLen16 - 1
	case 0xFE:
		l = MaxVarintLen32 - 1
	case 0xFF:
		l = MaxVarintLen64 - 1
	default:
		return uint64(b), nil
	}

	var buf [8]byte
	for i := 0; i < l; i++ {
		b, err = r.ReadByte()
		if err != nil {
			return 0, err
		}
		buf[i] = b
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}
