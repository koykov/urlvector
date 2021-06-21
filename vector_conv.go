package urlvector

// Alias of QueryEscape().
//
// Uses built-in buffer.
func (vec *Vector) QueryEscape(p []byte) []byte {
	return escape(vec, p)
}

// Alias of QueryUnescape().
func (vec *Vector) QueryUnescape(p []byte) []byte {
	return unescape(p)
}
