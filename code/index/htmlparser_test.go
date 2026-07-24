package index

import "testing"
import "strings"

func TestParser(t *testing.T) {
	var input = "<a href=\"collins\" blackButton=yes> <!-- comment --> hello my nigga </a><a><b>"
	var NextEvent = Parse(strings.NewReader(input))
	var event = ParserEvent{
		Buffer: make([]rune, 250),
		Error:  nil,
	}

	var read = func(eventtype eventType, expected string) {
		NextEvent(&event)
		if event.Error != nil {
			t.Logf("error: %v", event.Error)
		}
		if string(event.Buffer) != expected || (eventtype != event.Type) {
			t.Fatalf("[ERROR]: \nEXPECTED: %v, %s, \nGOT: %v, %s\n", eventtype, expected, event.Type, string(event.Buffer))
		}
	}

	read(EventTypeOpeningTag, "a")
	read(EventTypeAttributeKey, "href")
	read(EventTypeAttributeValue, "collins")
	read(EventTypeAttributeKey, "blackButton")
	read(EventTypeAttributeValue, "yes")
	read(EventTypeComment, "")
	read(EventTypeTextNode, "hello")
	read(EventTypeTextNode, "my")
	read(EventTypeTextNode, "nigga")
	read(EventTypeClosingTag, "a")
	read(EventTypeOpeningTag, "a")
	read(EventTypeOpeningTag, "b")
}
