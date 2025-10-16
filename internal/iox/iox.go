package iox

import (
	"io"
	"os"
)

// Stream copies data from a reader to a writer.
func Stream(r io.Reader, w io.Writer) {
	// If the writer is nil, default to os.Stdout.
	if w == nil {
		w = os.Stdout
	}
	_, _ = io.Copy(w, r)
}
