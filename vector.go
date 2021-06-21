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
	flagQuerySorted = 11
	flagQueryMod    = 12
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

// Bytes reassembles the vector into a valid URL bytes array.
func (vec *Vector) Bytes() []byte {
	return vec.bytes(false)
}

// String reassembles the vector into a valid URL string.
func (vec *Vector) String() string {
	return fastconv.B2S(vec.Bytes())
}

// Escaped version of Bytes().
//
// In addition, escapes host and hash part.
func (vec *Vector) BytesEscaped() []byte {
	return vec.bytes(true)
}

// Escaped version of String().
//
// In addition, escapes host and hash part.
func (vec *Vector) StringEscaped() string {
	return fastconv.B2S(vec.BytesEscaped())
}

// Internal marshaller.
func (vec *Vector) bytes(esc bool) []byte {
	// Bytes uses internal buffer as destination array to assemble the URL. So we need to save current length of the
	// buffer and use it further as offset.
	offset := vec.BufLen()

	if scheme := vec.SchemeBytes(); len(scheme) > 0 {
		vec.BufAppend(scheme)
		vec.BufAppend(bSchemaSep)
	} else if vec.Slashes() {
		vec.BufAppend(bSlashes)
	}

	if username := vec.UsernameBytes(); len(username) > 0 {
		vec.BufAppend(username)
		if password := vec.PasswordBytes(); len(password) > 0 {
			vec.BufAppend(bColon)
			vec.BufAppend(password)
		}
		vec.BufAppend(bAt)
	}

	if hostname := vec.HostnameBytes(); len(hostname) > 0 {
		vec.BufAppend(hostname)
		if port := vec.getByIdx(idxPort); port.Value().Len() > 0 {
			vec.BufAppend(bColon)
			vec.BufAppend(port.Bytes())
		}
	}

	if path := vec.PathBytes(); len(path) > 0 {
		if path[0] != '/' {
			vec.BufAppend(bSlash)
		}
		if esc {
			vecEscape(vec, path, modePath)
		} else {
			vec.BufAppend(path)
		}
	}

	if query := vec.QueryBytes(); len(query) > 0 {
		if query[0] != '?' {
			vec.BufAppend(bQM)
		}
		vec.BufAppend(query)
	}

	if hash := vec.HashBytes(); len(hash) > 0 {
		if hash[0] != '#' {
			vec.BufAppend(bHash)
		}
		if esc {
			vecEscape(vec, hash, modeHash)
		} else {
			vec.BufAppend(hash)
		}
	}

	return vec.Buf()[offset:]
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

// Get query node with origin query params order.
func (vec *Vector) Query() *vector.Node {
	query := vec.getByIdx(idxQuery)
	if !vec.CheckBit(flagQueryParsed) {
		vec.SetBit(flagQueryParsed, true)
		vec.parseQueryParams(query)
	}
	return query
}

// Sort query params an AB order.
func (vec *Vector) QuerySort() *Vector {
	query := vec.Query()
	if !vec.CheckBit(flagQuerySorted) {
		vec.SetBit(flagQuerySorted, true)
		children := query.Children()
		quickSort(children, 0, len(children)-1)
		vec.SetBit(flagQueryMod, true)
	}
	return vec
}

// Internal query getter.
func (vec *Vector) queryOrigin() *vector.Node {
	queryOrigin := vec.getByIdx(idxQueryOrigin)

	if vec.CheckBit(flagQueryMod) {
		vec.SetBit(flagQueryMod, false)
		offset := vec.BufLen()
		limit := 1
		vec.BufAppendStr("?")
		query := vec.getByIdx(idxQuery)
		query.Each(func(idx int, node *vector.Node) {
			if idx > 0 {
				vec.BufAppendStr("&")
				limit++
			}
			key := vecEscape(vec, node.KeyBytes(), modeQuery)
			vec.BufAppendStr("=")
			val := vecEscape(vec, node.Bytes(), modeQuery)
			limit += len(key) + len(val) + 1
		})
		vec.SetBit(flagBufMod, true)

		queryOrigin.Value().Init(vec.Buf(), offset, limit)
		queryOrigin.Value().SetBit(flagBufSrc, true)
		vec.PutNode(queryOrigin.Index(), queryOrigin)
	}
	return queryOrigin
}

// Get hash node.
func (vec *Vector) Hash() *vector.Node {
	return vec.getByIdx(idxHash)
}

// Get node by index considering flags.
func (vec *Vector) getByIdx(idx int) *vector.Node {
	node := vec.GetByIdx(idx)
	if node.Value().CheckBit(flagBufSrc) && vec.CheckBit(flagBufMod) {
		node.Value().TakeAddr(vec.Buf())
	}
	return node
}
