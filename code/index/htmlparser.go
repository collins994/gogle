/*
type eventType int
	errors I'm worried About:
		a < in a character string
*/
package index

import (
	"bufio"
	"fmt"
	"errors"
	"io"
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
	EventTypeTextNode       eventType = iota // "my", "name", "is", "collins" - an eventTextNode is triggered for each word that is not in an end tag or a start tag
	EventTypeEndDocument                     // triggered right after the > symbol of </html> is read
	EventTypeOpeningTag                      // triggered right after the h in <html> is read (and other tags too)
	EventTypeClosingTag                      // triggered right after the / in </html> is read (and other tags too)
	EventTypeAttributeKey                    // triggered right after the h in <a href="example.com"> is read
	EventTypeAttributeValue                  // triggered right after the = in <a href="example.com"> is read
	EventTypeComment                         //
	EventTypeUnknown                         // used to signal that an event.type is not set
)

var (
	eventErrorInvalidWhiteSpace = errors.New("Invalid whitespace")
	eventErrorInvalidTag        = errors.New("Invalid tag")
	eventErrorInvalidEndOfFile  = errors.New("Invalid end of file")
	eventErrorInvalidNewLine    = errors.New("Invalid new line")
)

type ParserEvent struct {
	Type   eventType
	Buffer []rune
	Error  error
}
type parserState byte

const (
	parserStateOpeningTag parserState = iota
	parserStateUnknown
	parserStateAttributeKey
	parserStateAttributeValue
)

func Parse(input io.Reader) func(*ParserEvent) {
	var (
		nextRune              rune
		reader                            = bufio.NewReader(input)
		currentState          parserState = parserStateUnknown
		numberOfSpacesSkipped int         = 0
		lineNumber            int         = 1
	)

	// skip whitespace and read the first character
	// at the end of input, returns 0
	var skipWhiteSpace = func() rune {
		numberOfSpacesSkipped = 0
	read:
		character, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return 0
			} else {
				panic("Unknown Error, [DEBUGGING]: file: index/htmlparser.go, function: skipWhiteSpace")
			}
		}
		if nextRune == '\n' {
			lineNumber++
		}
		if !unicode.IsSpace(character) {
			return character
		} else {
			numberOfSpacesSkipped++
		}
		goto read
	}

	return func(event *ParserEvent) {
		event.Buffer = event.Buffer[:0]
		event.Error = nil

	functionStart:
		nextRune = skipWhiteSpace() //skipWhiteSpace between tokens
		if nextRune == 0 {
			event.Type = EventTypeEndDocument
			println("END")
			return
		}

		// these states are used to make sure we don't break out of an opening tag in between calls untill we've read the entire tag
		if currentState == parserStateOpeningTag {
			reader.UnreadRune()
			goto attributeKey
		}

		if currentState == parserStateAttributeKey {
			if nextRune == '=' {
				goto attributeValue
			}
			goto attributeKey
		}

		if currentState == parserStateAttributeValue {
			reader.UnreadRune()
			goto attributeKey
		}

		// possible tag
		if nextRune == '<' {
			nextRune = skipWhiteSpace() // allow <a < a
			if nextRune == '/' {
				goto closingTag
			}

			if unicode.IsLetter(nextRune) {
				reader.UnreadRune()
				goto openingTag
			}

			if nextRune == '!' {
				// discard comment or document declaration
				event.Type = EventTypeComment
				var secondRune = skipWhiteSpace()
				for {
					nextRune = skipWhiteSpace()
					if nextRune == 0 {
						event.Error = eventErrorInvalidEndOfFile
						return
					}

					if nextRune == '-' {
						nextRune, _, _ = reader.ReadRune()
						if nextRune == '-' {
							nextRune, _, _ = reader.ReadRune()
							if nextRune == '>' && secondRune == '-' {
								break
							}
						}
					}
					if nextRune == '>' && secondRune != '-' {
						break
					}
				}
				return
			}

			event.Error = fmt.Errorf("%w, line: %v", eventErrorInvalidTag, lineNumber)
			return
		}
		reader.UnreadRune()
		goto textNode;

	openingTag:
		{
			currentState = parserStateOpeningTag
			event.Type = EventTypeOpeningTag
			// read up to first space or >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 {
					if nextRune == '>' { // <a >
						currentState = parserStateUnknown
						break
					}
					reader.UnreadRune() // we have a rune, but we skipped a few spaces to get it, it is an attributeKey
					break
				}

				if nextRune == '>' { // <a>
					currentState = parserStateUnknown
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	attributeKey:
		{
			currentState = parserStateAttributeKey
			event.Type = EventTypeAttributeKey
			// read up until a space or '=' or >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 && nextRune != '=' && nextRune != '>' { // a second key
					currentState = parserStateAttributeKey
					break
				}

				// NOTE: numberOfSpacesSkipped makes no difference here
				if nextRune == '=' { // <a href = ..
					reader.UnreadRune() // NOTE: '=' is used by the next call to go to attributeValue state (see functionStart label)
					currentState = parserStateAttributeKey
					break
				}
				if nextRune == '>' { // <a blackButton  >,
					currentState = parserStateUnknown
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	attributeValue:
		{
			currentState = parserStateAttributeValue
			event.Type = EventTypeAttributeValue
			// read up until '>' or space, or quotes (if it started with a quote)
			var firstRune = skipWhiteSpace()
			if !(firstRune == '\'' || firstRune == '"') {
				reader.UnreadRune()
			}
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if nextRune == '>' {
					currentState = parserStateUnknown
					break
				}

				if numberOfSpacesSkipped > 0 && !(firstRune == '\'' || firstRune == '"') {
					reader.UnreadRune();
					break
				}
				if nextRune == firstRune && (firstRune == '\'' || firstRune == '"') {
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}

	textNode:
		{
			event.Type = EventTypeTextNode
			// read up to a space or <
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if numberOfSpacesSkipped > 0 {
					reader.UnreadRune()
					if len(event.Buffer) < 1 {
						// NOTE: so we dont return empty space
						goto functionStart
					}
					break
				}
				if nextRune == '<' {
					reader.UnreadRune();
					break;
				}

				event.Buffer = append(event.Buffer, nextRune)
			}
			if len(event.Buffer) < 1 {
				// NOTE: so we dont return empty space
				goto functionStart
			}
			return
		}

	closingTag:
		{
			event.Type = EventTypeClosingTag
			// read up to >
			for {
				nextRune = skipWhiteSpace()
				if nextRune == 0 {
					event.Error = eventErrorInvalidEndOfFile
					break
				}
				if nextRune == '>' {
					break
				}
				event.Buffer = append(event.Buffer, nextRune)
			}
			return
		}
	}
}

/*
	this is an implmentation of an html event-based parser
*/

/*
func parseHTMLFile(file *os.File) func(*parserEvent) {
	var nextbyte byte
	var numberOfSpacesSkipped int
	var reader = newFileReader(file)
	var currentState parserState = parserStateUnknown

	return func(event *parserEvent) {
		event.Buffer = event.Buffer[:0]
		event.Error = nil

	functionStart:
		nextbyte, _ = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter) // skip whitespace between tokens
		if nextbyte == 0 {
			event.eventType = EventTypeEndDocument
			return
		}

		// these states are used to make sure we don't break out of an opening tag in between calls untill we've read the entire tag
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

		// tag
		if nextbyte == '<' {
			reader.readLowerCase(skipWhiteSpace, consumeFirstCharacter) // discard delimiter
			nextbyte, numberOfSpacesSkipped = reader.readLowerCase(skipWhiteSpace, dontConsumeFirstCharacter)
			if numberOfSpacesSkipped > 0 {
				event.eventError = ErrorInvalidWhiteSpace
				return
			}

			// the second byte after < has to be a letter,  /,  !
			// discard comments and document declaration
			if nextbyte == '!' {
				event.eventType = EventTypeComment
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
		// tag
		goto textNode

	openingTag:
		{
			currentState = parserStateOpeningTag
			event.eventType = EventTypeOpeningTag
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
				event.Buffer = append(event.Buffer, b)
			}
			return
		}

	attributeKey:
		{
			currentState = parserStateAttributeKey
			event.eventType = EventTypeAttributeKey
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
				event.Buffer = append(event.Buffer, b)
			}
			return
		}

	attributeValue:
		{
			currentState = parserStateAttributeValue
			event.eventType = EventTypeAttributeValue
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
				event.Buffer = append(event.Buffer, b)
			}
			return
		}

	closingTag:
		{
			event.eventType = EventTypeClosingTag
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
				event.Buffer = append(event.Buffer, b)
			}
			return
		}

	textNode:
		{
			event.eventType = EventTypeTextNode
			// read up to space || <
			for {
				nextbyte, _ := reader.readLowerCase(dontSkipWhiteSpace, dontConsumeFirstCharacter)
				if nextbyte == 0 {
					event.eventError = eventErrorInvalidEndOfFile
					break
				}
				if unicode.IsSpace(rune(nextbyte)) || nextbyte == '<' {
					// EDGE CASE: skip returning an empty textNode
					if !(len(event.Buffer) > 0) {
						goto functionStart
					}
					break
				}
				b, _ := reader.readLowerCase(dontSkipWhiteSpace, consumeFirstCharacter)
				// ignore punctuations
				if !ignorePunctuation[b] {
					event.Buffer = append(event.Buffer, b)
				}
			}
			return
		}

		return
	}
}

type fileReader struct {
	buf  *bufio.Reader
	line uint64
	// cache byte
	cache rune
}

func newFileReader(file *os.File) *fileReader {
	return &fileReader{
		line:  1,
		buf:   bufio.NewReader(file),
		cache: 0,
	}
}

// if skipWhiteSpace is true,  nextByte will skip any tabs, spaces, newlines, carriage returns to get to the firs non whitespace character
// if skipWhiteSpace is false, nextByte will return the first byte it encounters, whether a whitespace character or not
// if consumeFirstCharacter is true, nextByte will consume the first character it encounters before returning it, meaning subsequent calls to nextByte will return the bytes after the consumed one
// if consumeFirstCharacter is false, nextByte will peek and return the byte, meaning a subsequent call will return the same byte as the one before
// at end of file, nextByte will return 0
const (
	skipWhiteSpace            = true
	dontSkipWhiteSpace        = false
	consumeFirstCharacter     = true
	dontConsumeFirstCharacter = false
)

// func (reader *fileReader) readLowerCase(skipWhiteSpace bool, consumeFirstCharacter bool) (byte, int) {
// 	var (
// 		numberOfSpacesSkipped int = 0
// 		nextbyte              byte
// 		ok                    bool = false
// 		nextCharacter         rune
// 		size                  int
// 		err                   error
// 	)
//
// 	if reader.cache != 0 {
// 		nextbyte = reader.cache
// 		if skipWhiteSpace && unicode.IsSpace(rune(nextbyte)) {
// 			reader.cache = 0
// 			goto read
// 		}
// 		if consumeFirstCharacter {
// 			reader.cache = 0
// 		}
// 		goto returnByte
// 	}
//
// read:
// 	nextCharacter, size, err = reader.buf.ReadRune()
// 	if err == io.EOF {
// 		return 0, 0
// 	}
// 	if nextCharacter == '\n' || nextCharacter == '\v' {
// 		reader.line++
// 	}
// 	if unicode.IsSpace(nextCharacter) && skipWhiteSpace {
// 		numberOfSpacesSkipped++
// 		goto read
// 	}
// 	if size == 1 {
// 		nextbyte = byte(nextCharacter)
// 		// if !consumeFirstCharacter {
// 		// 	reader.cache = nextbyte
// 		// }
// 		goto returnByte
// 	}
//
// 	// NOTE: if you change , make sure to change the normalize function too
// 	if nextbyte, ok = baseRune[nextCharacter]; ok {
// 		// if !consumeFirstCharacter {
// 		// 	reader.cache = nextbyte
// 		// }
// 		goto returnByte
// 	}
// 	// NOTE: if you change , make sure to change the normalize function too
// 	goto read // ignore any character we cannot decompose
//
// returnByte:
// 	// if nextbyte >= 'A' && nextbyte <= 'Z' {
// 	// 	nextbyte = nextbyte + 32 // lowercase
// 	// }
// 	if !consumeFirstCharacter {
// 		reader.cache = nextbyte
// 	}
// 	return nextbyte, numberOfSpacesSkipped
// }

func (reader *fileReader) nextCharacter(skipWhiteSpace bool, consumeFirstCharacter bool) (rune, int) {
	var (
		numberOfSpacesSkipped int  = 0
		ok                    bool = false
		nextCharacter         rune
		size                  int
		err                   error
	)

	if reader.cache != 0 {
		// nextbyte = reader.cache
		nextCharacter = reader.Cache
		if skipWhiteSpace && unicode.IsSpace(nextCharacter) {
			reader.cache = 0
			goto read
		}
		if consumeFirstCharacter {
			reader.cache = 0
		}
		goto returnCharacter
	}

read:
	nextCharacter, size, err = reader.buf.ReadRune()
	if err == io.EOF {
		return 0, 0
	}
	if nextCharacter == '\n' || nextCharacter == '\v' {
		reader.line++
	}
	if unicode.IsSpace(nextCharacter) && skipWhiteSpace {
		numberOfSpacesSkipped++
		goto read
	}
	// if size == 1 {
	// 	nextbyte = byte(nextCharacter)
	// 	// if !consumeFirstCharacter {
	// 	// 	reader.cache = nextbyte
	// 	// }
	// 	goto returnCharacter
	// }

	// NOTE: if you change , make sure to change the normalize function too
	// if nextbyte, ok = baseRune[nextCharacter]; ok {
	// 	// if !consumeFirstCharacter {
	// 	// 	reader.cache = nextbyte
	// 	// }
	// 	goto returnCharacter
	// }
	// NOTE: if you change , make sure to change the normalize function too
	// goto read // ignore any character we cannot decompose

returnCharacter:
	// if nextbyte >= 'A' && nextbyte <= 'Z' {
	// 	nextbyte = nextbyte + 32 // lowercase
	// }
	if !consumeFirstCharacter {
		reader.cache = nextCharacter
	}
	return nextCharacter, numberOfSpacesSkipped
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

var baseRune = map[rune]byte{
	// A
	'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
	'Ā': 'A', 'Ă': 'A', 'Ą': 'A', 'Ǎ': 'A',
	'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
	'ā': 'a', 'ă': 'a', 'ą': 'a', 'ǎ': 'a',

	// C
	'Ç': 'C', 'Ć': 'C', 'Ĉ': 'C', 'Ċ': 'C', 'Č': 'C',
	'ç': 'c', 'ć': 'c', 'ĉ': 'c', 'ċ': 'c', 'č': 'c',

	// D
	'Ď': 'D',
	'ď': 'd',

	// E
	'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
	'Ē': 'E', 'Ĕ': 'E', 'Ė': 'E', 'Ę': 'E', 'Ě': 'E',
	'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
	'ē': 'e', 'ĕ': 'e', 'ė': 'e', 'ę': 'e', 'ě': 'e',

	// G
	'Ĝ': 'G', 'Ğ': 'G', 'Ġ': 'G', 'Ģ': 'G',
	'ĝ': 'g', 'ğ': 'g', 'ġ': 'g', 'ģ': 'g',

	// H
	'Ĥ': 'H',
	'ĥ': 'h',

	// I
	'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
	'Ĩ': 'I', 'Ī': 'I', 'Ĭ': 'I', 'Į': 'I', 'Ǐ': 'I',
	'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
	'ĩ': 'i', 'ī': 'i', 'ĭ': 'i', 'į': 'i', 'ǐ': 'i',

	// N
	'Ñ': 'N', 'Ń': 'N', 'Ņ': 'N', 'Ň': 'N',
	'ñ': 'n', 'ń': 'n', 'ņ': 'n', 'ň': 'n',

	// O
	'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
	'Ō': 'O', 'Ŏ': 'O', 'Ő': 'O', 'Ǒ': 'O',
	'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
	'ō': 'o', 'ŏ': 'o', 'ő': 'o', 'ǒ': 'o',

	// R
	'Ŕ': 'R', 'Ŗ': 'R', 'Ř': 'R',
	'ŕ': 'r', 'ŗ': 'r', 'ř': 'r',

	// S
	'Ś': 'S', 'Ŝ': 'S', 'Ş': 'S', 'Š': 'S',
	'ś': 's', 'ŝ': 's', 'ş': 's', 'š': 's',

	// T
	'Ţ': 'T', 'Ť': 'T',
	'ţ': 't', 'ť': 't',

	// U
	'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
	'Ũ': 'U', 'Ū': 'U', 'Ŭ': 'U', 'Ů': 'U', 'Ű': 'U', 'Ų': 'U', 'Ǔ': 'U',
	'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
	'ũ': 'u', 'ū': 'u', 'ŭ': 'u', 'ů': 'u', 'ű': 'u', 'ų': 'u', 'ǔ': 'u',

	// W
	'Ẁ': 'W', 'Ẃ': 'W', 'Ŵ': 'W', 'Ẅ': 'W',
	'ẁ': 'w', 'ẃ': 'w', 'ŵ': 'w', 'ẅ': 'w',

	// Y
	'Ỳ': 'Y', 'Ý': 'Y', 'Ŷ': 'Y', 'Ÿ': 'Y', 'Ỹ': 'Y',
	'ỳ': 'y', 'ý': 'y', 'ŷ': 'y', 'ÿ': 'y', 'ỹ': 'y',

	// Z
	'Ź': 'Z', 'Ż': 'Z', 'Ž': 'Z',
	'ź': 'z', 'ż': 'z', 'ž': 'z',
}*/
