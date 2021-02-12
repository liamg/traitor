package tml

import (
	"fmt"
)

// Printf works like fmt.Printf, but adds the option of using tags to apply colour or text formatting to the written text. For example "<red>some red text</red>".
// A full list of tags is available here: https://github.com/liamg/tml
func Printf(input string, a ...interface{}) error {
	format, err := Parse(input)
	if err != nil {
		return err
	}
	_, err = fmt.Printf(format, a...)
	return err
}

