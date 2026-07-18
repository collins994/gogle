package index

import (
	"fmt"
	"os"
	"sort"
)

func Index() func(*os.File) {
	var (
		event         = parserEvent{}
		terms         = make([]string, (0), 1024)
		index     int = 0
		filename  string
		NextEvent func(*parserEvent)
	)

	return func(file *os.File) {
		filename = file.Name()
		event.eventType = eventTypeUnknown
		event.eventBuffer = make([]byte, 1024)

		NextEvent = parseHTMLFile(file)
		for {
			NextEvent(&event)
			if event.eventError != nil {
				fmt.Printf("[ERROR]: %v\n", event.eventError)
				break
			}

			if event.eventType == eventTypeEndDocument {
				break
			}

			if event.eventType == eventTypeTextNode {
				terms = append(terms, string(event.eventBuffer))
				// terms = append(terms, string(event.eventBuffer))
				index++
			}
		}

		sort.Strings(terms)

		for index = 0; index < len(terms); index++ {
			fmt.Printf("(%s, %s)\n", terms[index], filename)
		}
	}
}
