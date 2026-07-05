package sax

import (
	"errors"
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

func ParseFileHTML(filename string, callbackFunction func(*Event, error)) {
	var hfile = htmlFile{}
	hfile.buffer = make([]byte, 1024)

	if f, err := os.Open(filename); err != nil {
		callbackFunction(nil, err)
		return // there's nothing we can do :)
	} else {
		hfile.file = f
		defer hfile.file.Close()
	}

	callbackFunction(&Event{Type: EventStartDocument}, nil)
	var nextChar = nextCharacter(&hfile)
	if nextChar == byte(0) {
		callbackFunction(&Event{Type: EventEndDocument}, nil)
		return
	}

	// process the data in buffer,
	// read a character at a time , decide the state
	for {
		nextChar = nextCharacter(&hfile)
		if nextChar == 0 {
			break
		}
		if nextChar == '\n' || nextChar == '\t' || nextChar == '\r' {
			continue
		}

		switch(nextChar) {
			case '<': {
				nextChar = nextCharacter(&hfile);
				if nextChar == '/' { // close tag
					println("close tag")
				} else  { // opening tag
					println("open tag")
				}
			}
		}
	}
}

// reads and returns the next character in the file
// returns 0 at the end of the file
func nextCharacter(file *htmlFile) byte {
	if file.index >= file.bufferLength {
		bytesRead, err := file.file.Read(file.buffer) // read a kilobyte into buffer
		if err != nil {
			return 0
		}
		file.bufferLength = bytesRead
		file.index = 0
	}

	var nextChar = file.buffer[file.index]
	file.index++
	return nextChar
}
