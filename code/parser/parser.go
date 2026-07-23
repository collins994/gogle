package parser

import (
	"github.com/collins994/gogle/code/reader"
	// "os"
	// "unicode"
)

type Parser struct {
	reader reader.Reader
}

// call this function to extract the next Event from the file
// NOTE: comments and document declarations are ignored entirely,
// NOTE: for eventTypeTextNode, the text in event.Buffer is normalized
func (parser *Parser) Next(event *Event) {
}
