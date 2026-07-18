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
	// var name = []byte("collins gygye")
	// var name = []byte("caresses")
	// var name = []byte("ponies")
	// var name = []byte("caress")
	// var name = []byte("cats")
	// var name = []byte("feed")
	// var name = []byte("agreed")
	// var name = []byte("bled")
	// var name = []byte("plastered")
	// var name = []byte("motoring")
	// var name = []byte("sing")
	// var name = []byte("conflated")
	// var name = []byte("hopping")
	var name = []byte("tanned")
	print(string(name), " ==> ")
	normalize(&name)
	println(string(name))
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
				currentByte = input[index+1]
				if !(currentByte == 'a' || currentByte == 'e' || currentByte == 'i' || currentByte == 'o' || currentByte == 'u') {
					return true
				}
			}
		}
		return false
	}
	// end of containsVowel function variable

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
			var lastTwoBytes = string((*input)[(inputLen - 2):inputLen])
			if !containsVowel([]byte(lastTwoBytes)) &&
				lastTwoBytes[0] == lastTwoBytes[1] &&
				!(lastTwoBytes == "ss" || lastTwoBytes == "ll" || lastTwoBytes == "zz") {
				*input = (*input)[:inputLen-1]
			} else {
				// (m = 1 and *o) -> E
				//*o -	the stem ends cvc, where the second c is not W, X or Y (e.g. -WIL, -HOP)
				// read from the end 
				var count = inputLen - 1;


				// var lastFourBytes = string((*input)[(inputLen - 4):inputLen])
				// if (measure((*input)[0:(inputLen-3)]) == 1) && () {
				// }
			}
		}
	endOfStep1b:
	}
}
