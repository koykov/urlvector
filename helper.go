package urlvector

import (
	"github.com/koykov/vector"
)

const (
	flagEscape = 0
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
