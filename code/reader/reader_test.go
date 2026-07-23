package reader

import "testing"

func TestReadByte(t *testing.T) {
	var testHTML = "<diŹ"
	var testReader = NewReader(testHTML)

	var b byte
	var read = func(expected byte) {
		testReader.ReadByte(&b)
		if b != expected {
			t.Fatalf("[ERROR]: Expected byte: %v, got byte: %v", string(expected), b)
		}
	}

	read('<')
	read('d')
	read('i')
	read('Z')
	read(0)
}

func TestReadNormalizedRange(t *testing.T) {
	var testHTML = "<!hello>Ò"
	var expectedOutput = "helloo"
	var testReader = NewReader(testHTML)
	var targetBuffer = make([]byte, 0)

	testReader.ReadNormalizedRange(&targetBuffer, 0)
	if string(targetBuffer) != expectedOutput {
		t.Fatalf("[ERROR]: Expected string: %v, got string: %v", expectedOutput, string(targetBuffer))
	}
}
