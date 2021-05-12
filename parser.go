package urlvector

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
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
)

var (
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

	bIndex = []byte("schemeslashesauthusernamepasswordhosthostnameportpathnamequeryoriginhashtruequery")
)

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bSpace)
	if len(s) >= 5 && bytes.Equal(s[:5], bBlob) {
		s = s[5:]
	}
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	h := (*reflect.SliceHeader)(unsafe.Pointer(&bIndex))
	vec.keyAddr = uint64(h.Data)

	offset := 0
	// Create root node and register it.
	root, i := vec.GetNodeWT(0, vector.TypeObj)
	root.SetOffset(vec.Index.Len(1))

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

	vec.PutNode(i, root)

	// Check unparsed tail.
	if offset < vec.SrcLen() {
		vec.SetErrOffset(offset)
		err = vector.ErrUnparsedTail
		return
	}

	return
}

func (vec *Vector) parseScheme(depth, offset int, node *vector.Node) (int, error) {
	var err error

	scheme, isc := vec.GetChildWT(node, depth, vector.TypeStr)
	slashes, isl := vec.GetChildWT(node, depth, vector.TypeBool)

	// max scheme length by https://en.wikipedia.org/wiki/List_of_URI_schemes#Official_IANA-registered_schemes
	limit := 29
	if sl := vec.SrcLen(); sl < limit {
		limit = sl
	}
	if limit < 2 {
		return offset, vector.ErrShortSrc
	}
	if pos := bytes.Index(vec.Src(), bSchemaSep); pos > 0 {
		scheme.Key().Set(vec.keyAddr+offsetScheme, lenScheme)
		scheme.Value().Set(vec.SrcAddr()+uint64(offset), pos)
		offset += pos + 3
	} else if bytes.Equal(vec.Src()[:2], bSlashes) {
		slashes.Key().Set(vec.keyAddr+offsetSlashes, lenSlashes)
		slashes.Value().Set(vec.keyAddr+offsetTrue, lenTrue)
		offset += 2
	}

	vec.PutNode(isc, scheme)
	vec.PutNode(isl, slashes)

	return offset, err
}

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
		auth.Key().Set(vec.keyAddr+offsetAuth, lenAuth)
		auth.Value().Set(vec.SrcAddr()+uint64(offset), posAt-offset)

		if posCol >= 0 {
			username.Key().Set(vec.keyAddr+offsetUsername, lenUsername)
			username.Value().Set(vec.SrcAddr()+uint64(offset), posCol-offset)
			offset = posCol + 1

			password.Key().Set(vec.keyAddr+offsetPassword, lenPassword)
			password.Value().Set(vec.SrcAddr()+uint64(offset), posAt-posCol-1)
		} else {
			username.Key().Set(vec.keyAddr+offsetUsername, lenUsername)
			username.Value().Set(vec.SrcAddr()+uint64(offset), posAt-offset)
		}
		offset = posAt + 1
	}

	vec.PutNode(ia, auth)
	vec.PutNode(iu, username)
	vec.PutNode(ip, password)

	return offset, err
}

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

	host.Key().Set(vec.keyAddr+offsetHost, lenHost)
	host.Value().Set(vec.SrcAddr()+uint64(offset), posSl-offset)

	if posCol >= 0 {
		hostname.Key().Set(vec.keyAddr+offsetHostname, lenHostname)
		hostname.Value().Set(vec.SrcAddr()+uint64(offset), posCol-offset)
		offset = posCol + 1

		port.Key().Set(vec.keyAddr+offsetPort, lenPort)
		port.Value().Set(vec.SrcAddr()+uint64(offset), posSl-offset)
	} else {
		hostname.Key().Set(vec.keyAddr+offsetHostname, lenHostname)
		hostname.Value().Set(vec.SrcAddr()+uint64(offset), posSl-offset)
	}

	vec.PutNode(ih, host)
	vec.PutNode(in, hostname)
	vec.PutNode(ip, port)

	offset = posSl

	return offset, err
}

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
		path.Key().Set(vec.keyAddr+offsetPath, lenPath)
		path.Value().Set(vec.SrcAddr()+uint64(offset), posQM-offset)
		offset = posQM
	}

	vec.PutNode(i, path)

	return offset, err
}

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
			hash.Key().Set(vec.keyAddr+offsetHash, lenHash)
			hash.Value().Set(vec.SrcAddr()+uint64(posHash), vec.SrcLen()-posHash)
		}
		query.Key().Set(vec.keyAddr+offsetQuery, lenQuery)
		queryOrig.Key().Set(vec.keyAddr+offsetQueryOrigin, lenQueryOrigin)
		queryOrig.Value().Set(vec.SrcAddr()+uint64(offset), posHash-offset)
		offset = vec.SrcLen()
	}

	vec.PutNode(iqo, queryOrig)
	vec.PutNode(ih, hash)
	vec.PutNode(iq, query)

	return offset, err
}

func (vec *Vector) parseQueryParams(query *vector.Node) {
	origin := vec.QueryBytes()
	if len(origin) == 0 {
		return
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&origin))
	originAddr := uint64(h.Data)
	var (
		node   *vector.Node
		idx, i int
		offset uint64
		k      = true
	)
	node, idx = vec.GetChildWT(query, 2, vector.TypeStr)
	_ = origin[len(origin)-1]
	for i = 0; i < len(origin); i++ {
		switch origin[i] {
		case '?':
			offset = uint64(i) + 1
		case '&':
			k = true
			node.Value().Set(originAddr+offset, i-int(offset))
			offset = uint64(i) + 1
			vec.PutNode(idx, node)

			node, idx = vec.GetChildWT(query, 2, vector.TypeStr)
		case '=':
			if k {
				k = false
				node.Key().Set(originAddr+offset, i-int(offset))
				offset = uint64(i) + 1
			}
		}
	}
	node.Value().Set(originAddr+offset, i-int(offset))
	vec.PutNode(idx, node)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
