package urlvector

import (
	"bytes"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
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
	bBSlash    = []byte("\\")
	bColon     = []byte(":")
	bAt        = []byte("@")
	bQM        = []byte("?")
	bHash      = []byte("#")
	bQB        = []byte("[]")

	// Keys source array and raw address of it.
	bKeys = []byte("schemeslashesauthusernamepasswordhosthostnameportpathnamequeryoriginhashtruequery")
)

// Main internal parser helper.
func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bSpace)
	// Remove blob: prefix if present.
	if len(s) >= 5 && bytes.Equal(s[:5], bBlob) {
		s = s[5:]
	}
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	offset := 0
	// Create root node and register it.
	root, i := vec.GetNodeWT(0, vector.TypeObj)
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
	vec.PutNode(i, root)

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

	scheme, isc := vec.GetChildWT(node, depth, vector.TypeStr)
	slashes, isl := vec.GetChildWT(node, depth, vector.TypeBool)

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

	vec.PutNode(isc, scheme)
	vec.PutNode(isl, slashes)

	return offset, err
}

// Parse auth (username + password) and separate username and password parts.
func (vec *Vector) parseAuth(depth, offset int, node *vector.Node) (int, error) {
	var err error

	auth, ia := vec.GetChildWT(node, depth, vector.TypeStr)
	username, iu := vec.GetChildWT(node, depth, vector.TypeStr)
	password, ip := vec.GetChildWT(node, depth, vector.TypeStr)

	posCol := bytealg.IndexAt(vec.Src(), bColon, offset)
	posAt := bytealg.IndexAt(vec.Src(), bAt, max(posCol, offset))
	if posSl := bytealg.IndexAt(vec.Src(), bSlash, offset); posSl >= 0 && posSl < posAt {
		posAt = -1
	}

	if posAt > 0 {
		auth.Key().Init(bKeys, offsetAuth, lenAuth)
		auth.Value().Init(vec.Src(), offset, posAt-offset)

		if posCol >= 0 {
			username.Key().Init(bKeys, offsetUsername, lenUsername)
			username.Value().Init(vec.Src(), offset, posCol-offset)
			offset = posCol + 1

			password.Key().Init(bKeys, offsetPassword, lenPassword)
			password.Value().Init(vec.Src(), offset, posAt-posCol-1)
		} else {
			username.Key().Init(bKeys, offsetUsername, lenUsername)
			username.Value().Init(vec.Src(), offset, posAt-offset)
		}
		offset = posAt + 1
	}

	vec.PutNode(ia, auth)
	vec.PutNode(iu, username)
	vec.PutNode(ip, password)

	return offset, err
}

// Parse host (hostname + port) and separate hostname and port parts.
func (vec *Vector) parseHost(depth, offset int, node *vector.Node) (int, error) {
	var err error

	host, ih := vec.GetChildWT(node, depth, vector.TypeStr)
	hostname, in := vec.GetChildWT(node, depth, vector.TypeStr)
	port, ip := vec.GetChildWT(node, depth, vector.TypeNum)

	posSl := bytealg.IndexAt(vec.Src(), bSlash, offset)
	if posSl < 0 {
		if posBSl := bytealg.IndexAt(vec.Src(), bBSlash, offset); posBSl >= 0 {
			posSl = posBSl
		} else {
			posSl = vec.SrcLen()
		}
	}
	posCol := -1
	i := offset
loop:
	i = bytealg.IndexAt(vec.Src(), bColon, i+1)
	if i >= 0 && i < posSl {
		posCol = i
		goto loop
	}

	host.Key().Init(bKeys, offsetHost, lenHost)
	host.Value().Init(vec.Src(), offset, posSl-offset)

	if posCol >= 0 {
		hostname.Key().Init(bKeys, offsetHostname, lenHostname)
		hostname.Value().Init(vec.Src(), offset, posCol-offset)
		offset = posCol + 1

		port.Key().Init(bKeys, offsetPort, lenPort)
		port.Value().Init(vec.Src(), offset, posSl-offset)
	} else {
		hostname.Key().Init(bKeys, offsetHostname, lenHostname)
		hostname.Value().Init(vec.Src(), offset, posSl-offset)
	}

	vec.PutNode(ih, host)
	vec.PutNode(in, hostname)
	vec.PutNode(ip, port)

	offset = posSl

	return offset, err
}

// Parse path part.
func (vec *Vector) parsePath(depth, offset int, node *vector.Node) (int, error) {
	var err error

	path, i := vec.GetChildWT(node, depth, vector.TypeStr)

	if offset < vec.SrcLen() {
		posQM := bytealg.IndexAt(vec.Src(), bQM, offset)
		posHash := bytealg.IndexAt(vec.Src(), bHash, offset)
		if posQM < 0 {
			if posHash >= 0 {
				posQM = posHash
			} else {
				posQM = vec.SrcLen()
			}
		}
		path.Key().Init(bKeys, offsetPath, lenPath)
		path.Value().Init(vec.Src(), offset, posQM-offset)
		offset = posQM
	}

	vec.PutNode(i, path)

	return offset, err
}

// Parse query part.
//
// Note that this method doesn't parse query to separate args for performance.
// Query will be parsed to separate args by first attempt to access to any query argument.
func (vec *Vector) parseQuery(depth, offset int, node *vector.Node) (int, error) {
	var err error

	queryOrig, iqo := vec.GetChildWT(node, depth, vector.TypeStr)
	hash, ih := vec.GetChildWT(node, depth, vector.TypeStr)
	query, iq := vec.GetChildWT(node, depth, vector.TypeObj)

	if offset < vec.SrcLen() {
		posHash := bytealg.IndexAt(vec.Src(), bHash, offset)
		if posHash < 0 {
			posHash = vec.SrcLen()
		} else {
			hash.Key().Init(bKeys, offsetHash, lenHash)
			hash.Value().Init(vec.Src(), posHash, vec.SrcLen()-posHash)
		}
		query.Key().Init(bKeys, offsetQuery, lenQuery)
		queryOrig.Key().Init(bKeys, offsetQueryOrigin, lenQueryOrigin)
		queryOrig.Value().Init(vec.Src(), offset, posHash-offset)
		offset = vec.SrcLen()
	}

	vec.PutNode(iqo, queryOrig)
	vec.PutNode(ih, hash)
	vec.PutNode(iq, query)

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
	for {
		kv, k, v = nil, nil, nil

		i := bytealg.IndexByteAtLR(origin, '&', offset)
		if i < 0 {
			i = len(origin)
		}
		kv = origin[offset:i]
		j := bytealg.IndexByteAtLR(kv, '=', 0)
		if j < 0 {
			k = kv
		} else {
			k = kv[:j]
			if j < len(kv)-1 {
				v = kv[j+1:]
			}
		}
		if len(k) == 0 {
			continue
		}

		if kl := len(k); kl > 2 && bytes.Equal(k[kl-2:], bQB) {
			if root = query.Get(fastconv.B2S(k)); root.Type() != vector.TypeArr {
				root, _ = vec.GetChildWT(query, 2, vector.TypeArr)
				root.SetOffset(vec.Index.Len(3))
				root.Key().Init(origin, offset, len(k))
			}
			node, idx = vec.GetChildWT(root, 3, vector.TypeStr)
			if len(v) > 0 {
				node.Value().Init(origin, offset+len(k)+1, len(v))
				node.Value().SetBit(flagEscape, bytealg.IndexByteAtLR(v, '%', 0) >= 0)
			}
			vec.PutNode(idx, node)
			vec.PutNode(root.Index(), root)
		} else {
			node, idx = vec.GetChildWT(query, 2, vector.TypeStr)
			node.Key().Init(origin, offset, len(k))
			if len(v) > 0 {
				node.Value().Init(origin, offset+len(k)+1, len(v))
				node.Value().SetBit(flagEscape, bytealg.IndexByteAtLR(v, '%', 0) >= 0)
			}
			vec.PutNode(idx, node)
		}

		offset = i + 1
		if offset >= len(origin) {
			break
		}
	}
	vec.PutNode(query.Index(), query)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
