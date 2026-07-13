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

	var Next = sax.ParseHTMLFile(file, true, "sample.html");
	var n = 4;
	for {
		if n == 0 {
			break;
		}
		n--;
		Next(&event);
		if event.EventError != nil {
			fmt.Printf("EventError: %s", event.EventError);
			return 
		}
		println(event.Type, string(event.EventBuffer));
	}
}
