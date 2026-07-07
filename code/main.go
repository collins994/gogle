package main

import (
	"github.com/collins994/gogle/code/sax"
)

func main() {
	sax.ParseHTMLFile("sample.html", func(event *sax.Event, err error){
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
