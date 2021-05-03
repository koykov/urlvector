package urlvector

import (
	"testing"

	"github.com/koykov/vector"
)

type testTarget struct {
	scheme, auth, username, password, host, hostname, path, query, hash string

	slashes bool
	port    int
	err     error
	errOff  int
}

type testTargets struct {
	url    string
	target testTarget
}

var (
	cases = []testTargets{
		{
			"",
			testTarget{err: vector.ErrEmptySrc, errOff: 0},
		},
		{
			"/foo",
			testTarget{path: "/foo"},
		},
		{
			"http://example.com",
			testTarget{scheme: "http", host: "example.com", hostname: "example.com"},
		},
		{
			"//foo/bar",
			testTarget{slashes: true, host: "foo", hostname: "foo", path: "/bar"},
		},
		{
			"foo\\nbar\\rbaz\\u2028qux\\u2029",
			testTarget{err: vector.ErrUnparsedTail, errOff: 0},
		},
	}
	// url0 = []byte("https://john_ruth:hangman17@99.99.99.99:3306/foo/bar?that\\'s#all, folks")

	vec = NewVector()
)

func TestVector_Parse(t *testing.T) {
	for i, tst := range cases {
		vec.Reset()
		err := vec.ParseStr(tst.url)
		if err != nil {
			if err != tst.target.err {
				t.Error(i, err)
			}
			continue
		}

		if len(tst.target.scheme) > 0 && vec.SchemeStr() != tst.target.scheme {
			t.Error("scheme mismatch", vec.SchemeStr(), "vs", tst.target.scheme)
		}
		if tst.target.slashes && vec.Slashes() != tst.target.slashes {
			t.Error("slashes mismatch", vec.Slashes(), "vs", tst.target.slashes)
		}
		if len(tst.target.auth) > 0 && vec.AuthStr() != tst.target.auth {
			t.Error("auth mismatch", vec.AuthStr(), "vs", tst.target.auth)
		}
		if len(tst.target.username) > 0 && vec.UsernameStr() != tst.target.username {
			t.Error("username mismatch", vec.UsernameStr(), "vs", tst.target.username)
		}
		if len(tst.target.password) > 0 && vec.PasswordStr() != tst.target.password {
			t.Error("password mismatch", vec.PasswordStr(), "vs", tst.target.password)
		}
		if len(tst.target.host) > 0 && vec.HostStr() != tst.target.host {
			t.Error("host mismatch", vec.HostStr(), "vs", tst.target.host)
		}
		if len(tst.target.hostname) > 0 && vec.HostnameStr() != tst.target.hostname {
			t.Error("hostname mismatch", vec.HostnameStr(), "vs", tst.target.hostname)
		}
		if tst.target.port > 0 && vec.Port() != tst.target.port {
			t.Error("port mismatch", vec.Port(), "vs", tst.target.port)
		}
		if len(tst.target.path) > 0 && vec.PathStr() != tst.target.path {
			t.Error("path mismatch", vec.PathStr(), "vs", tst.target.path)
		}
		if len(tst.target.query) > 0 && vec.QueryStr() != tst.target.query {
			t.Error("query mismatch", vec.QueryStr(), "vs", tst.target.query)
		}
		if len(tst.target.hash) > 0 && vec.HashStr() != tst.target.hash {
			t.Error("hash mismatch", vec.HashStr(), "vs", tst.target.hash)
		}
	}
}
