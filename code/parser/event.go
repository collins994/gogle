package parser

import (
	"errors"
)

type eventType int

const (
	EventTypeTextNode eventType = iota
	EventTypeEndDocument
	EventTypeOpeningTag
	EventTypeClosingTag
	EventTypeComment
	EventTypeDeclaration
	EventTypeUnknown
)

// NOTE: for EventTypeOpeningTag, use the Event.GetAttributeValue(key) to extract the value of key
// for any other event, the text needed will be in event.Buffer
type Event struct {
	Type   eventType
	Buffer []byte
	Error  error
}

var (
	EventErrorInvalidEndOfFile = errors.New("Invalid end of file")
	EventErrorInvalidTag = errors.New("Invalid Tag");
)
