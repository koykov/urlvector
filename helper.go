package urlvector

import (
	"github.com/koykov/vector"
)

const (
	flagEscape = uint8(1)
)

type URLHelper struct{}

var (
	urlHelper = &URLHelper{}
)

func (h *URLHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.CheckFlag(flagEscape) {
		return unescape(b)
	}
	return b
}
