package urlvector

// SchemeBytes returns scheme as bytes.
func (vec *Vector) SchemeBytes() []byte {
	return vec.Scheme().Bytes()
}

// SchemeString returns scheme as string.
func (vec *Vector) SchemeString() string {
	return vec.Scheme().String()
}

// AuthBytes returns auth (username:password) as bytes.
func (vec *Vector) AuthBytes() []byte {
	return vec.Auth().Bytes()
}

// AuthString returns auth (username:password as string.
func (vec *Vector) AuthString() string {
	return vec.Auth().String()
}

// UsernameBytes returns username as bytes.
func (vec *Vector) UsernameBytes() []byte {
	return vec.Username().Bytes()
}

// UsernameString returns username as string.
func (vec *Vector) UsernameString() string {
	return vec.Username().String()
}

// PasswordBytes returns password as bytes.
func (vec *Vector) PasswordBytes() []byte {
	return vec.Password().Bytes()
}

// PasswordString returns password as string.
func (vec *Vector) PasswordString() string {
	return vec.Password().String()
}

// HostBytes returns host (hostname:port) as bytes.
func (vec *Vector) HostBytes() []byte {
	return vec.Host().Bytes()
}

// HostString returns host (hostname:port) as string.
func (vec *Vector) HostString() string {
	return vec.Host().String()
}

// HostnameBytes returns hostname as bytes.
func (vec *Vector) HostnameBytes() []byte {
	return vec.Hostname().Bytes()
}

// HostnameString returns hostname as string.
func (vec *Vector) HostnameString() string {
	return vec.Hostname().String()
}

// PathBytes returns password as bytes.
func (vec *Vector) PathBytes() []byte {
	return vec.Path().Bytes()
}

// PathString returns password as string.
func (vec *Vector) PathString() string {
	return vec.Path().String()
}

// QueryBytes returns query as bytes.
func (vec *Vector) QueryBytes() []byte {
	return vec.queryOrigin().Bytes()
}

// QueryString returns query as string.
func (vec *Vector) QueryString() string {
	return vec.queryOrigin().String()
}

// QueryLen returns length of the raw query without question mark symbol.
func (vec *Vector) QueryLen() int {
	if l := len(vec.QueryBytes()) - 1; l >= 0 {
		return l
	}
	return 0
}

// HashBytes returns hash as bytes.
func (vec *Vector) HashBytes() []byte {
	return vec.Hash().Bytes()
}

// HashString returns hash as string.
func (vec *Vector) HashString() string {
	return vec.Hash().String()
}
