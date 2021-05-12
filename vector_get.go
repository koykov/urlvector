package urlvector

func (vec *Vector) SchemeBytes() []byte {
	return vec.Scheme().Bytes()
}

func (vec *Vector) SchemeString() string {
	return vec.Scheme().String()
}

func (vec *Vector) AuthBytes() []byte {
	return vec.Auth().Bytes()
}

func (vec *Vector) AuthString() string {
	return vec.Auth().String()
}

func (vec *Vector) UsernameBytes() []byte {
	return vec.Username().Bytes()
}

func (vec *Vector) UsernameString() string {
	return vec.Username().String()
}

func (vec *Vector) PasswordBytes() []byte {
	return vec.Password().Bytes()
}

func (vec *Vector) PasswordString() string {
	return vec.Password().String()
}

func (vec *Vector) HostBytes() []byte {
	return vec.Host().Bytes()
}

func (vec *Vector) HostString() string {
	return vec.Host().String()
}

func (vec *Vector) HostnameBytes() []byte {
	return vec.Hostname().Bytes()
}

func (vec *Vector) HostnameString() string {
	return vec.Hostname().String()
}

func (vec *Vector) PathBytes() []byte {
	return vec.Path().Bytes()
}

func (vec *Vector) PathString() string {
	return vec.Path().String()
}

func (vec *Vector) QueryBytes() []byte {
	return vec.queryOrigin().Bytes()
}

func (vec *Vector) QueryString() string {
	return vec.queryOrigin().String()
}

func (vec *Vector) HashBytes() []byte {
	return vec.Hash().Bytes()
}

func (vec *Vector) HashString() string {
	return vec.Hash().String()
}
