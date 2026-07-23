package reader

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Reader struct {
	buffer io.RuneReader
}

// at end of input, Reader.readByte will set the target byte to 0
// NOTE: incase readByte encounters a multibyte character, it will try to decompose it (see the baseRune variable)
// if it can't be decomposed, readByte will ignore it
func (reader *Reader) ReadByte(target *byte) {
labelReadARune:
	character, size, err := reader.buffer.ReadRune()
	if err != nil {
		if err == io.EOF {
			*target = 0
			return
		}
		panic(err) // TODO(collins): maybe dont panic eh!
	}

	if size > 1 {
		if b, ok := baseRune[character]; ok {
			*target = b
			return
		}
		goto labelReadARune
	}

	*target = byte(character)
}

// ReadRange will read a range from the Reader.buffer, 
// ReadNormalizedRange will return the number of bytes read into the target buffer.
// at end of input, Reader.ReadRange will return a -1
// NOTE: to read up to the end of input, set stopAtDelimiter to 0
// NOTE: the delimiter is not read into the targetBuffer
// NOTE: incase ReadRange encounters a multibyte character, it will try to decompose it (see the baseRune variable)
// if it can't be decomposed, ReadRange will ignore it
// func(reader *Reader) ReadRange (target *[]byte, stopAtDelimiter byte) int {
// 	var (
// 		byteRead byte;
// 		numberOfBytesRead int = 0;
// 	)
// 
// 	character
// }



// ReadNormalizedRange will read and normalize Reader.buffer on the fly
//		- the text is lowercased,
//    - the punctuations are ignored (see ignorePunctuation variable)
// 		- any multibyte characters are decomposed to the base rune if possible, otherwise it is ignored (see the baseRune variable)
// ReadNormalizedRange will return the number of bytes read into the target buffer.
// at end of input, Reader.ReadNormalizedRange will return a -1
// NOTE: to read up to the end of input, set stopAtDelimiter to 0
// NOTE: the delimiter is not read into the targetBuffer
func (reader *Reader) ReadNormalizedRange(target *[]byte, stopAtDelimiter byte) int {
	var (
		byteRead          byte
		numberOfBytesRead int = 0
	)
labelReadARune:
	character, size, err := reader.buffer.ReadRune()
	if err != nil {
		if err == io.EOF {
			return -1
		}
		panic(err) // TODO(collins): maybe dont panic eh!
	}
	numberOfBytesRead++

	if size > 1 {
		b, ok := baseRune[character]
		if !ok {
			goto labelReadARune
		}
		byteRead = b
		goto lowercase
	}
	byteRead = byte(character)

lowercase:
	if byteRead >= 'A' && byteRead <= 'Z' {
		byteRead = byteRead + 32 // lowercase
	}

	if ignorePunctuation[byteRead] {
		goto labelReadARune
	}

	*target = append(*target, byteRead)
	if !(byteRead == stopAtDelimiter) {
		goto labelReadARune
	}

	return numberOfBytesRead
}

func NewReader[T *os.File | string](input T) *Reader {
	var r = Reader{}

	switch v := any(input).(type) {
	case string:
		{
			r.buffer = strings.NewReader(v)
		}
	case *os.File:
		{
			r.buffer = bufio.NewReader(v)
		}
	default:
		{
			return nil
		}
	}

	return &r
}

