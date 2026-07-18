package index

// lowercase all the letters
// remove punctuations
func normalize(input []byte) {
	var measure = func(input []byte) int {
		return 1
	}

	var m = measure([]byte("colllins"))
	println(m)
}
