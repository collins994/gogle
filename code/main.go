package main

import (
	"fmt"
	"github.com/collins994/gogle/code/index"
	"os"
)

func main() {
	var event = index.ParserEvent{
		Buffer: make([]rune, 250),
		Error:  nil,
	}
	// for count := 1; count <= 1; count++ {
	for count := 1; count < 352; count++ {
		file, err := os.Open(fmt.Sprintf("gl2/%d.xhtml", count))
		println("[READING]: ", file.Name())
		if err != nil {
			fmt.Printf("[ERROR]: %v\n", err)
			continue
		}
		var next = index.Parse(file)
		for {
			next(&event);
			if event.Type == index.EventTypeEndDocument {
				break;
			}
			// println(event.Type, string(event.Buffer));
		}
	}
}
