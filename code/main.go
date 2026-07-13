package main

import (
	"fmt"
	"github.com/collins994/gogle/code/parser"
	"os"
)

func main() {
	file, _ := os.Open("sample.html");
	defer file.Close()
	var event = sax.ParserEvent{ }
	event.Type = sax.ParserEventTypeUnknown;
	event.EventBuffer = make([]byte, 1024);

	var Next = sax.ParseHTMLFile(file, "sample.html");
	for {
		Next(&event);
		if event.Type == sax.ParserEventTypeEndDocument {
			println("done parsing");
			break;
		}
		if event.EventError != nil {
			fmt.Printf("EventError: %s, ", event.EventError);
		}
		println(event.Type, string(event.EventBuffer));
	}
}
