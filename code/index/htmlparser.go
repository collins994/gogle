package index

import (
	"errors"
	"os"
	"unicode"
)

/*
	consider:
	<html>
		<p> my name is collins </p>
		<a href="example.com"> go to example.com </a>
	</html>
*/
type eventType int

const (
	eventTypeTextNode       eventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	eventTypeEndDocument                     // triggered right after the > symbol of </html> is read
	eventTypeOpeningTag                      // triggered right after the h in <html> is read (and other tags too)
	eventTypeClosingTag                      // triggered right after the / in </html> is read (and other tags too)
	eventTypeAttributeKey                    // triggered right after the h in <a href="example.com"> is read
	eventTypeAttributeValue                  // triggered right after the = in <a href="example.com"> is read
	eventTypeComment                         //
	eventTypeUnknown                         // used to signal that an event.type is not set
)

var (
	eventErrorInvalidWhiteSpace = errors.New("Invalid whitespace")
	eventErrorInvalidTag        = errors.New("Invalid tag")
	eventErrorInvalidEndOfFile  = errors.New("Invalid end of file")
	eventErrorInvalidNewLine    = errors.New("Invalid new line")
)

type parserEvent struct {
	eventType   eventType
	eventBuffer []byte
	eventError  error
}

type parserState byte

const (
	parserStateOpeningTag parserState = iota
	parserStateUnknown
	parserStateAttributeKey
	parserStateAttributeValue
)

/*
	this is an implmentation of an html event-based parser
*/
// func parseHTMLFile(file *os.File) func(*parserEvent) {
// 	var nextbyte byte
// 	var numberOfSpacesSkipped int
// 	var reader = newFileReader(file)
// 	var currentState parserState = parserStateUnknown
// 
// 	return func(event *parserEvent) {
// 		event.eventBuffer = event.eventBuffer[:0]
// 		event.eventError = nil
// 
// 	functionStart:
// 		nextbyte, _ = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
// 		if nextbyte == 0 {
// 			event.eventType = eventTypeEndDocument
// 			return
// 		}
// 
// 		/*these states are used to make sure we don't break out of an opening tag in between calls untill we've read the entire tag*/
// 		if currentState == parserStateOpeningTag {
// 			goto attributeKey
// 		}
// 
// 		if currentState == parserStateAttributeKey {
// 			if nextbyte == '=' {
// 				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 				goto attributeValue
// 			}
// 			goto attributeKey
// 		}
// 
// 		if currentState == parserStateAttributeValue {
// 			goto attributeKey
// 		}
// 
// 		/* tag */
// 		if nextbyte == '<' {
// 			reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 			nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
// 			if numberOfSpacesSkipped > 0 {
// 				event.eventError = eventErrorInvalidWhiteSpace
// 				return
// 			}
// 
// 			// the second byte after < has to be a letter,  /,  !
// 			// discard comments and document declaration
// 			if nextbyte == '!' {
// 				event.eventType = eventTypeComment
// 				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 				secondbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
// 				for {
// 					nextbyte, _ = reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter)
// 					if nextbyte == '-' {
// 						nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 						if nextbyte == '-' {
// 							nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 							if nextbyte == '>' {
// 								break
// 							}
// 						}
// 					}
// 					if nextbyte == '>' && secondbyte != '-' {
// 						break
// 					}
// 					if nextbyte == 0 {
// 						event.eventError = eventErrorInvalidEndOfFile
// 						break
// 					}
// 				}
// 				return
// 			}
// 
// 			if nextbyte == '/' { // a closing tag
// 				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 				goto closingTag
// 			}
// 
// 			if unicode.IsLetter(rune(nextbyte)) { // an openingTag
// 				goto openingTag
// 			}
// 
// 			event.eventError = eventErrorInvalidTag
// 			return
// 		}
// 		/* tag */
// 		goto textNode
// 
// 	openingTag:
// 		{
// 			currentState = parserStateOpeningTag
// 			event.eventType = eventTypeOpeningTag
// 			// read upto the first space, or >
// 			for {
// 				nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 				if nextbyte == 0 {
// 					event.eventError = eventErrorInvalidEndOfFile
// 					break
// 				}
// 				if nextbyte == '>' {
// 					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 					currentState = parserStateUnknown
// 					break
// 				}
// 				if unicode.IsSpace(rune(nextbyte)) {
// 					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 					nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 					if nextbyte == '>' {
// 						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 						currentState = parserStateUnknown
// 						break
// 					}
// 					break
// 				}
// 
// 				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 				event.eventBuffer = append(event.eventBuffer, b)
// 			}
// 			return
// 		}
// 
// 	attributeKey:
// 		{
// 			currentState = parserStateAttributeKey
// 			event.eventType = eventTypeAttributeKey
// 			// read up until a space or '=' or >
// 			for {
// 				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 				if nextbyte == 0 {
// 					event.eventError = eventErrorInvalidEndOfFile
// 					return
// 				}
// 				if nextbyte == '=' || unicode.IsSpace(rune(nextbyte)) {
// 					currentState = parserStateAttributeKey
// 					nextbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
// 					if nextbyte == '>' {
// 						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 						currentState = parserStateUnknown
// 					}
// 					return
// 				}
// 				if nextbyte == '>' {
// 					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 					currentState = parserStateUnknown                               // end of the openingTag
// 					return
// 				}
// 				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 				event.eventBuffer = append(event.eventBuffer, b)
// 			}
// 			return
// 		}
// 
// 	attributeValue:
// 		{
// 			currentState = parserStateAttributeValue
// 			event.eventType = eventTypeAttributeValue
// 			// read up until '>' || space (outside quotes)
// 			// discard quotes if any
// 			firstbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 			if firstbyte == '\'' || firstbyte == '"' {
// 				reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 			}
// 			for {
// 				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 				if nextbyte == 0 || nextbyte == '>' {
// 					currentState = parserStateUnknown
// 					break
// 				}
// 				if unicode.IsSpace(rune(nextbyte)) && (firstbyte != '"' && firstbyte != '\'') {
// 					currentState = parserStateAttributeValue
// 					break
// 				}
// 				if (nextbyte == '\'' && firstbyte == '\'') || (nextbyte == '"' && firstbyte == '"') {
// 					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 					nextbyte, _ := reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
// 					if nextbyte == '>' {
// 						reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 						currentState = parserStateUnknown
// 					} else {
// 						currentState = parserStateAttributeValue
// 					}
// 					break
// 				}
// 				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 				event.eventBuffer = append(event.eventBuffer, b)
// 			}
// 			return
// 		}
// 
// 	closingTag:
// 		{
// 			event.eventType = eventTypeClosingTag
// 			// read up to >
// 			for {
// 				nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
// 				if numberOfSpacesSkipped > 0 {
// 					event.eventError = eventErrorInvalidWhiteSpace
// 				}
// 				if nextbyte == '>' {
// 					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
// 					break
// 				}
// 				if nextbyte == 0 {
// 					event.eventError = eventErrorInvalidEndOfFile
// 					break
// 				}
// 				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 				event.eventBuffer = append(event.eventBuffer, b)
// 			}
// 			return
// 		}
// 
// 	textNode:
// 		{
// 			event.eventType = eventTypeTextNode
// 			// read up to space || <
// 			for {
// 				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
// 				if nextbyte == 0 {
// 					event.eventError = eventErrorInvalidEndOfFile
// 					break
// 				}
// 				if unicode.IsSpace(rune(nextbyte)) || nextbyte == '<' {
// 					// EDGE CASE: skip returning an empty textNode
// 					if !(len(event.eventBuffer) > 0) {
// 						goto functionStart
// 					}
// 					break
// 				}
// 				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
// 				// ignore punctuations
// 				// if !ignorePunctuation[b] {
// 				// 	event.eventBuffer = append(event.eventBuffer, b)
// 				// }
// 				event.eventBuffer = append(event.eventBuffer, b)
// 			}
// 			return
// 		}
// 
// 		return
// 	}
// }












func parseHTMLFile(file *os.File) func(*parserEvent) {
	var nextbyte byte
	var numberOfSpacesSkipped int
	var reader = newFileReader(file)
	var currentState parserState = parserStateUnknown

	return func(event *parserEvent) {
		event.eventBuffer = event.eventBuffer[:0]
		event.eventError = nil

	functionStart:
		nextbyte, _ = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
		if nextbyte == 0 {
			event.eventType = eventTypeEndDocument
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
				event.eventError = eventErrorInvalidWhiteSpace
				return
			}

			// the second byte after < has to be a letter,  /,  !
			// discard comments and document declaration
			if nextbyte == '!' {
				event.eventType = eventTypeComment
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
						event.eventError = eventErrorInvalidEndOfFile
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

			event.eventError = eventErrorInvalidTag
			return
		}
		/* tag */
		goto textNode

	openingTag:
		{
			currentState = parserStateOpeningTag
			event.eventType = eventTypeOpeningTag
			// read upto the first space, or >
			for {
				nextbyte, _ = reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.eventError = eventErrorInvalidEndOfFile
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
				event.eventBuffer = append(event.eventBuffer, b)
			}
			return
		}

	attributeKey:
		{
			currentState = parserStateAttributeKey
			event.eventType = eventTypeAttributeKey
			// read up until a space or '=' or >
			for {
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.eventError = eventErrorInvalidEndOfFile
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
				event.eventBuffer = append(event.eventBuffer, b)
			}
			return
		}

	attributeValue:
		{
			currentState = parserStateAttributeValue
			event.eventType = eventTypeAttributeValue
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
				event.eventBuffer = append(event.eventBuffer, b)
			}
			return
		}

	closingTag:
		{
			event.eventType = eventTypeClosingTag
			// read up to >
			for {
				nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
				if numberOfSpacesSkipped > 0 {
					event.eventError = eventErrorInvalidWhiteSpace
				}
				if nextbyte == '>' {
					reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter) // discard delimiter
					break
				}
				if nextbyte == 0 {
					event.eventError = eventErrorInvalidEndOfFile
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				event.eventBuffer = append(event.eventBuffer, b)
			}
			return
		}

	textNode:
		{
			event.eventType = eventTypeTextNode
			// read up to space || <
			for {
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.eventError = eventErrorInvalidEndOfFile
					break
				}
				if unicode.IsSpace(rune(nextbyte)) || nextbyte == '<' {
					// EDGE CASE: skip returning an empty textNode
					if !(len(event.eventBuffer) > 0) {
						goto functionStart
					}
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				// ignore punctuations
				// if !ignorePunctuation[b] {
				// 	event.eventBuffer = append(event.eventBuffer, b)
				// }
				event.eventBuffer = append(event.eventBuffer, b)
			}
			return
		}

		return
	}
}

