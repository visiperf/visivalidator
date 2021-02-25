package visivalidator

import (
	"bytes"
	"fmt"
)

// ValidationErrors is custom error containing potentially multiple errors
type ValidationErrors []error

func (ve ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")

	for i, e := range ve {
		buff.WriteString(fmt.Sprintf("%d. %s", i+1, e.Error()))
		buff.WriteString("\n")
	}

	return buff.String()
}
