package urlvector

import (
	"bytes"
	"errors"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

const (
	// Offset of keys substrings in bKeys.
	offsetScheme      = 0
	offsetSlashes     = 6
	offsetAuth        = 13
	offsetUsername    = 17
	offsetPassword    = 25
	offsetHost        = 33
	offsetHostname    = 37
	offsetPort        = 45
	offsetPath        = 49
	offsetQueryOrigin = 57
	offsetHash        = 68
	offsetTrue        = 72
	offsetQuery       = 76

	// Length of keys substrings in bKeys.
	lenScheme      = 6
	lenSlashes     = 7
	lenAuth        = 4
	lenUsername    = 8
	lenPassword    = 8
	lenHost        = 4
	lenHostname    = 8
	lenPort        = 4
	lenPath        = 4
	lenQueryOrigin = 11
	lenHash        = 4
	lenTrue        = 4
	lenQuery       = 5

	// Max scheme length by https://en.wikipedia.org/wiki/List_of_URI_schemes#Official_IANA-registered_schemes.
	maxSchemaLen = 29
)

var (
	// Byte constants.
	bSpace     = []byte(" ")
	bBlob      = []byte("blob:")
	bSchemaSep = []byte("://")
	bSlashes   = []byte("//")
	bSlash     = []byte("/")
	bColon     = []byte(":")
	bAt        = []byte("@")
	bQM        = []byte("?")
	bAmp       = []byte("&")
	bHash      = []byte("#")
	bQB        = []byte("[]")
	bEq        = []byte("=")

	// Keys source array and raw address of it.
	bKeys = []byte("schemeslashesauthusernamepasswordhosthostnameportpathnamequeryoriginhashtruequery")

	errBadInit = errors.New("bad vector initialization, use urlvector.NewVector() or urlvector.Acquire()")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	if !vec.init {
		err = errBadInit
		return
	}

	s = bytealg.Trim(s, bSpace)
	// Remove blob: prefix if present.
	if len(s) >= 5 && bytes.Equal(s[:5], bBlob) {
		s = s[5:]
	}
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}
	vec.SetBit(flagCopy, copy)

	offset := 0
	// Create root node and register it.
	root, i := vec.AcquireNodeWithType(0, vector.TypeObject)
	root.SetOffset(vec.Index.Len(1))

	// Parse URL parts.
	if offset, err = vec.parseScheme(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}
	if offset, err = vec.parseAuth(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}
	if offset, err = vec.parseHost(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}
	if offset, err = vec.parsePath(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}
	if offset, err = vec.parseQuery(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}

	// Return root node back to the vector.
	vec.relNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		err = vector.ErrUnparsedTail
		return
	}

	return
}

// Parse schema and slashes parts.
func (vec *Vector) parseScheme(depth, offset int, node *vector.Node) (int, error) {
	var err error

	scheme, isc := vec.AcquireChildWithType(node, depth, vector.TypeString)
	slashes, isl := vec.AcquireChildWithType(node, depth, vector.TypeBool)

	limit := maxSchemaLen
	if sl := vec.SrcLen(); sl < limit {
		limit = sl
	}
	if limit < 2 {
		return offset, vector.ErrShortSrc
	}
	if pos := bytes.Index(vec.Src(), bSchemaSep); pos > 0 {
		scheme.Key().Init(bKeys, offsetScheme, lenScheme)
		scheme.Value().Init(vec.Src(), offset, pos)
		offset += pos + 3
	} else if bytes.Equal(vec.Src()[:2], bSlashes) {
		slashes.Key().Init(bKeys, offsetSlashes, lenSlashes)
		slashes.Value().Init(bKeys, offsetTrue, lenTrue)
		offset += 2
	}

	vec.relNode(isc, scheme)
	vec.relNode(isl, slashes)

	return offset, err
}

// Parse auth (username + password) and separate username and password parts.
func (vec *Vector) parseAuth(depth, offset int, node *vector.Node) (int, error) {
	var err error

	auth, ia := vec.AcquireChildWithType(node, depth, vector.TypeString)
	username, iu := vec.AcquireChildWithType(node, depth, vector.TypeString)
	password, ip := vec.AcquireChildWithType(node, depth, vector.TypeString)

	src := vec.Src()
	n := len(src)
	_ = src[n-1]

	posCol := bytealg.IndexByteAtBytes(src, ':', offset)
	posAt := bytealg.IndexByteAtBytes(src, '@', max_(posCol, offset))
	if posSl := bytealg.IndexByteAtBytes(src, '/', offset); posSl >= 0 && posSl < posAt {
		posAt = -1
	}

	if posAt > 0 {
		auth.Key().Init(bKeys, offsetAuth, lenAuth)
		auth.Value().Init(src, offset, posAt-offset)

		if posCol >= 0 {
			username.Key().Init(bKeys, offsetUsername, lenUsername)
			username.Value().Init(src, offset, posCol-offset)
			offset = posCol + 1

			password.Key().Init(bKeys, offsetPassword, lenPassword)
			password.Value().Init(src, offset, posAt-posCol-1)
		} else {
			username.Key().Init(bKeys, offsetUsername, lenUsername)
			username.Value().Init(src, offset, posAt-offset)
		}
		offset = posAt + 1
	}

	vec.relNode(ia, auth)
	vec.relNode(iu, username)
	vec.relNode(ip, password)

	return offset, err
}

// Parse host (hostname + port) and separate hostname and port parts.
func (vec *Vector) parseHost(depth, offset int, node *vector.Node) (int, error) {
	var err error

	host, ih := vec.AcquireChildWithType(node, depth, vector.TypeString)
	hostname, in := vec.AcquireChildWithType(node, depth, vector.TypeString)
	port, ip := vec.AcquireChildWithType(node, depth, vector.TypeNumber)

	src := vec.Src()
	n := len(src)
	_ = src[n-1]

	posSl := bytealg.IndexByteAtBytes(src, '/', offset)
	if posSl < 0 {
		if posBSl := bytealg.IndexByteAtBytes(src, '\\', offset); posBSl >= 0 {
			posSl = posBSl
		} else if posQM := bytealg.IndexByteAtBytes(src, '?', offset); posQM >= 0 {
			posSl = posQM
		} else {
			posSl = n
		}
	}
	posCol := -1
	i := offset
loop:
	i = bytealg.IndexByteAtBytes(src, ':', i+1)
	if i >= 0 && i < posSl {
		posCol = i
		goto loop
	}

	host.Key().Init(bKeys, offsetHost, lenHost)
	host.Value().Init(src, offset, posSl-offset)

	if posCol >= 0 {
		hostname.Key().Init(bKeys, offsetHostname, lenHostname)
		hostname.Value().Init(src, offset, posCol-offset)
		offset = posCol + 1

		port.Key().Init(bKeys, offsetPort, lenPort)
		port.Value().Init(src, offset, posSl-offset)
	} else {
		hostname.Key().Init(bKeys, offsetHostname, lenHostname)
		hostname.Value().Init(src, offset, posSl-offset)
	}

	vec.relNode(ih, host)
	vec.relNode(in, hostname)
	vec.relNode(ip, port)

	offset = posSl

	return offset, err
}

// Parse path part.
func (vec *Vector) parsePath(depth, offset int, node *vector.Node) (int, error) {
	var err error

	path, i := vec.AcquireChildWithType(node, depth, vector.TypeString)

	src := vec.Src()
	n := len(src)
	_ = src[n-1]

	if offset < n {
		posQM := bytealg.IndexByteAtBytes(vec.Src(), '?', offset)
		posHash := bytealg.IndexByteAtBytes(vec.Src(), '#', offset)
		if posQM >= 0 && posHash >= 0 && posQM > posHash {
			posQM = posHash
		}
		if posQM < 0 {
			if posHash >= 0 {
				posQM = posHash
			} else {
				posQM = n
			}
		}
		path.Key().Init(bKeys, offsetPath, lenPath)
		val := src[offset:posQM]
		path.Value().Init(src, offset, posQM-offset)
		path.Value().SetBit(flagEscape, bytealg.IndexByteAtBytes(val, '%', 0) >= 0)
		offset = posQM
	}

	vec.relNode(i, path)

	return offset, err
}

// Parse query part.
//
// Note that this method doesn't parse query to separate args for performance.
// Query will be parsed to separate args by first attempt to access to any query argument.
func (vec *Vector) parseQuery(depth, offset int, node *vector.Node) (int, error) {
	var err error

	queryOrig, iqo := vec.AcquireChildWithType(node, depth, vector.TypeString)
	hash, ih := vec.AcquireChildWithType(node, depth, vector.TypeString)
	query, iq := vec.AcquireChildWithType(node, depth, vector.TypeObject)

	src := vec.Src()
	n := len(src)
	_ = src[n-1]

	if offset < n {
		posHash := bytealg.IndexByteAtBytes(src, '#', offset)
		if posHash < 0 {
			posHash = n
		} else {
			hash.Key().Init(bKeys, offsetHash, lenHash)
			hash.Value().Init(src, posHash, n-posHash)
		}
		query.Key().Init(bKeys, offsetQuery, lenQuery)
		queryOrig.Key().Init(bKeys, offsetQueryOrigin, lenQueryOrigin)
		queryOrig.Value().Init(src, offset, posHash-offset)
		offset = n
	}

	vec.relNode(iqo, queryOrig)
	vec.relNode(ih, hash)
	vec.relNode(iq, query)

	return offset, err
}

// Parse query string to separate arguments.
func (vec *Vector) parseQueryParams(query *vector.Node) {
	origin := bytealg.TrimLeft(vec.QueryBytes(), bQM)
	if len(origin) == 0 {
		return
	}
	var (
		offset, idx int
		kv, k, v    []byte
		root, node  *vector.Node
	)
	if origin[0] == '&' {
		offset++
	}
	for {
		kv, k, v = nil, nil, nil

		i := bytealg.IndexByteAtBytes(origin, '&', offset)
		if i < 0 {
			i = len(origin)
		}
		kv = origin[offset:i]
		j := bytealg.IndexByteAtBytes(kv, '=', 0)
		if j < 0 {
			k = kv
		} else {
			k = kv[:j]
			if j < len(kv)-1 {
				v = kv[j+1:]
			}
		}
		if len(k) == 0 {
			offset = i + 1
			continue
		}

		if kl := len(k); kl > 2 && bytes.Equal(k[kl-2:], bQB) {
			if root = query.Get(byteconv.B2S(k)); root.Type() != vector.TypeArray {
				root, _ = vec.AcquireChildWithType(query, 2, vector.TypeArray)
				root.SetOffset(vec.Index.Len(3))
				root.Key().Init(origin, offset, len(k))
			}
			node, idx = vec.AcquireChildWithType(root, 3, vector.TypeString)
			if len(v) > 0 {
				v = unescape(v)
				node.Value().Init(origin, offset+len(k)+1, len(v))
			}
			vec.relNode(idx, node)
			vec.ReleaseNode(root.Index(), root)
		} else {
			node, idx = vec.AcquireChildWithType(query, 2, vector.TypeString)
			node.Key().Init(origin, offset, len(k))
			if len(v) > 0 {
				v = unescape(v)
				node.Value().Init(origin, offset+len(k)+1, len(v))
			}
			vec.relNode(idx, node)
		}

		offset = i + 1
		if offset >= len(origin) {
			break
		}
	}
	vec.relNode(query.Index(), query)
}

// Call vector.ReleaseNode() and set required flags before.
func (vec *Vector) relNode(idx int, node *vector.Node) {
	vec.ensureFlags(node)
	vec.ReleaseNode(idx, node)
}

// Consider source origin and set flags.
func (vec *Vector) ensureFlags(node *vector.Node) {
	if vec.CheckBit(flagCopy) {
		node.Value().SetBit(flagBufSrc, true)
	}
}

func max_(a, b int) int {
	if a > b {
		return a
	}
	return b
}
