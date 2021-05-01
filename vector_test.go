package urlvector

import "testing"

var (
	url0 = []byte("https://john_ruth:hangman17@github.com:3306/foo/bar")

	vec = NewVector()
)

func TestVector_Parse(t *testing.T) {
	vec.Reset()
	err := vec.Parse(url0)
	t.Log(vec.SchemeStr())
	t.Log(vec.Slashes())
	t.Log(vec.AuthStr())
	t.Log(vec.UsernameStr())
	t.Log(vec.PasswordStr())
	t.Log(vec.HostStr())
	t.Log(vec.HostnameStr())
	t.Log(vec.Port())
	if err != nil {
		// t.Error(err)
	}
}
