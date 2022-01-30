package logger

import (
	"fmt"

	"github.com/liamg/tml"
)

type Logger struct {
	title    string
	silenced bool
}

func New() Logger {
	return Logger{}
}

func (logger Logger) Silenced() Logger {
	logger.silenced = true
	return logger
}

func (logger Logger) WithTitle(title string) Logger {
	logger.title = title
	return logger
}

func (logger Logger) Printf(format string, args ...interface{}) {
	if logger.silenced {
		return
	}
	_ = tml.Printf("\r<blue>[</blue><yellow>+</yellow><blue>]</blue>")
	if logger.title != "" {
		_ = tml.Printf("<blue>[</blue><red>%s</red><blue>]</blue>", logger.title)
	}
	line := fmt.Sprintf(format, args...)
	_ = tml.Printf(" %s\r\n", line)
}
