package urlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

const (
	// Indexes of URL parts in the nodes array.
	idxScheme      = 1
	idxSlashes     = 2
	idxAuth        = 3
	idxUsername    = 4
	idxPassword    = 5
	idxHost        = 6
	idxHostname    = 7
	idxPort        = 8
	idxPath        = 9
	idxQueryOrigin = 10
	idxHash        = 11
	idxQuery       = 12

	// Vector level flags.
	flagCopy        = 8
	flagBufMod      = 9
	flagQueryParsed = 10
	// Byteptr level flags.
	flagEscape = 8
	flagBufSrc = 9
)

// Parser object.
type Vector struct {
	vector.Vector
}

// Make new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = urlHelper
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// Parse source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// Copy source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// Copy source string and parse it.
func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

// Get scheme node.
func (vec *Vector) Scheme() *vector.Node {
	return vec.getByIdx(idxScheme)
}

// Returns true if URL is started with slashes.
func (vec *Vector) Slashes() bool {
	return vec.getByIdx(idxSlashes).Bool()
}

// Get auth node (contains both username and password substrings).
func (vec *Vector) Auth() *vector.Node {
	return vec.getByIdx(idxAuth)
}

// Get username node.
func (vec *Vector) Username() *vector.Node {
	return vec.getByIdx(idxUsername)
}

// Get password node.
func (vec *Vector) Password() *vector.Node {
	return vec.getByIdx(idxPassword)
}

// Get host node (contains both hostname/IP and port substrings).
func (vec *Vector) Host() *vector.Node {
	return vec.getByIdx(idxHost)
}

// Get hostname node (similar to Host(), but excludes port).
func (vec *Vector) Hostname() *vector.Node {
	return vec.getByIdx(idxHostname)
}

// Get port as integer.
func (vec *Vector) Port() int {
	i, _ := vec.getByIdx(idxPort).Int()
	return int(i)
}

// Get path node.
func (vec *Vector) Path() *vector.Node {
	return vec.getByIdx(idxPath)
}

// Get query node.
func (vec *Vector) Query() *vector.Node {
	query := vec.getByIdx(idxQuery)
	if !vec.CheckBit(flagQueryParsed) {
		vec.SetBit(flagQueryParsed, true)
		vec.parseQueryParams(query)
	}
	return query
}

// Internal query getter.
func (vec *Vector) queryOrigin() *vector.Node {
	return vec.getByIdx(idxQueryOrigin)
}

// Get hash node.
func (vec *Vector) Hash() *vector.Node {
	return vec.getByIdx(idxHash)
}

func (vec *Vector) getByIdx(idx int) *vector.Node {
	node := vec.GetByIdx(idx)
	if node.Value().CheckBit(flagBufSrc) && vec.CheckBit(flagBufMod) {
		node.Value().TakeAddr(vec.Buf())
	}
	return node
}
