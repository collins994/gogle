package main

import (
	"fmt"
	"github.com/collins994/gogle/code/index"
	"os"
)

func main2() {
	var IndexFile = index.Index()

	// for count := 1; count < 352; count++ {
	for count := 1; count <= 1; count++ {
		file, err := os.Open(fmt.Sprintf("gl2/%d.xhtml", count))
		if err != nil {
			fmt.Printf("[ERROR]: %v\n", err)
			continue
		}
		IndexFile(file)
	}
}

func main() {
	var names = []string{"caresses", "ponies", "caress", "cats", "feed",
		"agreed", "bled", "plastered", "motoring", "sing", "conflated", "hopping", "tanned",
		"hissing", "falling", "failing", "filing", "fizzed", "sized", "happy",
		"relational", "conditional", "rational", "valenci", "digitizer",
		"conformabli", "radicalli", "differentli", "vileli", "analogousli", "vietnamization", "predication",
		"operator", "feudalism", "decisiveness", "hopefulness", "callousness", "formaliti", "sensitiviti",
		"sensibiliti", "triplicate", "formative", "formalize", "electriciti", "electrical", "hopeful", "goodness",
		"revival", "allowance", "inference", "airliner", "gyroscopic", "adjustable", "defensible", "irritant", "replacement",
		"adjustment", "dependent", "adoption", "homologou", "communism", "activate", "angulariti", "homologous", "effective", "bowdlerize",
		"probate", "rate", "cease", "controll", "roll",
	}

	for _, n := range names {
		var name = []byte(n)
		normalize(&name)
		println(n, " ==> ", string(name))
	}
}

func normalize(input *[]byte) {
	var measure = func(input []byte) int {
		var m int
		var currentByte byte
		var previousVowel bool = false
		for index := 0; index < len(input); index++ {
			currentByte = input[index]
			if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
				// read the next byte, if it is a consonant, increament the measure
				previousVowel = true
				continue
			}
			if currentByte == 'y' {
				currentByte = input[index+1]
				if !(currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u') {
					previousVowel = true
					continue
				}
			}
			if previousVowel == true {
				m++
			}
			previousVowel = false
		}
		return m
	}
	// end of measure function variable

	var containsVowel = func(input []byte) bool {
		var currentByte byte
		for index := 0; index < len(input); index++ {
			currentByte = input[index]
			if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
				return true
			}
			if currentByte == 'y' {
				if index == 0 {
					continue
				}
				currentByte = input[index-1]
				if !(currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u') {
					return true
				}
			}
		}
		return false
	}
	// end of containsVowel function variable

	var endsInCVC = func(input []byte) bool {
		// read from the end
		var inputLen = len(input)
		var index uint = uint(inputLen - 1)
		var currentByte byte

		/* check for the cvc condition */
		// last byte
		currentByte = (input)[index]
		if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
			//goto endOfStep1b
			return false
		} else {
			if currentByte == 'y' { // y is considered a vowel if it is preceeded by a consonant
				currentByte = (input)[index-1]
				if !(currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u') {
					/// goto endOfStep1b // the byte preceeding y is a consonant, therefore y is a vowel, breaking the cvc condition
					return false
				}
			}
		}

		// second last byte
		index--
		currentByte = (input)[index]
		if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
			// goto thirdLastByte
			return false
		} else if currentByte == 'y' {
			currentByte = (input)[index-1]
			if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
				// goto endOfStep1b // the character preceeding y is a vowel, therefore y is a consonant, breaking the cvc condition
				return false
			} else {
				// goto thirdLastByte
				return false
			}
		}

		// third last byte
		index--
		currentByte = (input)[index]
		if currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u' {
			// goto endOfStep1b
			return false
		} else if currentByte == 'y' {
			currentByte = (input)[index-1]
			if !(currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u') {
				// goto endOfStep1b // y is preceeded by a consonant, that makes it a vowel, breaking the cvc condition
				return false
			}
		}
		return true
	}

	var inputLen int
	// step 1a
	// NOTE: we start with the longest suffix, order matters
	{
		inputLen = len(*input)
		if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "sses" { // sses -> ss
			*input = (*input)[:inputLen-2]
		} else if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "ies" { // ies -> i
			*input = (*input)[:inputLen-2]
		} else if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "ss" { // ss -> ss
		} else if inputLen > 1 && string((*input)[(inputLen-1):inputLen]) == "s" { // s ->
			*input = (*input)[:inputLen-1]
		}
	}

	// step 1b
	{
		inputLen = len(*input)
		if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "eed" { // (m > 0) eed -> ee
			if measure((*input)[0:(inputLen-3)]) > 0 {
				*input = (*input)[:inputLen-1]
				goto endOfStep1b
			}
		} else if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "ed" { // (*v*) ed ->
			if containsVowel((*input)[0 : inputLen-2]) {
				*input = (*input)[:inputLen-2]
			}
		} else if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "ing" { // (*v*) ing ->
			if containsVowel((*input)[0 : inputLen-3]) {
				*input = (*input)[:inputLen-3]
			}
		}

		inputLen = len(*input)
		if inputLen > 2 && (string((*input)[(inputLen-2):inputLen]) == "at") { // at -> ate
			*input = append(*input, 'e')
		} else if inputLen > 2 && (string((*input)[(inputLen-2):inputLen]) == "bl") { // bl -> ble
			*input = append(*input, 'e')
		} else if inputLen > 2 && (string((*input)[(inputLen-2):inputLen]) == "iz") { // iz -> ize
			*input = append(*input, 'e')
		} else {
			// (*d and not (*L or *S or *Z)) -> single letter
			// the last two bytes are double consonants (same) but not "ss" || "ll" || "zz"
			if len(*input) < 2 {
				goto endOfStep1b
			}
			var lastTwoBytes = string((*input)[(inputLen - 2):inputLen])
			if !containsVowel([]byte(lastTwoBytes)) &&
				lastTwoBytes[0] == lastTwoBytes[1] &&
				!(lastTwoBytes == "ss" || lastTwoBytes == "ll" || lastTwoBytes == "zz") {
				*input = (*input)[:inputLen-1]
			} else {
				// (m = 1 and *o) -> E
				//*o -	the stem ends cvc, where the second c is not W, X or Y (e.g. -WIL, -HOP)
				// var currentByte byte
				if endsInCVC(*input) && measure(*input) == 1 {
					*input = append(*input, 'e')
				}
			}
		}
	endOfStep1b:
	}

	// step 1c
	{
		inputLen = len(*input)
		if inputLen > 1 && (*input)[inputLen-1] == 'y' && containsVowel((*input)[0:inputLen-2]) { // (*v*) Y	->	I
			(*input)[inputLen-1] = 'i'
		}
	}

	// step 2
	{
		inputLen = len(*input)
		var penultimateByte byte
		if inputLen < 2 {
			goto endOfStep2
		}
		penultimateByte = (*input)[inputLen-2]
		if penultimateByte == 'a' {
			if inputLen > 7 && string((*input)[(inputLen-7):inputLen]) == "ational" && measure((*input)[0:inputLen-7]) > 0 { // (m > 0) ational -> ate
				*input = (*input)[:inputLen-5]
				*input = append(*input, 'e')
			} else if inputLen > 6 && string((*input)[(inputLen-6):inputLen]) == "tional" && measure((*input)[0:inputLen-6]) > 0 { // (m > 0) tional -> tion
				*input = (*input)[:inputLen-2]
			}
		}

		if penultimateByte == 'c' {
			if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "enci" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) ENCI->ENCE
				(*input)[inputLen-1] = 'e'
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "anci" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) ANCI->ANCE
				(*input)[inputLen-1] = 'e'
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "izer" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) izer -> ize
				*input = (*input)[:inputLen-1]
			}
		}

		if penultimateByte == 'l' {
			if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "abli" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) abli -> able
				(*input)[inputLen-1] = 'e'
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "alli" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) alli -> al
				*input = (*input)[:inputLen-2]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "entli" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) entli -> ent
				*input = (*input)[:inputLen-2]
			} else if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "eli" && measure((*input)[0:inputLen-3]) > 0 { // (m>0) eli -> e
				*input = (*input)[:inputLen-2]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "ousli" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) ousli -> ous
				*input = (*input)[:inputLen-2]
			}
		}

		if penultimateByte == 'o' {
			if inputLen > 7 && string((*input)[(inputLen-7):inputLen]) == "ization" && measure((*input)[0:inputLen-7]) > 0 { // (m>0) ization -> ize
				*input = (*input)[:inputLen-5]
				*input = append(*input, 'e')
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "ation" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) ation -> ate
				*input = (*input)[:inputLen-3]
				*input = append(*input, 'e')
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ator" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) ator -> ate
				*input = (*input)[:inputLen-2]
				*input = append(*input, 'e')
			}
		}

		if penultimateByte == 's' {
			if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "alism" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) alism -> al
				*input = (*input)[:inputLen-3]
			} else if inputLen > 7 && string((*input)[(inputLen-7):inputLen]) == "iveness" && measure((*input)[0:inputLen-7]) > 0 { // (m>0) iveness -> ive
				*input = (*input)[:inputLen-4]
			} else if inputLen > 7 && string((*input)[(inputLen-7):inputLen]) == "fulness" && measure((*input)[0:inputLen-7]) > 0 { // (m>0) fulness -> ful
				*input = (*input)[:inputLen-4]
			} else if inputLen > 7 && string((*input)[(inputLen-7):inputLen]) == "ousness" && measure((*input)[0:inputLen-7]) > 0 { // (m>0) ousness -> ous
				*input = (*input)[:inputLen-4]
			}
		}

		if penultimateByte == 't' {
			if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "aliti" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) aliti -> al
				*input = (*input)[:inputLen-3]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "iviti" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) iviti -> ive
				*input = (*input)[:inputLen-3]
				*input = append(*input, 'e')
			} else if inputLen > 6 && string((*input)[(inputLen-6):inputLen]) == "biliti" && measure((*input)[0:inputLen-6]) > 0 { // (m>0) biliti -> ble
				*input = (*input)[:inputLen-5]
				*input = append(*input, 'l')
				*input = append(*input, 'e')
			}
		}
	endOfStep2:
	}

	// step 3
	{
		inputLen = len(*input)
		var ultimateByte = (*input)[inputLen-1]
		if ultimateByte == 'e' {
			if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "icate" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) icate -> ic
				*input = (*input)[:inputLen-3]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "ative" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) ative ->
				*input = (*input)[:inputLen-5]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "alize" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) alize -> al
				*input = (*input)[:inputLen-3]
			}
		} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "iciti" && measure((*input)[0:inputLen-5]) > 0 { // (m>0) iciti -> ic
			*input = (*input)[:inputLen-3]
		} else if ultimateByte == 'l' {
			if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ical" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) ical -> ic
				*input = (*input)[:inputLen-2]
			} else if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "ful" && measure((*input)[0:inputLen-3]) > 0 { // (m>0) ful ->
				*input = (*input)[:inputLen-3]
			}
		} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ness" && measure((*input)[0:inputLen-4]) > 0 { // (m>0) ness ->
			*input = (*input)[:inputLen-4]
		}
	}

	// step 4
	{
		inputLen = len(*input)
		var penultimateByte byte
		if inputLen < 2 {
			goto endOfStep4
		}
		penultimateByte = (*input)[inputLen-2]
		if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "al" && measure((*input)[0:inputLen-2]) > 1 { // (m>1) al ->
			*input = (*input)[:inputLen-2]
		} else if penultimateByte == 'c' {
			if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ance" && measure((*input)[0:inputLen-4]) > 1 { // (m>1) ance ->
				*input = (*input)[:inputLen-4]
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ence" && measure((*input)[0:inputLen-4]) > 1 { // (m>1) ence ->
				*input = (*input)[:inputLen-4]
			}
		} else if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "er" && measure((*input)[0:inputLen-2]) > 1 { // (m>1) er ->
			*input = (*input)[:inputLen-2]
		} else if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "ic" && measure((*input)[0:inputLen-2]) > 1 { // (m>1) ic ->
			*input = (*input)[:inputLen-2]
		} else if penultimateByte == 'l' {
			if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "able" && measure((*input)[0:inputLen-4]) > 1 { // (m>1) able ->
				*input = (*input)[:inputLen-4]
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ible" && measure((*input)[0:inputLen-4]) > 1 { // (m>1) ible ->
				*input = (*input)[:inputLen-4]
			}
		} else if penultimateByte == 'n' {
			if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "ant" && measure((*input)[0:inputLen-3]) > 1 { // (m>1) ant ->
				*input = (*input)[:inputLen-3]
			} else if inputLen > 5 && string((*input)[(inputLen-5):inputLen]) == "ement" && measure((*input)[0:inputLen-5]) > 1 { // (m>1) ement ->
				*input = (*input)[:inputLen-5]
			} else if inputLen > 4 && string((*input)[(inputLen-4):inputLen]) == "ment" && measure((*input)[0:inputLen-4]) > 1 { // (m>1) ment ->
				*input = (*input)[:inputLen-4]
			} else if inputLen > 3 && string((*input)[(inputLen-3):inputLen]) == "ent" && measure((*input)[0:inputLen-3]) > 1 { // (m>1) ent ->
				*input = (*input)[:inputLen-3]
			}
		} else if inputLen > 3 && ((*input)[inputLen-(1+3)] == 's' || (*input)[inputLen-(1+3)] == 't') && string((*input)[(inputLen-3):inputLen]) == "ion" && measure((*input)[0:inputLen-3]) > 1 { // (m>1 and (*S or *T)) ION ->
			*input = (*input)[:inputLen-3]
		} else if inputLen > 2 && string((*input)[(inputLen-2):inputLen]) == "ou" && measure((*input)[0:inputLen-2]) > 1 { // (m>1) ou ->
			*input = (*input)[:inputLen-2]
		} else if inputLen > 3 && measure((*input)[0:inputLen-3]) > 1 &&
			(string((*input)[(inputLen-3):inputLen]) == "ism" || // (m > 1) ism ->
				string((*input)[(inputLen-3):inputLen]) == "ate" || // (m>1) ate ->
				string((*input)[(inputLen-3):inputLen]) == "iti" || // (m>1) iti ->
				string((*input)[(inputLen-3):inputLen]) == "ous" || // (m>1) ous ->
				string((*input)[(inputLen-3):inputLen]) == "ive" || // (m>1) ive ->
				string((*input)[(inputLen-3):inputLen]) == "ize") { // (m>1) ize ->
			*input = (*input)[:inputLen-3]
		}
	endOfStep4:
	}

	// step 5a
	{
		inputLen = len(*input)
		if inputLen > 1 && (*input)[inputLen-1] == 'e' && measure((*input)[0:inputLen-1]) > 1 { // (m>1) e ->
			*input = (*input)[:inputLen-1]
		} else {
			// (m = 1 and not *o) e ->
			//*o -	the stem ends cvc, where the second c is not W, X or Y (e.g. -WIL, -HOP)
			if measure((*input)[0:inputLen-1]) == 1 && endsInCVC(*input) {
				*input = (*input)[:inputLen-1]
			}
		}
	}

	// steb 5b
	{
		inputLen = len(*input)
		var lastTwoBytes string
		// (m > 1 and *d and *L) -> single letter
		if measure(*input) <= 1 || inputLen < 2 {
			goto endOfStep5b
		}
		lastTwoBytes = string((*input)[(inputLen - 2):inputLen])
		// *d and *L - the input ends in double consonant "ll"
		if lastTwoBytes == "ll" {
			*input = (*input)[:inputLen-1]
		}
		endOfStep5b:
	}
}
