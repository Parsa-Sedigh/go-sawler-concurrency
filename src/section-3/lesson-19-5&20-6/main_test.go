package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	// In the main function, we're writing to the console, so we need to keep a copy of whatever os standard out is. Let's call it `stdOut`.
	stdOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	// read the results into a variable called result:
	result, _ := io.ReadAll(r)

	// convert it into a string:
	output := string(result)

	// set os.Stdout back to it's original value
	os.Stdout = stdOut

	if !strings.Contains(output, "$34320.00") {
		t.Error("wrong balance returned")
	}
}
