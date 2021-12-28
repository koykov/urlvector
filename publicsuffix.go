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

func (m *PublicSuffixDB) add(ps []byte) {
	var e pse
	lo := uint32(len(m.buf))
	hi := uint32(len(ps)) + lo
	e.encode(lo, hi)
	m.idx = append(m.idx, e)
	m.buf = append(m.buf, ps...)
}

func (m *PublicSuffixDB) AddStr(ps string) {
	m.Add(fastconv.S2B(ps))
}

func (m *PublicSuffixDB) Reset() {
	m.idx = m.idx[:0]
	m.buf = m.buf[:0]
}

func psMustSkip(line []byte) bool {
	if len(line) == 0 || line[0] == '/' || line[0] == '!' {
		return true
	}
	return false
}
