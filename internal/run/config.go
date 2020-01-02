package run

import (
	"io"
)

// Config contains configuration applicable to all helm commands
type Config struct {
	Debug     bool
	Namespace string
	Stdout    io.Writer
	Stderr    io.Writer
}
