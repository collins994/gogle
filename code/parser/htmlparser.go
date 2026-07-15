/*
	errors I'm worried About:
		a < in a character string
*/
package parser

import (
	"errors"
	"os"
	"unicode"
)

type ParserEventType int

/*
	consider:
	<html>
		<p> my name is collins </p>
		<a href="example.com"> go to example.com </a>
	</html>
*/
const (
	ParserEventTypeTextNode       ParserEventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	ParserEventTypeEndDocument                           // triggered right after the > symbol of </html> is read
	ParserEventTypeOpeningTag                            // triggered right after the h in <html> is read (and other tags too)
	ParserEventTypeClosingTag                            // triggered right after the / in </html> is read (and other tags too)
	ParserEventTypeAttributeKey                          // triggered right after the h in <a href="example.com"> is read
	ParserEventTypeAttributeValue                        // triggered right after the = in <a href="example.com"> is read
	ParserEventTypeComment                               //
	ParserEventTypeUnknown                               // used to signal that an event.type is not set
)

var (
	EventErrorInvalidWhiteSpace = errors.New("Invalid whitespace")
	EventErrorInvalidTag        = errors.New("Invalid tag")
	EventErrorInvalidEndOfFile  = errors.New("Invalid end of file")
	EventErrorInvalidNewLine    = errors.New("Invalid new line")
)

type ParserEvent struct {
	Type        ParserEventType
	EventBuffer []byte
	EventError  error
}

type parserState byte

const (
	parserStateOpeningTag parserState = iota
	parserStateUnknown
	parserStateAttributeKey
	parserStateAttributeValue
)

func ParseHTMLFile(file *os.File) func(*ParserEvent) {
	var nextbyte byte
	var numberOfSpacesSkipped int
	var reader = newFileReader(file)
	var currentState parserState = parserStateUnknown

	return func(event *ParserEvent) {
		event.EventBuffer = event.EventBuffer[:0]
		event.EventError = nil

		nextbyte, _ = reader.read(skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
		if nextbyte == 0 {
			event.Type = ParserEventTypeEndDocument
			return
		}
		/**/
		if currentState == parserStateOpeningTag {
			goto attributeKey
		}

		if currentState == parserStateAttributeKey {
			if nextbyte == '=' {
				reader.read(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
				goto attributeValue
			}
			goto attributeKey
		}

		if currentState == parserStateAttributeValue {
			goto attributeKey
		}

		/* tag */
		if nextbyte == '<' {
			reader.read(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			nextbyte, numberOfSpacesSkipped = reader.read(skipWhiteSpace, dontConsumeFirstCharacter)
			if numberOfSpacesSkipped > 0 {
				event.EventError = EventErrorInvalidWhiteSpace
				return
			}

			// the second byte after < has to be a letter,  /,  !
			// discard comments and document declaration
			if nextbyte == '!' {
				event.Type = ParserEventTypeComment
				reader.read(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
				secondbyte, _ := reader.read(skipWhiteSpace, dontConsumeFirstCharacter)
				for {
					nextbyte, _ = reader.read(skipWhiteSpace, consumeFirstCharacter)
					if nextbyte == '-' {
						nextbyte, _ = reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
						if nextbyte == '-' {
							nextbyte, _ = reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
							if nextbyte == '>' {
								break
							}
						}
					}
					if nextbyte == '>' && secondbyte != '-' {
						break
					}
					if nextbyte == 0 {
						event.EventError = EventErrorInvalidEndOfFile
						break
					}
				}
				return
			}

			if nextbyte == '/' { // a closing tag
				reader.read(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
				goto closingTag
			}

			if unicode.IsLetter(rune(nextbyte)) { // an openingTag
				goto openingTag
			}

			event.EventError = EventErrorInvalidTag
			return
		}
		/* tag */
		goto textNode

	openingTag:
		{
			currentState = parserStateOpeningTag
			event.Type = ParserEventTypeOpeningTag
			// read upto the first space, or >
			for {
				nextbyte, _ = reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				if nextbyte == '>' {
					reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					currentState = parserStateUnknown
					break
				}
				if unicode.IsSpace(rune(nextbyte)) {
					reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					nextbyte, _ = reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
						break
					}
					break
				}

				b, _ := reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	attributeKey:
		{
			currentState = parserStateAttributeKey
			event.Type = ParserEventTypeAttributeKey
			// read up until a space or '=' or >
			for {
				nextbyte, _ := reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					return
				}
				if nextbyte == '=' || unicode.IsSpace(rune(nextbyte)) {
					currentState = parserStateAttributeKey
					nextbyte, _ := reader.read(skipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
					}
					return
				}
				if nextbyte == '>' {
					reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					currentState = parserStateUnknown                      // end of the openingTag
					return
				}
				b, _ := reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	attributeValue:
		{
			currentState = parserStateAttributeValue
			event.Type = ParserEventTypeAttributeValue
			// read up until '>' || space (outside quotes)
			// discard quotes if any
			firstbyte, _ := reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
			if firstbyte == '\'' || firstbyte == '"' {
				reader.read(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			}
			for {
				nextbyte, _ := reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 || nextbyte == '>' {
					currentState = parserStateUnknown
					break
				}
				if unicode.IsSpace(rune(nextbyte)) && (firstbyte != '"' && firstbyte != '\'') {
					currentState = parserStateAttributeValue
					break
				}
				if (nextbyte == '\'' && firstbyte == '\'') || (nextbyte == '"' && firstbyte == '"') {
					reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					nextbyte, _ := reader.read(skipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
					} else {
						currentState = parserStateAttributeValue
					}
					break
				}
				b, _ := reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	closingTag:
		{
			event.Type = ParserEventTypeClosingTag
			// read up to >
			for {
				nextbyte, numberOfSpacesSkipped = reader.read(skipWhiteSpace, dontConsumeFirstCharacter)
				if numberOfSpacesSkipped > 0 {
					event.EventError = EventErrorInvalidWhiteSpace
				}
				if nextbyte == '>' {
					reader.read(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					break
				}
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				b, _ := reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	textNode:
		{
			event.Type = ParserEventTypeTextNode
			// read up to space || <
			for {
				nextbyte, _ := reader.read(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				if unicode.IsSpace(rune(nextbyte)) || nextbyte == '<' {
					break
				}
				b, _ := reader.read(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

		return
	}
}

/*
if skipWhiteSpace is true,  nextByte will skip any tabs, spaces, newlines, carriage returns to get to the firs non whitespace character
if skipWhiteSpace is false, nextByte will return the first byte it encounters, whether a whitespace character or not
if consumeFirstCharacter is true, nextByte will consume the first character it encounters before returning it, meaning subsequent calls to nextByte will return the bytes after the consumed one
if consumeFirstCharacter is false, nextByte will peek and return the byte, meaning a subsequent call will return the same byte as the one before
at end of file, nextByte will return 0
*/
const (
	skipWhiteSpace            = true
	dontSkipWhiteSpace        = false
	consumeFirstCharacter     = true
	dontConsumeFirstCharacter = false
)
