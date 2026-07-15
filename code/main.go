package main

import (
	"fmt"
	"github.com/collins994/gogle/code/parser"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	indexFolder("gl2")
}

func indexFolder(root string) {
	filepath.WalkDir("gl2", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("[ERROR]: can't access path: ", path)
			return err
		}
		// skip directories and non-html files
		if d.IsDir() {
			return nil
		}
		if ok, _ := filepath.Match("*html", filepath.Ext(path)); !ok {
			return nil
		}

		println("[INDEXING]: ", "gl2/"+filepath.Base(path))
		indexFile("gl2/" + filepath.Base(path))
		return nil
	})
}

func indexFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	var event = parser.ParserEvent{}
	event.Type = parser.ParserEventTypeUnknown
	event.EventBuffer = make([]byte, 1024)
	var Next = parser.ParseHTMLFile(file)

	for {
		event.Type = parser.ParserEventTypeUnknown
		Next(&event)
		if event.EventError != nil {
			println("[EVENT ERROR]: ", event.EventError.Error())
			return err
		}
		if event.Type == parser.ParserEventTypeEndDocument {
			break
		}
		if event.Type == parser.ParserEventTypeTextNode {
			println(string(event.EventBuffer))
		}
	}
	println()
	return nil
}
