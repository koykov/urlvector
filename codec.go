package urlvector

import (
	"io"

	"github.com/koykov/vector"
)

type Codec struct {
	vector.BaseCodec
}

func (h Codec) Decode(p *vector.Byteptr) ([]byte, error) {
	b := p.RawBytes()
	if p.CheckBit(flagEscape) {
		p.SetBit(flagEscape, false)
		b = unescape(b)
		p.SetLen(len(b))
	}
	return b, nil
}

func (h Codec) Beautify(_ io.Writer, _ *vector.Node) error {
	return nil
}

func (h Codec) Marshal(w io.Writer, node *vector.Node) error {
	_, err := w.Write(node.Bytes())
	return err
}
