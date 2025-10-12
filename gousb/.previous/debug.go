package gousb

// To enable internal debugging, set the GOUSB_DEBUG environment variable.

import (
	"io"
	"log" // TODO(kevlar): make a logger
	"os"
)

var debug *log.Logger

const debugEnvVarName = "GOUSB_DEBUG"

func init() {
	out := io.Writer(io.Discard)
	if os.Getenv(debugEnvVarName) != "" {
		out = os.Stderr
	}
	debug = log.New(out, "gousb: ", log.LstdFlags|log.Lshortfile)
}
