package urlvector

import (
	"bufio"
	"bytes"
	"net/http"
	"os"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

const (
	fullURL = "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"
)

type PublicSuffixDB struct {
	idx    []pse
	idxl   int
	buf, r []byte
}

var (
	bPrefixAllPS = []byte("*.")
)

func (m *PublicSuffixDB) Load(dbFile string) (err error) {
	var file *os.File
	if file, err = os.OpenFile(dbFile, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := bytealg.TrimLeft(scan.Bytes(), bSpace)
		if psMustSkip(line) {
			continue
		}
		m.Add(line)
	}
	if len(m.r) > 0 {
		m.add(m.r)
	}
	err = scan.Err()

	return
}

func (m *PublicSuffixDB) Fetch(dbURL string) (err error) {
	var resp *http.Response
	if resp, err = http.Get(dbURL); err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	scan := bufio.NewScanner(resp.Body)
	for scan.Scan() {
		line := bytealg.TrimLeft(scan.Bytes(), bSpace)
		if psMustSkip(line) {
			continue
		}
		m.Add(line)
	}
	if len(m.r) > 0 {
		m.add(m.r)
	}
	err = scan.Err()

	return
}

func (m *PublicSuffixDB) FetchFull() error {
	return m.Fetch(fullURL)
}

func (m *PublicSuffixDB) Add(ps []byte) {
	if bytes.Equal(ps[:2], bPrefixAllPS) {
		ps = ps[2:]
	}
	if len(ps) == 0 {
		return
	}
	if !bytealg.HasByteLR(ps, '.') {
		if len(m.r) > 0 {
			m.add(m.r)
		}
		m.r = append(m.r[:0], ps...)
		return
	}
	m.add(ps)
	return
}

func (m *PublicSuffixDB) AddStr(ps string) {
	m.Add(fastconv.S2B(ps))
}

func (m *PublicSuffixDB) add(ps []byte) {
	var e pse
	lo := uint32(len(m.buf))
	hi := uint32(len(ps)) + lo
	e.encode(lo, hi)
	m.idxl++
	m.idx = append(m.idx, e)
	m.buf = append(m.buf, ps...)
}

func (m PublicSuffixDB) Get(hostname []byte) (ps []byte) {
	ps, _ = m.GetWP(hostname)
	return
}

func (m PublicSuffixDB) GetWP(hostname []byte) ([]byte, int) {
	hl := len(hostname)
	if hl < 2 || m.idxl == 0 {
		return nil, -1
	}
	var (
		ps  []byte
		pos int
		ok  bool
	)
	b := m.idx
	_ = b[len(b)-1]
	for len(b) >= 8 {
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[0]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[1]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[2]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[3]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[4]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[5]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[6]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[7]); ok {
			return ps, pos
		}
		b = b[8:]
	}
	for len(b) >= 4 {
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[0]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[1]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[2]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[3]); ok {
			return ps, pos
		}
		b = b[4:]
	}
	for len(b) >= 2 {
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[0]); ok {
			return ps, pos
		}
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[1]); ok {
			return ps, pos
		}
		b = b[2:]
	}
	if len(b) == 1 {
		if ps, pos, ok = m.hostHasPSE(hostname, hl, b[0]); ok {
			return ps, pos
		}
	}
	return nil, -1
}

func (m PublicSuffixDB) GetStr(hostname string) (ps string) {
	ps, _ = m.GetStrWP(hostname)
	return
}

func (m PublicSuffixDB) GetStrWP(hostname string) (ps string, pos int) {
	x, p := m.GetWP(fastconv.S2B(hostname))
	if p == -1 {
		return
	}
	ps, pos = fastconv.B2S(x), p
	return
}

func (m *PublicSuffixDB) Reset() {
	m.idxl = 0
	m.idx = m.idx[:0]
	m.buf = m.buf[:0]
}

func (m PublicSuffixDB) hostHasPSE(hostname []byte, hl int, e pse) (ps []byte, pos int, ok bool) {
	pos = -1
	lo, hi := e.decode()
	if hi-lo >= uint32(hl) {
		return
	}
	a := m.buf[lo:hi]
	b := hostname[uint32(hl)-(hi-lo)-1:]
	if b[0] == '.' && bytes.Equal(a, b[1:]) {
		ps = b[1:]
		pos = hl - int(hi-lo)
		ok = true
	}
	return
}

func psMustSkip(line []byte) bool {
	if len(line) == 0 || line[0] == '/' || line[0] == '!' {
		return true
	}
	return false
}
