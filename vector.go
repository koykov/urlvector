package urlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

type Vector struct {
	vector.Vector
	keyAddr uint64
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
	return vec.Get("scheme")
}

func (vec *Vector) Slashes() bool {
	return vec.Get("slashes").Bool()
}

func (vec *Vector) Auth() *vector.Node {
	return vec.Get("auth")
}

func (vec *Vector) Username() *vector.Node {
	return vec.Get("username")
}

func (vec *Vector) Password() *vector.Node {
	return vec.Get("password")
}

func (vec *Vector) Host() *vector.Node {
	return vec.Get("host")
}

func (vec *Vector) Hostname() *vector.Node {
	return vec.Get("hostname")
}

func (vec *Vector) Port() int {
	i, _ := vec.Get("port").Int()
	return int(i)
}

func (vec *Vector) Path() *vector.Node {
	return vec.Get("path")
}

func (vec *Vector) Query() *vector.Node {
	return vec.Get("query")
}

func (vec *Vector) Hash() *vector.Node {
	return vec.Get("hash")
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.keyAddr = 0
}
