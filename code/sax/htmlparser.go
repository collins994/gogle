/*
	errors I'm worried About:
		a < in a character string
*/
package sax

import (
	"errors"
	"os"
	"unicode"
)

type ParserEventType int

/*
	consider:
	<html>
		<p> my name is collins </p>
		<a href="example.com"> go to example.com </a>
	</html>
*/
const (
	ParserEventTypeTextNode       ParserEventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	ParserEventTypeEndDocument                           // triggered right after the > symbol of </html> is read
	ParserEventTypeOpeningTag                            // triggered right after the h in <html> is read (and other tags too)
	ParserEventTypeClosingTag                            // triggered right after the / in </html> is read (and other tags too)
	ParserEventTypeAttributeKey                          // triggered right after the h in <a href="example.com"> is read
	ParserEventTypeAttributeValue                        // triggered right after the = in <a href="example.com"> is read
	ParserEventTypeUnknown                               // used to signal that an event.type is not set
)

var (
	EventErrorInvalidWhiteSpace = errors.New("Invalid whitespace")
	EventErrorInvalidTag = errors.New("Invalid tag")
)

type ParserEvent struct {
	Type        ParserEventType
	EventBuffer []byte
	EventError  error
}

func ParseHTMLFile(file *os.File, strict bool, filename string) func(*ParserEvent) {
	var nextbyte byte
	var numberOfSpacesSkipped int
	var fs = fileStruct{
		file:         file,
		index:        0,
		buffer:       make([]byte, 1024),
		bufferLength: 0,
	}

	return func(event *ParserEvent) {
		event.EventBuffer = event.EventBuffer[:0]
		event.EventError = nil;

		nextbyte, _ = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
		if nextbyte == 0 {
			event.Type = ParserEventTypeEndDocument;
			return;
		}
		if nextbyte == '<' {
			nextByte(&fs, skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			nextbyte, numberOfSpacesSkipped = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
			if strict && numberOfSpacesSkipped > 0{
				event.EventError = EventErrorInvalidWhiteSpace;
				return;
			}

			// the second byte after < has to be a letter,  /,  ! 
			switch nextbyte {
			case '!':
				{
					nextByte(&fs, skipWhiteSpace, consumeFirstCharacter) // discard delimiter (!)
					// after ! we expect a letter (for a document declaration), or - (for a comment);
					nextbyte, numberOfSpacesSkipped := nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
					if strict && numberOfSpacesSkipped > 0{
						event.EventError = EventErrorInvalidWhiteSpace;
						return;
					}
					switch(nextbyte) {
						case '-': {
							nextByte(&fs, skipWhiteSpace, consumeFirstCharacter) // discard delimiter (-)
							if !strict{
								goto comment;
							}
							// after the first -, we should get a second - (preferably no spaces)
							nextbyte, numberOfSpacesSkipped := nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
							if numberOfSpacesSkipped > 0{
								event.EventError = EventErrorInvalidWhiteSpace;
								return;
							}
							if nextbyte != '-' {
								event.EventError = EventErrorInvalidTag;
								return;
							}
							nextByte(&fs, skipWhiteSpace, consumeFirstCharacter) // discard delimiter (-)
							goto comment;
						}
						default: {
							if !(unicode.IsLetter(rune(nextbyte))) {
								event.EventError = EventErrorInvalidTag;
								return;
							} else {
								goto documentDeclaration;
							}
						}
					}
				}
			case '/':
				{
					goto closingTag
				}
			default:
				{
					if !(unicode.IsLetter(rune(nextbyte))) {
						event.EventError = EventErrorInvalidTag;
						return;
					}
					goto openingTag
				}
			}
		} else {
			goto textNode
		}

	documentDeclaration:
		{
			println("documentDeclaration")
			return 
		}

	comment:
		{
			println("comment")
			// discard upto the closing sequence -->
			return 
		}

	closingTag:
		{
			println("closing tag")
			return 
		}

	openingTag:
		{
			println("Opening tag")
			return 
		}

	textNode:
		{
			println("text node")
			return 
		}

		return 
	}
}

type fileStruct struct {
	file         *os.File
	index        int
	buffer       []byte
	bufferLength int
}

/*
if skipWhiteSpace is true,  nextByte will skip any tabs, spaces, newlines, carriage returns to get to the firs non whitespace character
if skipWhiteSpace is false, nextByte will return the first byte it encounters, whether a whitespace character or not
if consumeFirstCharacter is true, nextByte will consume the first character it encounters before returning it, meaning subsequent calls to nextByte will return the bytes after the consumed one
if consumeFirstCharacter is false, nextByte will peek and return the byte, meaning a subsequent call will return the same byte as the one before
at end of file, nextByte will return 0
*/
const (
	skipWhiteSpace            = true
	dontSkipWhiteSpace        = false
	consumeFirstCharacter     = true
	dontConsumeFirstCharacter = false
)

func nextByte(fs *fileStruct, skipWhiteSpace bool, consumeFirstCharacter bool) (byte, int) {
	goto readFile
readFile:
	if fs.index >= fs.bufferLength {
		bytesRead, err := fs.file.Read(fs.buffer)
		if err != nil {
			return 0, 0
		}
		fs.bufferLength = bytesRead
		fs.index = 0
	}

	var numberOfSpacesSkipped = 0
	if skipWhiteSpace {
		var temporaryByte byte
		for {
			if fs.index >= fs.bufferLength {
				goto readFile
			}
			temporaryByte = fs.buffer[fs.index]
			if temporaryByte == ' ' || temporaryByte == '\n' || temporaryByte == '\t' || temporaryByte == '\r' {
				fs.index++
				numberOfSpacesSkipped++
			} else {
				break
			}
		}
	}

	var nextbyte = fs.buffer[fs.index]
	if consumeFirstCharacter {
		fs.index++
	}
	return nextbyte, numberOfSpacesSkipped
}

// type parser struct{
// 	Mode ParserMode
// 	File *os.File
//
// 	buffer       []byte
// 	index        uint
// 	bufferLength uint
//
// 	previousState parserState
// }

// type ParserMode int8
//
// const (
// 	ParserModeStrict    ParserMode = 0
// 	ParserModeNotStrict ParserMode = 1
// )
//
// type parserState int
//
// const (
// 	parserStateStart parserState = iota
// )
//
// type HTMLParser struct {
// 	Mode ParserMode
// 	File *os.File
//
// 	buffer       []byte
// 	index        uint
// 	bufferLength uint
//
// 	previousState parserState
// }
//
// func NewHTMLParser(filename string, Mode ParserMode) (*HTMLParser, error) {
// 	var parser = &HTMLParser{}
//
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	parser.File = file
// 	parser.Mode = Mode
// 	parser.EventBuffer = make([]byte, 1024)
// 	return parser
// }
//
// func (parser *HTMLParser) Next(event *ParserEvent) error {
// 	event.Type = ParserEventTypeUnknown
// 	event.EventBuffer = event.EventBuffer[:0]
// 	var nextbyte byte
// 	for {
// 		nextbyte = parser.nextByte(dontSkipWhiteSpace, consumeFirstCharacter)
// 		println(string(nextbyte))
// 		if nextbyte == 0 {
// 			break
// 		}
// 	}
// 	// determine the current state
// 	return nil
// }
//
// /*
// if skipWhiteSpace is true,  nextByte will skip any tabs, spaces, newlines, carriage returns to get to the firs non whitespace character
// if skipWhiteSpace is false, nextByte will return the first byte it encounters, whether a whitespace character or not
// if consumeFirstCharacter is true, nextByte will consume the first character it encounters before returning it, meaning subsequent calls to nextByte will return the bytes after the consumed one
// if consumeFirstCharacter is false, nextByte will peek and return the byte, meaning a subsequent call will return the same byte as the one before
// at end of file, nextByte will return 0
// */
// const (
// 	skipWhiteSpace            = true
// 	dontSkipWhiteSpace        = false
// 	consumeFirstCharacter     = true
// 	dontConsumeFirstCharacter = false
// )
//
// func (parser *HTMLParser) nextByte(skipWhiteSpace bool, consumeFirstCharacter bool) byte {
// 	goto readFile
// readFile:
// 	println("parser.index: ", parser.index)
// 	if parser.index >= parser.bufferLength {
// 		bytesRead, err := parser.File.Read(parser.buffer)
// 		println("parser.File: ", parser.File)
// 		println("bytesRead: ", bytesRead)
// 		if err != nil {
// 			return 0
// 		}
// 		parser.bufferLength = uint(bytesRead)
// 		parser.index = 0
// 	}
//
// 	if skipWhiteSpace {
// 		var temporaryByte byte
// 		for {
// 			if parser.index >= parser.bufferLength {
// 				goto readFile
// 			}
// 			temporaryByte = parser.buffer[parser.index]
// 			if !(temporaryByte == ' ' || temporaryByte == '\n' || temporaryByte == '\t' || temporaryByte == '\r') {
// 				break
// 			}
// 			parser.index++
// 		}
// 	}
//
// 	println("parser.bufferLength: ", parser.bufferLength)
// 	var nextbyte = parser.buffer[parser.index]
// 	if consumeFirstCharacter {
// 		parser.index++
// 	}
// 	return nextbyte
// }
//
// func (parser *HTMLParser) reportError(err string) {
// }
//
