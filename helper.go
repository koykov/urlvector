package urlvector

import (
	"github.com/koykov/vector"
)

var (
	urlHelper = &URLHelper{}
)

type URLHelper struct{}

func (h *URLHelper) ConvertByteptr(p *vector.Byteptr) []byte {
	b := p.RawBytes()
	if p.GetFlag(vector.FlagEscape) {
		return unescape(b)
	}
	return b
}
