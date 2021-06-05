package urlvector

import (
	"github.com/koykov/fastconv"
	"github.com/koykov/vector"
)

// Replace scheme with bytes.
func (vec *Vector) SetSchemeBytes(scheme []byte) *Vector {
	return vec.set(vec.Scheme(), scheme)
}

// Replace scheme with string.
func (vec *Vector) SetSchemeString(scheme string) *Vector {
	return vec.SetSchemeBytes(fastconv.S2B(scheme))
}

// Replace auth with bytes.
func (vec *Vector) SetAuthBytes(auth []byte) *Vector {
	return vec.set(vec.Auth(), auth)
}

// Replace auth with string.
func (vec *Vector) SetAuthString(auth string) *Vector {
	return vec.SetAuthBytes(fastconv.S2B(auth))
}

// Replace username with bytes.
func (vec *Vector) SetUsernameBytes(username []byte) *Vector {
	return vec.set(vec.Username(), username)
}

// Replace username with string.
func (vec *Vector) SetUsernameString(username string) *Vector {
	return vec.SetUsernameBytes(fastconv.S2B(username))
}

// Replace password with bytes.
func (vec *Vector) SetPasswordBytes(password []byte) *Vector {
	return vec.set(vec.Password(), password)
}

// Replace password with string.
func (vec *Vector) SetPasswordString(password string) *Vector {
	return vec.SetPasswordBytes(fastconv.S2B(password))
}

// Replace host with bytes.
func (vec *Vector) SetHostBytes(host []byte) *Vector {
	return vec.set(vec.Host(), host)
}

// Replace host with string.
func (vec *Vector) SetHostString(host string) *Vector {
	return vec.SetHostBytes(fastconv.S2B(host))
}

// Replace hostname with bytes.
func (vec *Vector) SetHostnameBytes(hostname []byte) *Vector {
	return vec.set(vec.Hostname(), hostname)
}

// Replace hostname with string.
func (vec *Vector) SetHostnameString(hostname string) *Vector {
	return vec.SetHostnameBytes(fastconv.S2B(hostname))
}

// Replace path with bytes.
func (vec *Vector) SetPathBytes(path []byte) *Vector {
	return vec.set(vec.Path(), path)
}

// Replace path with string.
func (vec *Vector) SetPathString(path string) *Vector {
	return vec.SetPathBytes(fastconv.S2B(path))
}

// Replace query with bytes.
func (vec *Vector) SetQueryBytes(query []byte) *Vector {
	vec.SetBit(flagQueryParsed, false)
	vec.ForgetFrom(idxQuery + 1)
	node := vec.GetByIdx(idxQuery)
	vec.Index.Reset(node.Depth(), node.Offset())
	node.SetLimit(0)
	return vec.set(vec.queryOrigin(), query)
}

// Replace query with string.
func (vec *Vector) SetQueryString(query string) *Vector {
	return vec.SetQueryBytes(fastconv.S2B(query))
}

// Replace hash with bytes.
func (vec *Vector) SetHashBytes(hash []byte) *Vector {
	return vec.set(vec.Hash(), hash)
}

// Replace hash with string.
func (vec *Vector) SetHashString(hash string) *Vector {
	return vec.SetHashBytes(fastconv.S2B(hash))
}

// Internal setter.
func (vec *Vector) set(node *vector.Node, s []byte) *Vector {
	vec.SetBit(flagBufMod, true)
	offset := vec.BufLen()
	vec.BufAppend(s)
	node.Value().Init(vec.Buf(), offset, len(s))
	node.Value().SetBit(flagBufSrc, true)
	return vec
}
