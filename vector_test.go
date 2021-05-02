package urlvector

import "testing"

var (
	url0 = []byte("https://john_ruth:hangman17@99.99.99.99:3306/foo/bar?that\\'s#all, folks")

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
	t.Log(vec.PathStr())
	t.Log(vec.QueryStr())
	if err != nil {
		// t.Error(err)
	}
}
