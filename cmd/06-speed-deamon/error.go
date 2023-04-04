package main

import "io"

func WriteErrorMessage(w io.Writer, msg string) {
	// return an error about writing the error message does not make any sense so this function returns nothing
	if err := WriteMessageType(w, ErrorMessageType); err != nil {
		return
	}
	if err := WriteString(w, msg); err != nil {
		return
	}

	return
}
