package urlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

const (
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
)

type Vector struct {
	vector.Vector
	keyAddr     uint64
	queryParsed bool
}

func NewVector() *Vector {
	vec := &Vector{}
	return vec
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseStr(s string) error {
	return vec.parse(fastconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyStr(s string) error {
	return vec.parse(fastconv.S2B(s), true)
}

func (vec *Vector) Scheme() *vector.Node {
	return vec.GetByIdx(idxScheme)
}

func (vec *Vector) Slashes() bool {
	return vec.GetByIdx(idxSlashes).Bool()
}

func (vec *Vector) Auth() *vector.Node {
	return vec.GetByIdx(idxAuth)
}

func (vec *Vector) Username() *vector.Node {
	return vec.GetByIdx(idxUsername)
}

func (vec *Vector) Password() *vector.Node {
	return vec.GetByIdx(idxPassword)
}

func (vec *Vector) Host() *vector.Node {
	return vec.GetByIdx(idxHost)
}

func (vec *Vector) Hostname() *vector.Node {
	return vec.GetByIdx(idxHostname)
}

func (vec *Vector) Port() int {
	i, _ := vec.GetByIdx(idxPort).Int()
	return int(i)
}

func (vec *Vector) Path() *vector.Node {
	return vec.GetByIdx(idxPath)
}

func (vec *Vector) Query() *vector.Node {
	query := vec.GetByIdx(idxQuery)
	if !vec.queryParsed {
		vec.queryParsed = true
		vec.parseQueryParams(query)
	}
	return query
}

func (vec *Vector) queryOrigin() *vector.Node {
	return vec.GetByIdx(idxQueryOrigin)
}

func (vec *Vector) Hash() *vector.Node {
	return vec.GetByIdx(idxHash)
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.keyAddr = 0
	vec.queryParsed = false
}
