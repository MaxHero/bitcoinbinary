package bitcoinbinary_test

import (
	"bytes"
	"github.com/maxhero/bitcoinbinary"
	"testing"
)

var varstrCases = []struct {
	str []byte
	buf []byte
}{
	{[]byte{}, []byte{0x00}},
	{[]byte{0xAA}, []byte{0x01, 0xAA}},
	{bytes.Repeat([]byte{0xAA}, 0xFC), append([]byte{0xFC}, bytes.Repeat([]byte{0xAA}, 0xFC)...)},
	{bytes.Repeat([]byte{0xAA}, 0xFD), append([]byte{0xFD, 0xFD, 0x00}, bytes.Repeat([]byte{0xAA}, 0xFD)...)},
	{bytes.Repeat([]byte{0xAA}, 0xFE), append([]byte{0xFD, 0xFE, 0x00}, bytes.Repeat([]byte{0xAA}, 0xFE)...)},
	{bytes.Repeat([]byte{0xAA}, 0xFF), append([]byte{0xFD, 0xFF, 0x00}, bytes.Repeat([]byte{0xAA}, 0xFF)...)},
	{bytes.Repeat([]byte{0xAA, 0xBB}, 128), append([]byte{0xFD, 0x00, 0x01}, bytes.Repeat([]byte{0xAA, 0xBB}, 128)...)},
}

func TestPutVarstr(t *testing.T) {
	var buf [1000]byte
	for _, c := range varstrCases {
		written := bitcoinbinary.PutVarstr(buf[:], c.str)
		if written != len(c.buf) {
			t.Errorf("Case %d, failed. Expected %v len, got %v", c.str, len(c.buf), written)
		}
		if !bytes.Equal(buf[:written], c.buf) {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.str, c.buf, buf[:written])
		}
	}
}

func TestVarstr(t *testing.T) {
	for _, c := range varstrCases {
		str, n := bitcoinbinary.Varstr(c.buf)
		if n != len(c.buf) {
			t.Errorf("Case %d, failed. Expected %v len, got %v", c.str, len(c.buf), n)
		}
		if bytes.Compare(str, c.str) != 0 {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.buf, c.str, str)
		}
	}
}

func TestReadVarstr(t *testing.T) {
	for _, c := range varstrCases {
		r := bytes.NewReader(c.buf)

		str, err := bitcoinbinary.ReadVarstr(r)
		if err != nil {
			t.Errorf("Case %d, failed. Got error %v", c.buf, err)
		}
		if bytes.Compare(str, c.str) != 0 {
			t.Errorf("Case %d, failed. Expected %v, got %v", c.buf, c.str, str)
		}
	}
}

func TestPutVarstrReadme(t *testing.T) {
	buf := make([]byte, 15)
	written := bitcoinbinary.PutVarstr(buf, []byte("Hello World!")) //[]byte{0x0C, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x00, 0x00} in buf, 13 in written
	if written != 13 {
		t.Errorf("Something wrong with written!")
	}
	if bytes.Compare(buf, []byte{0x0C, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x00, 0x00}) != 0 {
		t.Errorf("Something wrong with buf!")
	}
}

func TestVarstrReadme(t *testing.T) {
	buf := []byte{0x0C, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x00, 0x00}
	b, read := bitcoinbinary.Varstr(buf) //13 in read
	str := string(b)                     //"Hello world!" in str
	if read != 13 {
		t.Errorf("Something wrong with read!")
	}
	if str != "Hello World!" {
		t.Errorf("Something wrong wuth buf!")
	}
}
