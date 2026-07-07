package sax

import (
	"errors"
	"os"
)

/*
	consider:
	<html>
		<p> my name is collins </p>
		<a href="example.com"> go to example.com </a>
	</html>
*/
type EventType int

const (
	EventTypeTextNode      EventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	EventTypeStartDocument                  // triggered right before the < symbol of <html> is read
	EventTypeEndDocument                    // triggered right after the > symbol of </html> is read
	EventTypeOpeningTag                     // triggered right after the h in <html> is read (and other tags too)
	EventTypeClosingTag                     // triggered right after the / in </html> is read (and other tags too)
	EventTypeAttribute                      // triggered right after the h in <a href="example.com"> is read
)

type Event struct {
	Type      EventType
	Text      []byte            // defined/changes at EventTypeTextNode
	Tag       []byte            // defined/changes at EventTypeOpeningTag, and EventTypeClosingTag
	Attribute map[string]string // defined/changes at EventTypeAttribute
}

type fileStruct struct {
	file         *os.File
	index        int
	buffer       []byte
	bufferLength int
}

var (
	InvalidFilePath = errors.New("Invalid file path")
)

const (
	parserStateStart int = iota
)

func ParseHTMLFile(filename string, callbackFunction func(*Event, error)) {
	var fs fileStruct
	var parserState = parserStateStart
	var nextbyte byte;

	if file, err := os.Open(filename); err != nil {
		callbackFunction(nil, InvalidFilePath)
		return
	} else {
		fs.file = file
		fs.index = 0
		fs.buffer = make([]byte, 1024)
	}
}

/*
if skipWhiteSpace is true,  nextByte will skip any tabs, spaces, newlines, carriage returns to get to the firs non whitespace character
if skipWhiteSpace is false, nextByte will return the first byte it encounters, whether a whitespace character or not
if consumeFirstCharacter is true, nextByte will consume the first character it encounters before returning it, meaning subsequent calls to nextByte will return the bytes after the consumed one
if consumeFirstCharacter is false, nextByte will peek and return the byte, meaning a subsequent call will return the same byte as the one before
at end of file, nextByte will return 0
*/
func nextByte(fs *fileStruct, skipWhiteSpace bool, consumeFirstCharacter bool) byte {
	goto readFile
readFile:
	if fs.index >= fs.bufferLength {
		bytesRead, err := fs.file.Read(fs.buffer)
		if err != nil {
			return 0
		}
		fs.bufferLength = bytesRead
		fs.index = 0
	}

	if skipWhiteSpace {
		var temporaryByte byte
		for {
			if fs.index >= fs.bufferLength {
				goto readFile
			}
			temporaryByte = fs.buffer[fs.index]
			if temporaryByte == ' ' || temporaryByte == '\n' || temporaryByte == '\t' || temporaryByte == '\r' {
				fs.index++
			} else {
				break
			}
		}
	}
	
	var nextbyte = fs.buffer[fs.index];
	if consumeFirstCharacter { fs.index++; }
	return nextbyte;
}

// type EventType int
//
// const (
// 	EventCharacters EventType = iota
// 	EventStartDocument
// 	EventEndDocument
// 	EventStartTag
// 	EventEndTag
// )
//
// type Event struct {
// 	Type       EventType
// 	Characters strings.Builder   // only defined at EventCharacters
// 	Attributes map[string]string // only defined at EventStartTag
// 	Tag        string            // only defined at EventStartTag, and EventEndTag
// }
//
// type htmlFile struct {
// 	file         *os.File
// 	index        int
// 	buffer       []byte
// 	bufferLength int
// }
//
// const (
// 	peekChar    = 0
// 	consumeChar = 1
// )
//
// var (
// 	ErrorUnexpectedEndOfFile = errors.New("Unexpected end of file")
// 	// ErrorUnexpectedEndOfFile = func(filename string) error {return errors.New("")}
// )
//
// func ParseFileHTML(filename string, callbackFunction func(*Event, error)) {
// 	var hfile = htmlFile{}
// 	hfile.buffer = make([]byte, 1024)
// 	var nextChar byte
// 	var event = Event{
// 		Characters: strings.Builder{},
// 		Attributes: map[string]string{},
// 	}
//
// 	if f, err := os.Open(filename); err != nil {
// 		callbackFunction(nil, err)
// 		return // there's nothing we can do :)
// 	} else {
// 		hfile.file = f
// 		defer hfile.file.Close()
// 	}
// 	callbackFunction(&Event{Type: EventStartDocument}, nil)
//
// 	// process the data in buffer,
// 	// read a character at a time , decide the state
// 	goto determineState
// determineState:
// 	{
// 		nextChar = nextCharacter(&hfile, peekChar)
// 		if nextChar == 0 {
// 			event.Type = EventEndDocument
// 			callbackFunction(&event, nil)
// 			return
// 		}
// 		// consume and skip whitespace
// 		if nextChar == ' ' || nextChar == '\n' || nextChar == '\t' || nextChar == '\r' {
// 			nextCharacter(&hfile, consumeChar)
// 			goto determineState
// 		}
// 		switch nextChar {
// 		case '<':
// 			goto parseTag
// 		default:
// 			goto parsePlainCharacters
// 		}
//
// 	}
//
// parseTag:
// 	{
// 		// event.Characters = strings.Builder{}
// 		var err = event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we are consuming a character that we have already peeked,
// 		if err != nil {
// 			fmt.Printf("[FATAL]: can't write to event.Characters.Builder\n")
// 		}
// 		nextChar = nextCharacter(&hfile, peekChar)
// 		// TODO: write code to handle and end of file before the closing symbol >
// 		if nextChar == 0 {
// 			callbackFunction(&event, ErrorUnexpectedEndOfFile)
// 			goto determineState
// 		}
//
// 		if nextChar == '/' { // parsing a closing tag
// 			event.Type = EventEndTag
// 			// read up untill the > symbol
// 			for {
// 				event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we've already peeked this character
// 				nextChar = nextCharacter(&hfile, peekChar)
// 				// TODO: write code to handle and end of file before the closing symbol >
// 				// maybe we should just call the call back function with EventEndTag
// 				if nextChar == 0 {
// 					callbackFunction(&event, fmt.Errorf("(%s) %w", filename, ErrorUnexpectedEndOfFile))
// 					goto determineState
// 				}
// 				if nextChar == '>' {
// 					event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we've already peeked this character
// 					callbackFunction(&event, nil)
// 					goto determineState
// 				}
// 			}
// 		}
//
// 		// parsing a comment
// 		if nextChar == '!' {
// 			// read up untill the > symbol
// 			for {
// 				nextCharacter(&hfile, consumeChar)
// 				nextChar = nextCharacter(&hfile, peekChar)
// 				if nextChar == '>' {
// 					nextCharacter(&hfile, consumeChar)
// 					goto determineState
// 				}
// 			}
// 		}
//
// 		// parsing an opening tag
// 		event.Type = EventStartTag
// 		// event.Attributes = map[string]string{};
// 		// read up untill the > symbol
// 		for {
// 			event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we've already peeked this character
// 			nextChar = nextCharacter(&hfile, peekChar)
// 			if nextChar == 0 {
// 				callbackFunction(&event, fmt.Errorf("(%s) %w", filename, ErrorUnexpectedEndOfFile))
// 				goto determineState
// 			}
// 			if nextChar == '>' {
// 				event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we've already peeked this character
// 				callbackFunction(&event, nil)
// 				goto determineState
// 			}
// 		}
// 	}
//
// parsePlainCharacters:
// 	{
// 		event.Type = EventCharacters
// 		// event.Characters = strings.Builder{};
// 		// read until an < symbol
// 		for {
// 			event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)) // we've already peeked this character
// 			nextChar = nextCharacter(&hfile, peekChar)
// 			// if we get an eof before the next tag opening symbol <
// 			if nextChar == 0 {
// 				callbackFunction(&event, fmt.Errorf("(%s) %w", filename, ErrorUnexpectedEndOfFile))
// 				goto determineState
// 			}
// 			if nextChar == '<' {
// 				callbackFunction(&event, nil)
// 				goto determineState
// 			}
// 		}
// 		goto determineState
// 	}
// }
//
// // reads and returns the next character in the file
// // returns 0 at the end of the file
// func nextCharacter(file *htmlFile, peekOrConsumeChar int) byte {
// 	if file.index >= file.bufferLength {
// 		bytesRead, err := file.file.Read(file.buffer) // read a kilobyte into buffer
// 		if err != nil {
// 			return 0
// 		}
// 		file.bufferLength = bytesRead
// 		file.index = 0
// 	}
//
// 	var nextChar = file.buffer[file.index]
// 	if peekOrConsumeChar == consumeChar {
// 		file.index++
// 	}
// 	return nextChar
// }
