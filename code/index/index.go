package index

import (
	"fmt"
	"os"
	// "sort"
)

func Index() func(*os.File) {
	var (
		event         = parserEvent{}
		NextEventFunc func(*parserEvent)
	)

	return func(file *os.File) {
		event.eventType = eventTypeUnknown
		event.eventBuffer = make([]byte, 1024)

		println("[INDEXING]: ", file.Name())
		NextEventFunc = parseHTMLFile(file)
		for {
			NextEventFunc(&event)
			if event.eventError != nil {
				fmt.Printf("[ERROR]: %v\n", event.eventError)
				break
			}

			if event.eventType == eventTypeEndDocument {
				break
			}

			if event.eventType == eventTypeTextNode {
				// porterStem(&event.eventBuffer)
				println(string(event.eventBuffer))
			}
		}
<<<<<<< HEAD
=======

		sort.Strings(terms)

		for index = 0; index < len(terms); index++ {
			fmt.Printf("(%s, %s)\n", terms[index], filename)
		}
>>>>>>> parent of 12bfbf1 (stemming)
	}
}
