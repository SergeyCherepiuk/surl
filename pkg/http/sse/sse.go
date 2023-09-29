package sse

import (
	"fmt"
	"io"
)

func Send(w io.Writer, event string, message []byte) (int, error) {
	return fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, message)
}
