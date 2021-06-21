package urlvector

// Alias of QueryEscape().
//
// Uses built-in buffer.
func (vec *Vector) QueryEscape(p []byte) []byte {
	return vecEscape(vec, p, modeQuery)
}

// Alias of QueryUnescape().
func (vec *Vector) QueryUnescape(p []byte) []byte {
	return vecUnescape(vec, p, modeQuery)
}

// Alias of PathEscape().
//
// Uses built-in buffer.
func (vec *Vector) PathEscape(p []byte) []byte {
	return vecEscape(vec, p, modePath)
}

// Alias of PathUnescape().
func (vec *Vector) PathUnescape(p []byte) []byte {
	return vecUnescape(vec, p, modePath)
}
