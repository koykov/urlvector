package urlvector

import "testing"

var (
	url0 = []byte("https://github.com/foo/bar")

	vec = NewVector()
)

func TestVector_Parse(t *testing.T) {
	vec.Reset()
	err := vec.Parse(url0)
	if err != nil {
		t.Error(err)
	}
}
