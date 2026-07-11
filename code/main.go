package main

import (
	"fmt"
	"github.com/collins994/gogle/code/sax"
	"os"
)

func main() {
	file, _ := os.Open("sample.html");
	defer file.Close()
	var event = sax.ParserEvent{ }
	event.Type = sax.ParserEventTypeUnknown;
	event.EventBuffer = make([]byte, 1024);

	var Next = sax.ParseHTMLFile(file, true, "sample.html");
	Next(&event);
	if event.EventError != nil {
		fmt.Printf("EventError: %s", event.EventError);
		return 
	}
	print(string(event.EventBuffer));
}
