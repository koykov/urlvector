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
		{
			" javascript://foo",
			testTarget{scheme: "javascript", host: "foo"},
		},
		{
			"http://google.com/?foo=bar",
			testTarget{scheme: "http", host: "google.com", query: "?foo=bar"},
		},
		{
			"http://google.com/?lolcakes",
			testTarget{scheme: "http", host: "google.com", query: "?lolcakes"},
		},
		{
			"blob:https://gist.github.com/3f272586-6dac-4e29-92d0-f674f2dde618",
			testTarget{scheme: "https", host: "gist.github.com", path: "/3f272586-6dac-4e29-92d0-f674f2dde618"},
		},
		{
			"https://www.mozilla.org/en-US/firefox/34.0/whatsnew/?oldversion=33.1",
			testTarget{scheme: "https", host: "www.mozilla.org", path: "/en-US/firefox/34.0/whatsnew/", query: "?oldversion=33.1"},
		},
		{
			"http://example.com:80",
			testTarget{host: "example.com:80"},
		},
		{
			"http://example.com:80/",
			testTarget{host: "example.com:80", path: "/"},
		},
		{
			"http://x.com/path?that\\'s#all, folks",
			testTarget{host: "x.com", path: "/path", query: "?that\\'s", hash: "#all, folks"},
		},
		{
			"http://google.com:80\\\\@yahoo.com/#what\\\\is going on",
			testTarget{username: "google.com", password: "80\\\\", host: "yahoo.com", hash: "#what\\\\is going on"},
		},
		// {
		// 	"http://yolo.com\\\\what-is-up.com",
		// 	testTarget{path: "/what-is-up.com"},
		// },
		{
			"HTTP://example.com",
			testTarget{scheme: "HTTP"}, // fixme must be "http"!
		},
		{
			"google.com/foo",
			testTarget{scheme: "", host: "google.com", path: "/foo"},
		},
		{
			"http://[1080:0:0:0:8:800:200C:417A]:61616/foo/bar?q=z",
			testTarget{host: "[1080:0:0:0:8:800:200C:417A]:61616", hostname: "[1080:0:0:0:8:800:200C:417A]", port: 61616},
		},
		{
			"http://user:password@[3ffe:2a00:100:7031::1]:8080/",
			testTarget{username: "user", password: "password", hostname: "[3ffe:2a00:100:7031::1]", port: 8080},
		},
		{
			"http://222.148.142.13:61616/foo/bar?q=z",
			testTarget{hostname: "222.148.142.13", port: 61616, path: "/foo/bar", query: "?q=z"},
		},
		{
			"HTTP://USER:PASS@EXAMPLE.COM",
			testTarget{
				scheme:   "HTTP", // fixme must be "http"!
				username: "USER", password: "PASS", host: "EXAMPLE.COM",
			},
		},
		{
			"http://mt0.google.com/vt/lyrs=m@114&hl=en&src=api&x=2&y=2&z=3&s=",
			testTarget{path: "/vt/lyrs=m@114&hl=en&src=api&x=2&y=2&z=3&s="},
		},
		{
			"http://user@www.example.com/",
			testTarget{username: "user", host: "www.example.com"},
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
