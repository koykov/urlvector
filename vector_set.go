package urlvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

// SetSchemeBytes replaces scheme with bytes.
func (vec *Vector) SetSchemeBytes(scheme []byte) *Vector {
	return vec.set(vec.Scheme(), scheme)
}

// SetSchemeString replaces scheme with string.
func (vec *Vector) SetSchemeString(scheme string) *Vector {
	return vec.SetSchemeBytes(byteconv.S2B(scheme))
}

// SetAuthBytes replaces auth with bytes.
func (vec *Vector) SetAuthBytes(auth []byte) *Vector {
	return vec.set(vec.Auth(), auth)
}

// SetAuthString replaces auth with string.
func (vec *Vector) SetAuthString(auth string) *Vector {
	return vec.SetAuthBytes(byteconv.S2B(auth))
}

// SetUsernameBytes replaces username with bytes.
func (vec *Vector) SetUsernameBytes(username []byte) *Vector {
	return vec.set(vec.Username(), username)
}

// SetUsernameString replaces username with string.
func (vec *Vector) SetUsernameString(username string) *Vector {
	return vec.SetUsernameBytes(byteconv.S2B(username))
}

// SetPasswordBytes replaces password with bytes.
func (vec *Vector) SetPasswordBytes(password []byte) *Vector {
	return vec.set(vec.Password(), password)
}

// SetPasswordString replaces password with string.
func (vec *Vector) SetPasswordString(password string) *Vector {
	return vec.SetPasswordBytes(byteconv.S2B(password))
}

// SetHostBytes replaces host with bytes.
func (vec *Vector) SetHostBytes(host []byte) *Vector {
	return vec.set(vec.Host(), host)
}

// SetHostString replaces host with string.
func (vec *Vector) SetHostString(host string) *Vector {
	return vec.SetHostBytes(byteconv.S2B(host))
}

// SetHostnameBytes replaces hostname with bytes.
func (vec *Vector) SetHostnameBytes(hostname []byte) *Vector {
	return vec.set(vec.Hostname(), hostname)
}

// SetHostnameString replaces hostname with string.
func (vec *Vector) SetHostnameString(hostname string) *Vector {
	return vec.SetHostnameBytes(byteconv.S2B(hostname))
}

// SetPort replaces port.
func (vec *Vector) SetPort(port int) *Vector {
	vec.SetBit(flagBufMod, true)
	offset := vec.BufLen()
	vec.BufAppendInt(int64(port))
	l := vec.BufLen() - offset
	node := vec.getByIdx(idxPort)
	node.Value().Init(vec.Buf(), offset, l)
	node.Value().SetBit(flagBufSrc, true)
	return vec
}

// SetPathBytes replaces path with bytes.
func (vec *Vector) SetPathBytes(path []byte) *Vector {
	return vec.set(vec.Path(), path)
}

// SetPathString replaces path with string.
func (vec *Vector) SetPathString(path string) *Vector {
	return vec.SetPathBytes(byteconv.S2B(path))
}

// SetQueryBytes replaces query with bytes.
func (vec *Vector) SetQueryBytes(query []byte) *Vector {
	vec.SetBit(flagQueryParsed, false)
	vec.ForgetFrom(idxQuery + 1)
	node := vec.GetByIdx(idxQuery)
	vec.Index.Reset(node.Depth(), node.Offset())
	node.SetLimit(0)
	return vec.set(vec.queryOrigin(), query)
}

// SetQueryString replaces query with string.
func (vec *Vector) SetQueryString(query string) *Vector {
	return vec.SetQueryBytes(byteconv.S2B(query))
}

// SetHashBytes replaces hash with bytes.
func (vec *Vector) SetHashBytes(hash []byte) *Vector {
	return vec.set(vec.Hash(), hash)
}

// SetHashString replaces hash with string.
func (vec *Vector) SetHashString(hash string) *Vector {
	return vec.SetHashBytes(byteconv.S2B(hash))
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
