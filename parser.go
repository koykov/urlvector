package urlvector

import (
	"bytes"
	"reflect"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/vector"
)

const (
	offsetScheme   = 0
	offsetSlashes  = 6
	offsetAuth     = 13
	offsetUsername = 17
	offsetPassword = 25
	offsetHost     = 33
	offsetHostname = 37
	offsetPort     = 45
	offsetPathname = 49
	offsetQuery    = 57
	offsetHash     = 62
	offsetHref     = 66
	offsetTrue     = 70

	lenScheme   = 6
	lenSlashes  = 7
	lenAuth     = 4
	lenUsername = 8
	lenPassword = 8
	lenHost     = 4
	lenHostname = 8
	lenPort     = 4
	lenPathname = 8
	lenQuery    = 5
	lenHash     = 4
	lenHref     = 4
	lenTrue     = 4
)

var (
	bSpace     = []byte(" ")
	bSchemaSep = []byte("://")
	bSlashes   = []byte("//")
	bSlash     = []byte("/")
	bColon     = []byte(":")
	bAt        = []byte("@")

	bIndex = []byte("schemeslashesauthusernamepasswordhosthostnameportpathnamequeryhashhreftrue")
)

func (vec *Vector) parse(s []byte, copy bool) (err error) {
	s = bytealg.Trim(s, bSpace)
	if err = vec.SetSrc(s, copy); err != nil {
		return
	}

	h := (*reflect.SliceHeader)(unsafe.Pointer(&bIndex))
	vec.keyAddr = uint64(h.Data)

	offset := 0
	// Create root node and register it.
	root := vec.GetNodeWT(0, vector.TypeObj)
	i := vec.Len() - 1
	vec.Index.Register(0, i)
	root.SetOffset(vec.Index.Len(1))

	if offset, err = vec.parseScheme(1, offset, root); err != nil {
		vec.SetErrOffset(offset)
		return
	}
	if offset, err = vec.parseAuth(1, offset, root); err != nil {
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

	// max scheme length by https://en.wikipedia.org/wiki/List_of_URI_schemes#Official_IANA-registered_schemes
	limit := 29
	if sl := vec.SrcLen(); sl < limit {
		limit = sl
	}
	if limit < 2 {
		return offset, vector.ErrShortSrc
	}
	if pos := bytes.Index(vec.Src(), bSchemaSep); pos > 0 {
		scheme := vec.GetNodeWT(depth, vector.TypeStr)
		i := vec.Len() - 1
		node.SetLimit(vec.Index.Register(depth, i))
		scheme.Key().Set(vec.keyAddr+offsetScheme, lenScheme)
		scheme.Value().Set(vec.SrcAddr()+uint64(offset), pos)
		vec.PutNode(i, scheme)
		offset += pos + 3
	} else if bytes.Equal(vec.Src()[:2], bSlashes) {
		scheme := vec.GetNodeWT(depth, vector.TypeBool)
		i := vec.Len() - 1
		node.SetLimit(vec.Index.Register(depth, i))
		scheme.Key().Set(vec.keyAddr+offsetSlashes, lenSlashes)
		scheme.Value().Set(vec.keyAddr+offsetTrue, lenTrue)
		vec.PutNode(i, scheme)
		offset += 2
	}

	return offset, err
}

func (vec *Vector) parseAuth(depth, offset int, node *vector.Node) (int, error) {
	var err error
	posCol := bytealg.IndexAt(vec.Src(), bColon, offset)
	if posCol < 0 {
		return offset, nil
	}
	posAt := bytealg.IndexAt(vec.Src(), bAt, posCol)
	if posAt < 0 || posAt <= posCol {
		return offset, nil
	}

	auth := vec.GetNodeWT(depth, vector.TypeStr)
	i := vec.Len() - 1
	node.SetLimit(vec.Index.Register(depth, i))
	auth.Key().Set(vec.keyAddr+offsetAuth, lenAuth)
	auth.Value().Set(vec.SrcAddr()+uint64(offset), posAt-offset)
	vec.PutNode(i, auth)

	username := vec.GetNodeWT(depth, vector.TypeStr)
	i = vec.Len() - 1
	node.SetLimit(vec.Index.Register(depth, i))
	username.Key().Set(vec.keyAddr+offsetUsername, lenUsername)
	username.Value().Set(vec.SrcAddr()+uint64(offset), posAt-offset)
	vec.PutNode(i, username)

	password := vec.GetNodeWT(depth, vector.TypeStr)
	i = vec.Len() - 1
	node.SetLimit(vec.Index.Register(depth, i))
	password.Key().Set(vec.keyAddr+offsetPassword, lenPassword)
	password.Value().Set(vec.SrcAddr()+uint64(offset+posCol), posAt-posCol-offset)
	vec.PutNode(i, password)

	offset += posAt

	return offset, err
}
