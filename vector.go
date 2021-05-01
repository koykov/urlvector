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

func (vec *Vector) Scheme() []byte {
	return vec.Get("scheme").Bytes()
}

func (vec *Vector) SchemeStr() string {
	return vec.Get("scheme").String()
}

func (vec *Vector) Slashes() bool {
	return vec.Get("slashed").Bool()
}

func (vec *Vector) Auth() []byte {
	return vec.Get("auth").Bytes()
}

func (vec *Vector) AuthStr() string {
	return vec.Get("auth").String()
}

func (vec *Vector) Username() []byte {
	return vec.Get("username").Bytes()
}

func (vec *Vector) UsernameStr() string {
	return vec.Get("username").String()
}

func (vec *Vector) Password() []byte {
	return vec.Get("password").Bytes()
}

func (vec *Vector) PasswordStr() string {
	return vec.Get("password").String()
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.keyAddr = 0
}
