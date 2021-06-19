package urlvector

import (
	"github.com/koykov/vector"
)

type URLHelper struct{}

var (
	urlHelper = &URLHelper{}
)

func (h *URLHelper) Indirect(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = unescape(b)
		p.SetLimit(len(b))
	}
	return b
}
