package urlvector

import (
	"github.com/koykov/vector"
)

type URLHelper struct{}

var (
	urlHelper = &URLHelper{}
)

func (h *URLHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		return unescape(b)
	}
	return b
}
