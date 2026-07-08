package main

import (
	"github.com/collins994/gogle/code/sax"
	"log"
)

func main() {
	// var filename = "build/sample.html";
	var filename = "sample.html";
	sax.ParseHTMLFile(filename, func(event *sax.Event, err error){
		if err != nil {
			log.Fatalf("%v", err);
		}

		switch(event.Type) {
			case sax.EventTypeStartDocument: println("start parsing", filename);
			case sax.EventTypeEndDocument: println("finish parsing", filename);
			case sax.EventTypeOpeningTag: println("start tag: ", string(event.Tag));
			case sax.EventTypeClosingTag: println("closing tag: ", string(event.Tag));
			case sax.EventTypeAttribute: println("Attribute, key: ", string(event.Attribute.Key), " value: ", string(event.Attribute.Value));
			case sax.EventTypeTextNode: println("Text node: ", string(event.Text));
			case sax.EventTypeUnknown: println("Uknown Event");
		}
	})
}

// func main() {
// 	var b bytes.Buffer;
// }
// 
// func main2() {
// 	var filename = "sample.html";
// 
// 	sax.ParseFileHTML(filename,func(parserEvent *sax.Event, parserError error){
// 		if parserError != nil {
// 			if errors.Is(parserError, sax.ErrorUnexpectedEndOfFile){
// 				fmt.Printf("[ERROR]: %v\n", parserError);
// 			}
// 		}
// 		switch(parserEvent.Type)
// 		{
// 			case sax.EventStartDocument: { 
// 				fmt.Printf("parsing\t%v\n", filename);
// 			} 
// 
// 			case sax.EventEndDocument: { 
// 				fmt.Printf("finish parsing\t%v\n", filename);
// 			} 
// 
// 			case sax.EventStartTag: { 
// 				fmt.Printf("start tag:\t%s\n", parserEvent.Characters.String());
// 			} 
// 
// 			case sax.EventEndTag: { 
// 				fmt.Printf("end tag:\t%s\n", parserEvent.Characters.String());
// 			} 
// 
// 			case sax.EventCharacters: { 
// 				fmt.Printf("Characters:\t%s\n", parserEvent.Characters.String());
// 			}
// 		}
// 	})
// }
