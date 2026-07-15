package parser

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

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
 */
func (reader *fileReader) read(skipWhiteSpace bool, consumeFirstCharacter bool) (byte, int) {
	var numberOfSpacesSkipped int = 0;
	if reader.cache != 0 {
		var b = reader.cache
		if skipWhiteSpace && unicode.IsSpace(rune(b)) {
			reader.cache = 0
			goto read;
		}
		if consumeFirstCharacter {
			reader.cache = 0
		}
		return b, 0;
	}

read:
	nextCharacter, size, err := reader.buf.ReadRune()
	if err == io.EOF {
		return 0, 0
	}
	if nextCharacter == '\n' || nextCharacter == '\v' {
		reader.line++
	}
	if unicode.IsSpace(nextCharacter) && skipWhiteSpace {
		numberOfSpacesSkipped++;
		goto read
	}
	if size == 1 {
		if !consumeFirstCharacter {
			reader.cache = byte(nextCharacter);
		}
		return byte(nextCharacter), numberOfSpacesSkipped
	}

	if b, ok := baseRune[nextCharacter]; ok {
		if !consumeFirstCharacter {
			reader.cache = b;
		}
		return b, numberOfSpacesSkipped;
	}
	goto read // ignore any character we cannot decompose
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
