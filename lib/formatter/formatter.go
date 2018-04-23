package formatter

import (
	"bytes"

	"github.com/twhiston/factd/lib/common"
)

// Formatter interface
type Formatter interface {
	Format(map[string]common.FactList) (*bytes.Buffer, error)
}
