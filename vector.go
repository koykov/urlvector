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

// Vector represents URL parser.
type Vector struct {
	vector.Vector
}

// NewVector makes new parser.
func NewVector() *Vector {
	vec := &Vector{}
	vec.Helper = helper
	return vec
}

// Parse source bytes.
func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

// ParseStr parses source string.
func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

// ParseCopyStr copies source string and parse it.
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

// BytesEscaped returns escaped URL bytes.
//
// In addition, escapes host and hash part.
func (vec *Vector) BytesEscaped() []byte {
	return vec.bytes(true)
}

// StringEscaped returns escaped URL string.
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

	if !vec.CheckBit(flagQueryParsed) {
		if query := vec.QueryBytes(); len(query) > 0 {
			if query[0] != '?' {
				vec.BufAppend(bQM)
			}
			vec.BufAppend(query)
		}
	} else {
		if query := vec.getByIdx(idxQuery); query.Limit() > 0 {
			vec.BufAppend(bQM)
			query.Each(func(idx int, node *vector.Node) {
				if idx > 0 {
					vec.BufAppend(bAmp)
				}
				vec.BufAppend(node.KeyBytes())
				vec.BufAppend(bEq)
				vecEscape(vec, node.Value().Bytes(), modeQuery)
			})
		}
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

// Scheme returns scheme node.
func (vec *Vector) Scheme() *vector.Node {
	return vec.getByIdx(idxScheme)
}

// Slashes indicates if URL is started with slashes.
func (vec *Vector) Slashes() bool {
	return vec.getByIdx(idxSlashes).Bool()
}

// Auth returns auth node (contains both username and password substrings).
func (vec *Vector) Auth() *vector.Node {
	return vec.getByIdx(idxAuth)
}

// Username returns username node.
func (vec *Vector) Username() *vector.Node {
	return vec.getByIdx(idxUsername)
}

// Password returns password node.
func (vec *Vector) Password() *vector.Node {
	return vec.getByIdx(idxPassword)
}

// Host returns host node (contains both hostname/IP and port substrings).
func (vec *Vector) Host() *vector.Node {
	return vec.getByIdx(idxHost)
}

// Hostname returns hostname node (similar to Host(), but excludes port).
func (vec *Vector) Hostname() *vector.Node {
	return vec.getByIdx(idxHostname)
}

// Port returns port as integer.
func (vec *Vector) Port() int {
	i, _ := vec.getByIdx(idxPort).Int()
	return int(i)
}

// Path returns path node.
func (vec *Vector) Path() *vector.Node {
	return vec.getByIdx(idxPath)
}

// Query returns query node with origin query params order.
func (vec *Vector) Query() *vector.Node {
	query := vec.getByIdx(idxQuery)
	if !vec.CheckBit(flagQueryParsed) {
		vec.SetBit(flagQueryParsed, true)
		vec.parseQueryParams(query)
	}
	return query
}

// QuerySort sorts query params an AB order.
func (vec *Vector) QuerySort() *Vector {
	query := vec.Query()
	if !vec.CheckBit(flagQuerySorted) {
		vec.SetBit(flagQuerySorted, true)
		children := query.ChildrenIndices()
		quickSort1(vec, children, 0, len(children)-1)
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

// Hash returns hash node.
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
