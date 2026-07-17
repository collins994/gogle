/*
type eventType int
	errors I'm worried About:
		a < in a character string
*/
package index

import (
	"bufio"
	"errors"
	"io"
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
	eventType        eventType
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
				if !ignorePunctuation[b] {
					event.eventBuffer = append(event.eventBuffer, b)
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

/** fileReader **/
type fileReader struct {
	buf   *bufio.Reader
	line  uint64
	cache byte
}

func newFileReader(file *os.File) *fileReader {
	return &fileReader{
		line:  1,
		buf:   bufio.NewReader(file),
		cache: 0,
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

func (reader *fileReader) readLowerCase(skipWhiteSpace bool, consumeFirstCharacter bool) (byte, int) {
	var (
		numberOfSpacesSkipped int = 0
		nextbyte              byte
		ok                    bool = false
		nextCharacter         rune
		size                  int
		err                   error
	)

	if reader.cache != 0 {
		nextbyte = reader.cache
		if skipWhiteSpace && unicode.IsSpace(rune(nextbyte)) {
			reader.cache = 0
			goto read
		}
		if consumeFirstCharacter {
			reader.cache = 0
		}
		goto returnByte
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
	if size == 1 {
		nextbyte = byte(nextCharacter)
		if !consumeFirstCharacter {
			reader.cache = nextbyte
		}
		goto returnByte
	}

	if nextbyte, ok = baseRune[nextCharacter]; ok {
		if !consumeFirstCharacter {
			reader.cache = nextbyte
		}
		goto returnByte
	}
	goto read // ignore any character we cannot decompose

returnByte:
	if nextbyte >= 'A' && nextbyte <= 'Z' {
		nextbyte = nextbyte + 32 // lowercase
	}
	return nextbyte, numberOfSpacesSkipped
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
}
