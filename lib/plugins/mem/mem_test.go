package mem

import (
	"github.com/twhiston/factd/lib/plugins"
	"testing"
)

func TestImplements(t *testing.T) {
	var _ plugins.Plugin = (*Mem)(nil)
}

func BenchmarkFactCollection(b *testing.B) {
	p := Mem{}
	for n := 0; n < b.N; n++ {
		_, err := p.Facts()
		if err != nil {
			b.Error(err)
		}
	}
}
