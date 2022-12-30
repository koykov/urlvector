package urlvector

// QueryEscape escapes the string so it can be safely placed inside a URL query.
//
// Uses built-in buffer.
func (vec *Vector) QueryEscape(p []byte) []byte {
	return vecEscape(vec, p, modeQuery)
}

// QueryUnescape does the inverse transformation of QueryEscape.
func (vec *Vector) QueryUnescape(p []byte) []byte {
	return vecUnescape(vec, p, modeQuery)
}

// PathEscape escapes the string so it can be safely placed inside a URL query.
//
// Uses built-in buffer.
func (vec *Vector) PathEscape(p []byte) []byte {
	return vecEscape(vec, p, modePath)
}

// PathUnescape does the inverse transformation of PathEscape.
func (vec *Vector) PathUnescape(p []byte) []byte {
	return vecUnescape(vec, p, modePath)
}
