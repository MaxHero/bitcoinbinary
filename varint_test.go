package bitcoinbinary_test

import (
	"bytes"
	"github.com/maxhero/bitcoinbinary"
	"testing"
)

var varintCases = []struct {
	i   uint64
	buf []byte
}{
	{0, []byte{0}},
	{1, []byte{1}},
	{0xFC, []byte{0xFC}},
	{0xFD, []byte{0xFD, 0xFD, 0x00}},
	{0xFE, []byte{0xFD, 0xFE, 0x00}},
	{0xFF, []byte{0xFD, 0xFF, 0x00}},
	{0x0100, []byte{0xFD, 0x00, 0x01}},
	{1344, []byte{0xFD, 0x40, 0x05}},
	{0xFFFF, []byte{0xFD, 0xFF, 0xFF}},
	{0x010000, []byte{0xFE, 0x00, 0x00, 0x01, 0x00}},
	{0xFFFFFFFE, []byte{0xFE, 0xFE, 0xFF, 0xFF, 0xFF}},
	{0x01FFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01, 0x00, 0x00, 0x00}},
	{0xFFFFFFFFFFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestPutUvarint(t *testing.T) {
	var buf [bitcoinbinary.MaxVarintLen64]byte
	for _, c := range varintCases {
		written := bitcoinbinary.PutUvarint(buf[:], c.i)
		if written != len(c.buf) {
			t.Errorf("Case %d, failed. Expected %v len, got %v", c.i, len(c.buf), written)
		}
		if !bytes.Equal(buf[:written], c.buf) {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.i, c.buf, buf[:written])
		}
	}
}

func TestUvarint(t *testing.T) {
	for _, c := range varintCases {
		i, n := bitcoinbinary.Uvarint(c.buf)
		if n != len(c.buf) {
			t.Errorf("Case %d, failed. Expected %v len, got %v", c.i, len(c.buf), n)
		}
		if i != c.i {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.buf, c.i, i)
		}
	}
}

func TestReadUvarint(t *testing.T) {
	for _, c := range varintCases {
		r := bytes.NewReader(c.buf)

		i, err := bitcoinbinary.ReadUvarint(r)
		if err != nil {
			t.Errorf("Case %d, failed. Got error %v", c.buf, err)
		}
		if i != c.i {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.buf, c.i, i)
		}
	}
}

func TestPutUvarintReadme(t *testing.T) {
	buf := make([]byte, 5)
	written := bitcoinbinary.PutUvarint(buf, uint64(1344)) //[]byte{0xFD, 0x40, 0x05, 0x00, 0x00} in buf, 3 in written
	if written != 3 {
		t.Errorf("Something wrong with written!")
	}
	if bytes.Compare(buf, []byte{0xFD, 0x40, 0x05, 0x00, 0x00}) != 0 {
		t.Errorf("Something wrong with buf!")
	}
}

func TestUvarintReadme(t *testing.T) {
	buf := []byte{0xFD, 0x40, 0x05, 0x00, 0x00}
	i, read := bitcoinbinary.Uvarint(buf) //1344 in i, 3 in read
	if read != 3 {
		t.Errorf("Something wrong with read!")
	}
	if i != 1344 {
		t.Errorf("Something wrong wuth buf!")
	}
}
