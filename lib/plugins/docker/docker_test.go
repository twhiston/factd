package docker

import (
	"github.com/twhiston/factd/lib/plugins"
	"testing"
)

func TestImplements(t *testing.T) {
	var _ plugins.Plugin = (*Docker)(nil)
}

func BenchmarkFactCollection(b *testing.B) {
	p := Docker{}
	for n := 0; n < b.N; n++ {
		_, err := p.Facts()
		if err != nil {
			b.Error(err)
		}
	}
}
