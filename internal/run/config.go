package run

import (
	"io"
)

// Config contains configuration applicable to all helm commands
type Config struct {
	Debug          bool
	KubeConfig     string
	Values         string
	StringValues   string
	ValuesFiles    []string
	Namespace      string
	Token          string
	SkipTLSVerify  bool
	Certificate    string
	APIServer      string
	ServiceAccount string
	Stdout         io.Writer
	Stderr         io.Writer
}
