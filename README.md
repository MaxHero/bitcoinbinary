# Bitcoin binary encoding library
Golang implementation of bitcoin's [variable integer](https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_integer) and [variable string](https://en.bitcoin.it/wiki/Protocol_documentation#Variable_length_string) binary encodings

## VarInt
### Encoding
```go
buf := make([]byte, 5)
written := bitcoinbinary.PutUvarint(buf, uint64(1344)) //[]byte{0xFD, 0x40, 0x05, 0x00, 0x00} in buf, 3 in written
```
### Decoding
```go
buf := []byte{0xFD, 0x40, 0x05, 0x00, 0x00}
i, read := bitcoinbinary.Uvarint(buf) //1344 in i, 3 in read
```

## VarStr
### Encoding
```go
buf := make([]byte, 15)
written := bitcoinbinary.PutVarstr(buf, []byte("Hello World!")) //[]byte{0x0C, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x00, 0x00} in buf, 13 in written
```
### Decoding 
```go
buf := []byte{0x0C, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x00, 0x00}
b, read := bitcoinbinary.Varstr(buf) //13 in read
str := string(b)                     //"Hello world!" in str
```