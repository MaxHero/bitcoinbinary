package bitcoinbinary

import (
	"io"
)

// This file implements "varstr" encoding of byte slices.
// The encoding corresponds to bitcoin varstr format
// https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_string:
// Field Size	Description	Data type	Comments
// 1+        	length     	varint   	Length of the string
// ?         	string     	[]byte   	The string itself (can be empty)

// PutVarstr encodes a str into buf and returns the number of bytes written.
// If the buffer is too small, PutVarstr will panic.
func PutVarstr(buf []byte, str []byte) int {
	written := PutUvarint(buf, uint64(len(str)))
	if len(buf) < written+len(str) {
		panic("bitcoinbinary: not enough buffer length")
	}
	copy(buf[written:], str)
	return written + len(str)
}

// Varstr decodes a str from buf and returns that value and the
// number of bytes read (> 0). If buf too small, the value is nil and the number
// of bytes n is 0
//
func Varstr(buf []byte) ([]byte, int) {
	l, n := Uvarint(buf)
	totalN := n + int(l)
	if n == 0 || len(buf) < totalN {
		return nil, 0
	}

	str := make([]byte, l)
	copy(str, buf[n:])
	return str, totalN
}

// ReadVarstr reads an encoded str from r and returns it as a []byte.
func ReadVarstr(r io.ByteReader) ([]byte, error) {
	l, err := ReadUvarint(r)
	if err != nil {
		return nil, err
	}

	str := make([]byte, l)
	for i := uint64(0); i < l; i++ {
		str[i], err = r.ReadByte()
		if err != nil {
			return nil, err
		}
	}
	return str, nil
}
