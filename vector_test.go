package urlvector

import "testing"

var (
	url0 = []byte("https://john_ruth:hangman17@github.com/foo/bar")

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
	if err != nil {
		// t.Error(err)
	}
}
