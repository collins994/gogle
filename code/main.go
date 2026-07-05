package main

import (
	"github.com/collins994/gogle/code/sax"
	"fmt"
)

func main() {
	var filename = "sample.html";

	sax.ParseFileHTML(filename,func(parserEvent *sax.Event, parserError error){
		if parserError != nil {
			fmt.Printf("error");
			return
		}
		switch(parserEvent.Type){
			case sax.EventStartDocument: { 
				fmt.Printf("parsing %v\n", filename);
			} 
			case sax.EventEndDocument: { 
				fmt.Printf("finish parsing %v\n", filename);
			} 

			case sax.EventStartTag: { 
				fmt.Printf("start tag: %s\n", parserEvent.Characters.String());
			} 

			case sax.EventEndTag: { 
				fmt.Printf("end tag: %s\n", parserEvent.Characters.String());
			} 

			case sax.EventCharacters: { 
				fmt.Printf("Characters: %s\n", parserEvent.Characters.String());
			}
		}
	})
}
