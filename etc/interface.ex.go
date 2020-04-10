package main

import (
	"bytes"
	"io"
	"reflect"
)

const debug = false

func main() {
	var buf *bytes.Buffer
	if debug {
		buf = new(bytes.Buffer) // enable collection of output
	}
	f(buf) // NOTE: subtly incorrect!
	if debug {
		// ...use buf...
	}
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
	// ...do something...
	if out != nil {
		if !reflect.ValueOf(out).IsNil() {
			out.Write([]byte("done!\n"))
		}
	}
}
