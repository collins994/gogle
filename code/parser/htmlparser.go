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

	functionStart:
		nextbyte, _ = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
		if nextbyte == 0 {
			event.Type = ParserEventTypeEndDocument
			return
		}

		/*these states are used to make sure we don't break out of an opening tag in between calls untill we've read the entire tag*/
		if currentState == parserStateOpeningTag {
			goto attributeKey
		}

		if currentState == parserStateAttributeKey {
			if nextbyte == '=' {
				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
				goto attributeValue
			}
			goto attributeKey
		}

		if currentState == parserStateAttributeValue {
			goto attributeKey
		}

		/* tag */
		if nextbyte == '<' {
			reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
			if numberOfSpacesSkipped > 0 {
				event.EventError = EventErrorInvalidWhiteSpace
				return
			}

			// the second byte after < has to be a letter,  /,  !
			// discard comments and document declaration
			if nextbyte == '!' {
				event.Type = ParserEventTypeComment
				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
				secondbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
				for {
					nextbyte, _ = reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter)
					if nextbyte == '-' {
						nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
						if nextbyte == '-' {
							nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
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
				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
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
				nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				if nextbyte == '>' {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					currentState = parserStateUnknown
					break
				}
				if unicode.IsSpace(rune(nextbyte)) {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
						break
					}
					break
				}

				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
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
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					return
				}
				if nextbyte == '=' || unicode.IsSpace(rune(nextbyte)) {
					currentState = parserStateAttributeKey
					nextbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
					}
					return
				}
				if nextbyte == '>' {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					currentState = parserStateUnknown                               // end of the openingTag
					return
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
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
			firstbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
			if firstbyte == '\'' || firstbyte == '"' {
				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			}
			for {
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 || nextbyte == '>' {
					currentState = parserStateUnknown
					break
				}
				if unicode.IsSpace(rune(nextbyte)) && (firstbyte != '"' && firstbyte != '\'') {
					currentState = parserStateAttributeValue
					break
				}
				if (nextbyte == '\'' && firstbyte == '\'') || (nextbyte == '"' && firstbyte == '"') {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					nextbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
					if nextbyte == '>' {
						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
						currentState = parserStateUnknown
					} else {
						currentState = parserStateAttributeValue
					}
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	closingTag:
		{
			event.Type = ParserEventTypeClosingTag
			// read up to >
			for {
				nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
				if numberOfSpacesSkipped > 0 {
					event.EventError = EventErrorInvalidWhiteSpace
				}
				if nextbyte == '>' {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					break
				}
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				event.EventBuffer = append(event.EventBuffer, b)
			}
			return
		}

	textNode:
		{
			event.Type = ParserEventTypeTextNode
			// read up to space || <
			for {
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.EventError = EventErrorInvalidEndOfFile
					break
				}
				if unicode.IsSpace(rune(nextbyte)) || nextbyte == '<' {
					// EDGE CASE: skip returning an empty textNode
					if !(len(event.EventBuffer) > 0){ 
						goto functionStart;
					}
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				// ignore punctuations
				if !ignorePunctuation[b] {
					event.EventBuffer = append(event.EventBuffer, b)
				}
			}
			return
		}

		return
	}
}

var ignorePunctuation = map[byte]bool{
	// brackets and braces
	'(': true, ')': true, '[': true, ']': true, '{': true, '}': true,

	// sentence punctuation
	',': true, '.': true, ':': true, ';': true, '?': true, '!': true,

	// quotes and apostrophes
	'\'': true, '"': true, '`': true,

	// mathematical and slashes
	'+': true, '-': true, '*': true, '/': true, '\\': true, '=': true, '<': true, '>': true,

	// financial and commercial
	'$': true, '%': true, '@': true, '&': true, '_': false,

	// connectors and miscellaneous
	'#': true, '^': true, '~': true, '|': true,
}
