package sax

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrorUnknown = errors.New("")
)

type EventType int

const (
	EventCharacters EventType = iota
	EventStartDocument
	EventEndDocument
	EventStartTag
	EventEndTag
)

type Event struct {
	Type       EventType
	Characters strings.Builder   // only defined at EventCharacters, and EventStartTag and EventEndTag
	Attributes map[string]string // only defined at EventStartTag
}

type htmlFile struct {
	file         *os.File
	index        int
	buffer       []byte
	bufferLength int
}

const (
	peekChar    = 0
	consumeChar = 1
)

func ParseFileHTML(filename string, callbackFunction func(*Event, error)) {
	var hfile = htmlFile{}
	hfile.buffer = make([]byte, 1024)
	var nextChar byte
	var event = Event{};

	if f, err := os.Open(filename); err != nil {
		callbackFunction(nil, err)
		return // there's nothing we can do :)
	} else {
		hfile.file = f
		defer hfile.file.Close()
	}
	callbackFunction(&Event{Type: EventStartDocument}, nil)

	// process the data in buffer,
	// read a character at a time , decide the state
	goto determineState;
	determineState:
	{
		nextChar = nextCharacter(&hfile, peekChar);
		if nextChar == 0 {
			event.Type = EventEndDocument;
			callbackFunction(&event, nil)
			return
		}
		// consume and skip whitespace
		if nextChar == ' ' || nextChar == '\n' || nextChar == '\t' || nextChar == '\r' {
			nextCharacter(&hfile, consumeChar);
			goto determineState 
		}
		switch nextChar {
		case '<': goto parseTag;
		default: goto parsePlainCharacters;
		}
	}
 
 parseTag:
 	{
  	event.Characters = strings.Builder{}
  	// we are consuming a character that we have already peeked, 
  	// so we do not have to check for eof error
  	var err = event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); 
  	if err != nil {
 			fmt.Printf("[FATAL]: can't write to event.Characters.Builder\n")
  	}
  	nextChar = nextCharacter(&hfile, peekChar);
  	// TODO: write code to handle and end of file before the closing symbol >
  	if nextChar == 0 {
  		fmt.Printf("TODO: handle this error!\n");
  	}

  	if nextChar == '/' { // parsing a closing tag
  		event.Type = EventEndTag;
  		// read up untill the > symbol
  		for {
  			event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); // we've already peeked this character
  			nextChar = nextCharacter(&hfile, peekChar);
  			// TODO: write code to handle and end of file before the closing symbol >
  			// maybe we should just call the call back function with EventEndTag 
  			if nextChar == 0 {
  				fmt.Printf("TODO: handle this error!\n");
  				goto determineState;
  			}
  			if nextChar == '>'{
  				event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); // we've already peeked this character
  				callbackFunction(&event, nil);
  				goto determineState;
  			}
  		}
  	}

		// parsing a comment
		if nextChar == '!' {
			// read up untill the > symbol
			for {
				nextCharacter(&hfile, consumeChar);
				nextChar = nextCharacter(&hfile, peekChar);
				if nextChar == '>' {
					nextCharacter(&hfile, consumeChar);
					goto determineState;
				}
			}
		}

		// parsing an opening tag
		event.Type = EventStartTag;
		// read up untill the > symbol
		for {
			event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); // we've already peeked this character
			nextChar = nextCharacter(&hfile, peekChar);
			// TODO: write code to handle and end of file before the closing symbol >
			// maybe we should just call the call back function with EventEndTag 
			if nextChar == 0 {
				fmt.Printf("TODO: handle this error!\n");
				goto determineState;
			}
			if nextChar == '>'{
				event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); // we've already peeked this character
				callbackFunction(&event, nil);
				goto determineState;
			}
		}
 	}

	parsePlainCharacters: 
	{
		event.Type = EventCharacters;
		event.Characters = strings.Builder{};
		// read until an < symbol
		for {
			event.Characters.WriteByte(nextCharacter(&hfile, consumeChar)); // we've already peeked this character
			nextChar = nextCharacter(&hfile, peekChar);
			// TODO: write code to handle and end of file before the closing symbol >
			// maybe we should just call the call back function with EventEndTag 
			if nextChar == 0 {
				fmt.Printf("TODO: handle this error!\n");
				goto determineState;
			}
			if nextChar == '<' {
				callbackFunction(&event, nil);
				goto determineState;
			}
		}
 		goto determineState
	}
}

// reads and returns the next character in the file
// returns 0 at the end of the file
func nextCharacter(file *htmlFile, peekOrConsumeChar int) byte {
	if file.index >= file.bufferLength {
		bytesRead, err := file.file.Read(file.buffer) // read a kilobyte into buffer
		if err != nil {
			return 0
		}
		file.bufferLength = bytesRead
		file.index = 0
	}

	var nextChar = file.buffer[file.index]
	if peekOrConsumeChar == consumeChar {
		file.index++
	}
	return nextChar
}
