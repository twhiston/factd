package factd

import (
	"github.com/twhiston/factd/lib/formatter"
	"github.com/twhiston/factd/lib/plugin"
	"io"
	"os"
)

// Config struct serves to pass Factd configuration
type Config struct {
	Formatter formatter.Formatter
	Plugins   map[string]plugin.Plugin
	Output    *io.Writer
}

// NewConfig returns a new config object with some default values
func NewConfig() *Config {
	c := new(Config)
	c.Plugins = make(map[string]plugin.Plugin)
	writer := io.Writer(os.Stdout)
	c.Output = &writer
	return c
}
