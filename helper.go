package urlvector

import (
	"io"

	"github.com/koykov/vector"
)

type Helper struct{}

var (
	helper = Helper{}
)

func (h Helper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = unescape(b)
		p.SetLen(len(b))
	}
	return b
}

func (h Helper) Beautify(_ io.Writer, _ *vector.Node) error {
	return nil
}
