package urlvector

type mode int

const (
	modePath mode = iota
	modeQuery
	modeHash

	hex = "\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x00\x01\x02\x03\x04\x05\x06\a\b\t\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\n\v\f\r\x0e\x0f\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10"
	// Hex digits in upper case.
	hexUp = "0123456789ABCDEF"
)

var (
	_, _, _, _ = QueryEscape, QueryUnescape, PathEscape, PathUnescape
)

// QueryEscape escapes the string so it can be safely placed inside a URL query.
func QueryEscape(dst, p []byte) []byte {
	return bufEscape(dst, p, modeQuery)
}

// QueryUnescape does the inverse transformation of QueryEscape.
func QueryUnescape(dst, p []byte) []byte {
	return bufUnescape(dst, p, modeQuery)
}

// PathEscape escapes the string so it can be safely placed inside a URL query.
func PathEscape(dst, p []byte) []byte {
	return bufEscape(dst, p, modePath)
}

// PathUnescape does the inverse transformation of PathEscape.
func PathUnescape(dst, p []byte) []byte {
	return bufUnescape(dst, p, modePath)
}

// Unescape byte array using itself as a destination.
func unescape(p []byte) []byte {
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

func bufEscape(dst, p []byte, mode mode) []byte {
	l := len(p)
	if l == 0 {
		return dst
	}
	_ = p[l-1]
	for i := 0; i < l; i++ {
		allow := p[i] >= 'a' && p[i] <= 'z' || p[i] >= 'A' && p[i] <= 'Z' || p[i] >= '0' && p[i] <= '9' ||
			p[i] == '-' || p[i] == '.' || p[i] == '_'
		switch mode {
		case modePath:
			allow = allow || p[i] == '/'
		case modeQuery:
			allow = allow && p[i] != '?'
		case modeHash:
			allow = allow || p[i] == '#' || p[i] == '?' || p[i] == '=' || p[i] == '!' || p[i] == '(' || p[i] == ')' || p[i] == '*'
		}
		if allow {
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

func bufUnescape(dst, p []byte, mode mode) []byte {
	_, _ = dst, mode
	dst = unescape(p)
	return dst
}

func vecEscape(vec *Vector, p []byte, mode mode) []byte {
	l := len(p)
	if l == 0 {
		return p
	}
	o := vec.BufLen()
	buf := bufEscape(vec.Buf(), p, mode)
	vec.BufUpdateWith(buf)
	return vec.Buf()[o:]
}

func vecUnescape(vec *Vector, p []byte, mode mode) []byte {
	_, _ = vec, mode
	return unescape(p)
}
