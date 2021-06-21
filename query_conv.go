package urlvector

const (
	hex = "\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x00\x01\x02\x03\x04\x05\x06\a\b\t\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10"
	// Hex digits in upper case.
	hexUp = "0123456789ABCDEF"
)

// QueryEscape escapes the string so it can be safely placed inside a URL query.
func QueryUnescape(dst, p []byte) []byte {
	dst = append(dst, p...)
	return queryUnescape(dst)
}

// QueryUnescape does the inverse transformation of QueryEscape.
func QueryEscape(dst, p []byte) []byte {
	l := len(p)
	if l == 0 {
		return dst
	}
	_ = p[l-1]
	for i := 0; i < l; i++ {
		if p[i] >= 'a' && p[i] <= 'z' || p[i] >= 'A' && p[i] <= 'Z' ||
			p[i] >= '0' && p[i] <= '9' || p[i] == '-' || p[i] == '.' || p[i] == '_' {
			dst = append(dst, p[i])
		} else if p[i] == ' ' {
			dst = append(dst, '+')
		} else {
			dst = append(dst, '%')
			dst = append(dst, hexUp[p[i]>>4])
			dst = append(dst, hexUp[p[i]&15])
		}
	}
	return dst
}

// PathEscape escapes the string so it can be safely placed inside a URL query.
func PathUnescape(dst, p []byte) []byte {
	dst = append(dst, p...)
	return pathUnescape(dst)
}

// PathUnescape does the inverse transformation of PathEscape.
func PathEscape(dst, p []byte) []byte {
	l := len(p)
	if l == 0 {
		return dst
	}
	_ = p[l-1]
	for i := 0; i < l; i++ {
		if (p[i] >= 'a' && p[i] <= 'z' || p[i] >= 'A' && p[i] <= 'Z' ||
			p[i] >= '0' && p[i] <= '9' || p[i] == '-' || p[i] == '.' || p[i] == '_') && p[i] != '?' {
			dst = append(dst, p[i])
		} else if p[i] == ' ' {
			dst = append(dst, '+')
		} else {
			dst = append(dst, '%')
			dst = append(dst, hexUp[p[i]>>4])
			dst = append(dst, hexUp[p[i]&15])
		}
	}
	return dst
}

// Unescape byte array using itself as a destination.
func queryUnescape(p []byte) []byte {
	l := len(p)
	n := len(p)
	if l < 3 {
		return p
	}
	_ = p[l-1]
	for i := 0; i < n; i++ {
		switch p[i] {
		case '%':
			if i+2 < l {
				x2 := hex[p[i+2]]
				x1 := hex[p[i+1]]
				if x1 != 16 || x2 != 16 {
					p[i] = x1<<4 | x2
					copy(p[i+1:], p[i+3:])
					n -= 2
				}
			}
		case '+':
			p[i] = ' '
		}
	}
	return p[:n]
}

// Query escape p using vec buffer.
func queryEscape(vec *Vector, p []byte) []byte {
	l := len(p)
	if l == 0 {
		return p
	}
	o := vec.BufLen()
	buf := QueryEscape(vec.Buf(), p)
	vec.BufUpdateWith(buf)
	return vec.Buf()[o:]
}

// Unescape byte array using itself as a destination.
func pathUnescape(p []byte) []byte {
	return queryUnescape(p)
}

// Query escape p using vec buffer.
func pathEscape(vec *Vector, p []byte) []byte {
	l := len(p)
	if l == 0 {
		return p
	}
	o := vec.BufLen()
	buf := PathEscape(vec.Buf(), p)
	vec.BufUpdateWith(buf)
	return vec.Buf()[o:]
}
