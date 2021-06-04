package urlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

func (vec *Vector) SetSchemeBytes(scheme []byte) *Vector {
	return vec.set(vec.Scheme(), scheme)
}

func (vec *Vector) SetSchemeString(scheme string) *Vector {
	return vec.SetSchemeBytes(fastconv.S2B(scheme))
}

func (vec *Vector) SetAuthBytes(auth []byte) *Vector {
	return vec.set(vec.Auth(), auth)
}

func (vec *Vector) SetAuthString(auth string) *Vector {
	return vec.SetAuthBytes(fastconv.S2B(auth))
}

func (vec *Vector) SetUsernameBytes(username []byte) *Vector {
	return vec.set(vec.Username(), username)
}

func (vec *Vector) SetUsernameString(username string) *Vector {
	return vec.SetUsernameBytes(fastconv.S2B(username))
}

func (vec *Vector) SetPasswordBytes(password []byte) *Vector {
	return vec.set(vec.Password(), password)
}

func (vec *Vector) SetPasswordString(password string) *Vector {
	return vec.SetPasswordBytes(fastconv.S2B(password))
}

func (vec *Vector) SetHostBytes(host []byte) *Vector {
	return vec.set(vec.Host(), host)
}

func (vec *Vector) SetHostString(host string) *Vector {
	return vec.SetHostBytes(fastconv.S2B(host))
}

func (vec *Vector) SetHostnameBytes(hostname []byte) *Vector {
	return vec.set(vec.Hostname(), hostname)
}

func (vec *Vector) SetHostnameString(hostname string) *Vector {
	return vec.SetHostnameBytes(fastconv.S2B(hostname))
}

func (vec *Vector) SetPathBytes(path []byte) *Vector {
	return vec.set(vec.Path(), path)
}

func (vec *Vector) SetPathString(path string) *Vector {
	return vec.SetPathBytes(fastconv.S2B(path))
}

func (vec *Vector) SetQueryBytes(query []byte) *Vector {
	vec.SetBit(flagQueryParsed, false)
	vec.ForgetFrom(idxQuery + 1)
	vec.GetByIdx(idxQuery).ResetIndex()
	return vec.set(vec.queryOrigin(), query)
}

func (vec *Vector) SetQueryString(query string) *Vector {
	return vec.SetQueryBytes(fastconv.S2B(query))
}

func (vec *Vector) SetHashBytes(hash []byte) *Vector {
	return vec.set(vec.Hash(), hash)
}

func (vec *Vector) SetHashString(hash string) *Vector {
	return vec.SetHashBytes(fastconv.S2B(hash))
}

func (vec *Vector) set(node *vector.Node, s []byte) *Vector {
	vec.SetBit(flagBufMod, true)
	offset := vec.BufLen()
	vec.BufAppend(s)
	node.Value().Init(vec.Buf(), offset, len(s))
	node.Value().SetBit(flagBufSrc, true)
	return vec
}
