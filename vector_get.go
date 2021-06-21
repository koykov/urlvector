package urlvector

// Get scheme as bytes.
func (vec *Vector) SchemeBytes() []byte {
	return vec.Scheme().Bytes()
}

// Get scheme as string.
func (vec *Vector) SchemeString() string {
	return vec.Scheme().String()
}

// Get auth (username:password) as bytes.
func (vec *Vector) AuthBytes() []byte {
	return vec.Auth().Bytes()
}

// Get auth (username:password as string.
func (vec *Vector) AuthString() string {
	return vec.Auth().String()
}

// Get username as bytes.
func (vec *Vector) UsernameBytes() []byte {
	return vec.Username().Bytes()
}

// Get username as string.
func (vec *Vector) UsernameString() string {
	return vec.Username().String()
}

// Get password as bytes.
func (vec *Vector) PasswordBytes() []byte {
	return vec.Password().Bytes()
}

// Get password as string.
func (vec *Vector) PasswordString() string {
	return vec.Password().String()
}

// Get host (hostname:port) as bytes.
func (vec *Vector) HostBytes() []byte {
	return vec.Host().Bytes()
}

// Get host (hostname:port) as string.
func (vec *Vector) HostString() string {
	return vec.Host().String()
}

// Get hostname as bytes.
func (vec *Vector) HostnameBytes() []byte {
	return vec.Hostname().Bytes()
}

// Get hostname as string.
func (vec *Vector) HostnameString() string {
	return vec.Hostname().String()
}

// Get password as bytes.
func (vec *Vector) PathBytes() []byte {
	return vec.Path().Bytes()
}

// Get password as string.
func (vec *Vector) PathString() string {
	return vec.Path().String()
}

// Get query as bytes.
func (vec *Vector) QueryBytes() []byte {
	return vec.queryOrigin().Bytes()
}

// Get query as string.
func (vec *Vector) QueryString() string {
	return vec.queryOrigin().String()
}

// Get length of the raw query without question mark symbol.
func (vec *Vector) QueryLen() int {
	if l := len(vec.QueryBytes()) - 1; l >= 0 {
		return l
	}
	return 0
}

// Get hash as bytes.
func (vec *Vector) HashBytes() []byte {
	return vec.Hash().Bytes()
}

// Get hash as string.
func (vec *Vector) HashString() string {
	return vec.Hash().String()
}
