package urlvector

// Alias of QueryEscape().
//
// Uses built-in buffer.
func (vec *Vector) QueryEscape(p []byte) []byte {
	return queryEscape(vec, p)
}

// Alias of QueryUnescape().
func (vec *Vector) QueryUnescape(p []byte) []byte {
	return queryUnescape(p)
}

// Alias of PathEscape().
//
// Uses built-in buffer.
func (vec *Vector) PathEscape(p []byte) []byte {
	return pathEscape(vec, p)
}

// Alias of PathUnescape().
func (vec *Vector) PathUnescape(p []byte) []byte {
	return pathUnescape(p)
}
