package sax

import (
	"errors"
	"os"
	"fmt"
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
	EventTypeUnknown                        // used to signal that an event.type is not set
)

type Event struct {
	Type      EventType
	Text      []byte   // defined/changes at EventTypeTextNode
	Tag       []byte   // defined/changes at EventTypeOpeningTag, and EventTypeClosingTag
	Attribute struct { // defined/changes at EventTypeAttribute
		Key      []byte
		Value    []byte
		HasValue bool
	}
}

type fileStruct struct {
	file         *os.File
	index        int
	buffer       []byte
	bufferLength int
}

var (
	ErrorInvalidFilePath       = errors.New("Invalid file path")
	ErrorInvalidAttributeValue = errors.New("Invalid attribute value")
)

var (
	skipWhiteSpace            = true
	dontSkipWhiteSpace        = false
	consumeFirstCharacter     = true
	dontConsumeFirstCharacter = false
)

func ParseHTMLFile(filename string, callbackFunction func(*Event, error)) {
	var fs fileStruct
	var event = Event{
		Type: EventTypeUnknown,
		Text: []byte{},
		Tag:  []byte{},
		// Attribute: map[string]string{},
		Attribute: struct {
			Key      []byte
			Value    []byte
			HasValue bool
		}{
			Key:      make([]byte, 1024), // an entire kilobyte for the key :)
			Value:    make([]byte, 1024), // an entire kilobyte for the value :)
			HasValue: false,
		},
	}

	if file, err := os.Open(filename); err != nil {
		callbackFunction(nil, ErrorInvalidFilePath)
		return
	} else {
		fs.file = file
		fs.index = 0
		fs.buffer = make([]byte, 1024)
	}
	event.Type = EventTypeStartDocument
	callbackFunction(&event, nil)
	event.Type = EventTypeUnknown

	goto stateStart
stateStart:
	{
		//decide the state based on the first non-whitespace byte in the stream
		var nextbyte = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
		if nextbyte == 0 { // end of file
			event.Type = EventTypeEndDocument
			callbackFunction(&event, nil)
			event.Type = EventTypeUnknown
			return
		}
		if nextbyte == '<' {
			goto stateOpeningTag;
		} else {
			goto stateCharacters;
		}

		// switch nextbyte {
		// case '<':
		// 	goto stateOpeningTag
		// default:
		// 	{
		// 		println("Unexpected character: ", string(nextbyte))
		// 	}
		// }
	}

stateOpeningTag:
	{
		nextByte(&fs, skipWhiteSpace, consumeFirstCharacter)
		var nextbyte = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter) // skip any space after the < eg  <   div>
		if nextbyte == '/' {
			nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
			goto stateClosingTag
		}
		if nextbyte == '!' {
			nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
			goto stateComment
		}
		// consider <div>, <div class="">
		/*
			TODO(collins994): read the bytes up to the first whitespace (marking the end of "div", and a possible begining of attributes) or >(marking the end of the entire tag),
		*/
		event.Tag = event.Tag[:0]
		for { // read "div"
			nextbyte = nextByte(&fs, dontSkipWhiteSpace, dontConsumeFirstCharacter)
			if nextbyte == ' ' || nextbyte == '\n' || nextbyte == '\t' || nextbyte == '\r' || nextbyte == '>' {
				event.Type = EventTypeOpeningTag
				callbackFunction(&event, nil)
				event.Type = EventTypeUnknown
				if nextbyte == '>' {
					nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter)
					goto stateStart // no attributes
				}
				goto stateAttribute
			}
			event.Tag = append(event.Tag, nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter))
		}

		goto stateStart
	}

stateClosingTag:
	{
		event.Tag = event.Tag[:0];
		// println("TODO stateClosingTag!")
		// TODO: find a way to inform the user if we get a space in the closing tag
		var nextbyte = nextByte(&fs, dontSkipWhiteSpace, dontConsumeFirstCharacter);
		if nextbyte == '/' {
			nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
			goto stateClosingTag;
		}
		for {
			nextbyte = nextByte(&fs, dontSkipWhiteSpace, dontConsumeFirstCharacter);
			if nextbyte == '>' {
				nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
				event.Type = EventTypeClosingTag;
				callbackFunction(&event, nil);
				event.Type = EventTypeUnknown;
				goto stateStart;
			}
			event.Tag = append(event.Tag, nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter));
		}
	}

stateComment:
	{
		println("TODO stateComment!")
	}

stateAttribute:
	{
		event.Attribute.Key = event.Attribute.Key[:0]
		event.Attribute.Value = event.Attribute.Value[:0]
		// consider <a    href="colins">
		// skip any whitespace before the attribute
		var nextbyte = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
		if nextbyte == '>' { // end of the tag,
			nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
			event.Type = EventTypeAttribute
			callbackFunction(&event, nil)
			event.Type = EventTypeUnknown
			goto stateStart
		}
		// we are reading an attribute key (href);
		// read untill we hit a '=' or a space or a >
		for {
			event.Attribute.Key = append(event.Attribute.Key, nextByte(&fs, skipWhiteSpace, consumeFirstCharacter))
			nextbyte = nextByte(&fs, dontSkipWhiteSpace, dontConsumeFirstCharacter)
			if nextbyte == '=' { // eg <a href="collins">
				nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter)
				goto stateAttributeValue
			}
			if nextbyte == '>' { // eg <a blackButton>
				// an attribute with no value; eg <a blackButton> TODO: find a way to merge this with if nextbyte ==
				nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
				event.Type = EventTypeAttribute
				event.Attribute.HasValue = false
				callbackFunction(&event, nil)
				event.Type = EventTypeUnknown
				goto stateStart
			}
			if nextbyte == ' ' {
				nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
				nextbyte = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == '=' { // eg <a href = "collins">
					nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter)
					goto stateAttributeValue
				}
				if nextbyte == '>' { // eg <a blackButton    >
					nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter)
					event.Type = EventTypeAttribute
					event.Attribute.HasValue = false
					callbackFunction(&event, nil)
					event.Type = EventTypeUnknown
					goto stateStart
				}
				// an attribute with no value; eg <a blackButton>
				event.Type = EventTypeAttribute
				event.Attribute.HasValue = false
				callbackFunction(&event, nil)
				event.Type = EventTypeUnknown
				goto stateAttribute
			}
		}

	stateAttributeValue:
		{
			// the first character of a value should be a '"'
			//TODO: find a way to inform the user when a value starts with any character other than '"'
			nextbyte = nextByte(&fs, skipWhiteSpace, dontConsumeFirstCharacter)
			if nextbyte == '"' {
				nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
				goto stateAttributeValue
			}
			// read up until '"' or >
			// TODO: find a way to inform the user  if we get a > before a closing '"'
			for {
				nextbyte = nextByte(&fs, dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == '>' || nextbyte == '"' {
					nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter) // discard the delimiter
					event.Type = EventTypeAttribute
					event.Attribute.HasValue = true
					if nextbyte == '>' {
						callbackFunction(&event, fmt.Errorf("%w, Missing closing \"", ErrorInvalidAttributeValue))
						event.Type = EventTypeUnknown
						goto stateStart
					} else { // nextbyte == '"'
						callbackFunction(&event, nil)
						event.Type = EventTypeUnknown
						goto stateAttribute
					}
				}
			event.Attribute.Value = append(event.Attribute.Value, nextByte(&fs, dontSkipWhiteSpace, consumeFirstCharacter))	
			}
		}
	}

stateCharacters:
	{
		println("TODO stateCharacters");
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

	var nextbyte = fs.buffer[fs.index]
	if consumeFirstCharacter {
		fs.index++
	}
	return nextbyte
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
