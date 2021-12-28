package urlvector

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

const (
	fullURL = "https://raw.githubusercontent.com/publicsuffix/list/master/public_suffix_list.dat"
)

type PublicSuffixDB struct {
	idx []pse
	buf []byte
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

var c int

func (m *PublicSuffixDB) Add(ps []byte) (lo, hi uint32) {
	if bytes.Equal(ps[:2], bPrefixAllPS) {
		ps = ps[2:]
	}
	if !bytealg.HasByteLR(ps, '.') && len(ps) < 4 {
		c++
		fmt.Println(string(ps), c)
	}
	if hi = uint32(len(ps)); hi == 0 {
		return
	}
	var e pse
	lo = uint32(len(m.buf))
	hi += lo
	e.encode(lo, hi)
	m.idx = append(m.idx, e)
	m.buf = append(m.buf, ps...)
	return
}

func (m *PublicSuffixDB) AddStr(ps string) (uint32, uint32) {
	return m.Add(fastconv.S2B(ps))
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
