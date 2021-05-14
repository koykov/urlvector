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
		{
			"http://yolo.com\\\\what-is-up.com",
			testTarget{path: "\\\\what-is-up.com"},
		},
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
	url0 = []byte("https://john_ruth:hangman17@99.99.99.99:3306/foo/bar?that\\'s#all, folks")

	query0       = []byte("http://localhost:8011/get_data?v=default&blockID=319385&page=https%3A%2F%2Fultra-software-base.ru%2Fsystem%2Fgoogle-chrome.html%3Fyclid%3D212247430717539672&domain=ultra-software-base.ru&uid=4f5d0edc-3a3e-48d0-9872-0b48a7998ac6&clientNotice=true&imgX=360&imgY=240&limit=1&subage_dt=2021-01-29&format=json&cur=RUB&ua=Mozilla%2F5.0+%28Windows+NT+6.1%3B+Win64%3B+x64%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F89.0.4389.105+YaBrowser%2F21.3.3.230+Yowser%2F2.5+Safari%2F537.36&ip=5.5.5.5&subage=102&language=ru")
	query0target = map[string]string{
		"v":            "default",
		"blockID":      "319385",
		"page":         "https://ultra-software-base.ru/system/google-chrome.html?yclid=212247430717539672",
		"domain":       "ultra-software-base.ru",
		"uid":          "4f5d0edc-3a3e-48d0-9872-0b48a7998ac6",
		"clientNotice": "true",
		"imgX":         "360",
		"imgY":         "240",
		"limit":        "1",
		"subage_dt":    "2021-01-29",
		"format":       "json",
		"cur":          "RUB",
		"ua":           "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 YaBrowser/21.3.3.230 Yowser/2.5 Safari/537.36",
		"ip":           "5.5.5.5",
		"subage":       "102",
		"language":     "ru",
	}

	query1 = []byte("http://x.com/1?x&y=1&z")
	query2 = []byte("http://x.com/x/y/z?arr[]=1&arr[]=2&arr[]=3&b=x")

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

		if len(tst.target.scheme) > 0 && vec.SchemeString() != tst.target.scheme {
			printErr(t, &tst, "scheme mismatch", vec.SchemeString(), "vs", tst.target.scheme)
		}
		if tst.target.slashes && vec.Slashes() != tst.target.slashes {
			printErr(t, &tst, "slashes mismatch", vec.Slashes(), "vs", tst.target.slashes)
		}
		if len(tst.target.auth) > 0 && vec.AuthString() != tst.target.auth {
			printErr(t, &tst, "auth mismatch", vec.AuthString(), "vs", tst.target.auth)
		}
		if len(tst.target.username) > 0 && vec.UsernameString() != tst.target.username {
			printErr(t, &tst, "username mismatch", vec.UsernameString(), "vs", tst.target.username)
		}
		if len(tst.target.password) > 0 && vec.PasswordString() != tst.target.password {
			printErr(t, &tst, "password mismatch", vec.PasswordString(), "vs", tst.target.password)
		}
		if len(tst.target.host) > 0 && vec.HostString() != tst.target.host {
			printErr(t, &tst, "host mismatch", vec.HostString(), "vs", tst.target.host)
		}
		if len(tst.target.hostname) > 0 && vec.HostnameString() != tst.target.hostname {
			printErr(t, &tst, "hostname mismatch", vec.HostnameString(), "vs", tst.target.hostname)
		}
		if tst.target.port > 0 && vec.Port() != tst.target.port {
			printErr(t, &tst, "port mismatch", vec.Port(), "vs", tst.target.port)
		}
		if len(tst.target.path) > 0 && vec.PathString() != tst.target.path {
			printErr(t, &tst, "path mismatch", vec.PathString(), "vs", tst.target.path)
		}
		if len(tst.target.query) > 0 && vec.QueryString() != tst.target.query {
			printErr(t, &tst, "query mismatch", vec.QueryString(), "vs", tst.target.query)
		}
		if len(tst.target.hash) > 0 && vec.HashString() != tst.target.hash {
			printErr(t, &tst, "hash mismatch", vec.HashString(), "vs", tst.target.hash)
		}
	}
}

func TestVector_ParseQuery(t *testing.T) {
	// vec.Reset()
	// _ = vec.Parse(query0)
	// query := vec.Query()
	// query.Each(func(_ int, node *vector.Node) {
	// 	k := node.KeyString()
	// 	if query0target[k] != node.String() {
	// 		t.Error("query0 mismatch query param", k, "need", query0target[k], "got", node.String())
	// 	}
	// })
	//
	// vec.Reset()
	// _ = vec.Parse(query1)
	// query = vec.Query()
	// if !query.Exists("x") || query.Get("x").String() != "" {
	// 	t.Error("query1 mismatch query param x")
	// }
	// if !query.Exists("z") || query.Get("z").String() != "" {
	// 	t.Error("query1 mismatch query param z")
	// }

	vec.Reset()
	_ = vec.Parse(query2)
	query := vec.Query()
	query.Each(func(_ int, node *vector.Node) {
		t.Log(node.KeyString(), "=", node.String())
	})
}

func BenchmarkVector_Parse(b *testing.B) {
	tst := testTargets{
		url: string(url0),
		target: testTarget{
			scheme:   "https",
			auth:     "john_ruth:hangman17",
			username: "john_ruth",
			password: "hangman17",
			host:     "99.99.99.99:3306",
			hostname: "99.99.99.99",
			path:     "/foo/bar",
			query:    "?that\\'s",
			hash:     "#all, folks",
			slashes:  false,
			port:     3306,
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		vec.Reset()
		if err := vec.ParseStr(tst.url); err != nil {
			b.Fatal(err)
		}

		if vec.SchemeString() != tst.target.scheme {
			printErr(b, &tst, "scheme mismatch", vec.SchemeString(), "vs", tst.target.scheme)
		}
		if tst.target.slashes && vec.Slashes() != tst.target.slashes {
			printErr(b, &tst, "slashes mismatch", vec.Slashes(), "vs", tst.target.slashes)
		}
		if len(tst.target.auth) > 0 && vec.AuthString() != tst.target.auth {
			printErr(b, &tst, "auth mismatch", vec.AuthString(), "vs", tst.target.auth)
		}
		if len(tst.target.username) > 0 && vec.UsernameString() != tst.target.username {
			printErr(b, &tst, "username mismatch", vec.UsernameString(), "vs", tst.target.username)
		}
		if len(tst.target.password) > 0 && vec.PasswordString() != tst.target.password {
			printErr(b, &tst, "password mismatch", vec.PasswordString(), "vs", tst.target.password)
		}
		if len(tst.target.host) > 0 && vec.HostString() != tst.target.host {
			printErr(b, &tst, "host mismatch", vec.HostString(), "vs", tst.target.host)
		}
		if len(tst.target.hostname) > 0 && vec.HostnameString() != tst.target.hostname {
			printErr(b, &tst, "hostname mismatch", vec.HostnameString(), "vs", tst.target.hostname)
		}
		if tst.target.port > 0 && vec.Port() != tst.target.port {
			printErr(b, &tst, "port mismatch", vec.Port(), "vs", tst.target.port)
		}
		if len(tst.target.path) > 0 && vec.PathString() != tst.target.path {
			printErr(b, &tst, "path mismatch", vec.PathString(), "vs", tst.target.path)
		}
		if len(tst.target.query) > 0 && vec.QueryString() != tst.target.query {
			printErr(b, &tst, "query mismatch", vec.QueryString(), "vs", tst.target.query)
		}
		if len(tst.target.hash) > 0 && vec.HashString() != tst.target.hash {
			printErr(b, &tst, "hash mismatch", vec.HashString(), "vs", tst.target.hash)
		}
	}
}

func BenchmarkVector_ParseQuery(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		vec.Reset()
		_ = vec.ParseCopy(query0)
		query := vec.Query()
		query.Each(func(_ int, node *vector.Node) {
			k := node.KeyString()
			if query0target[k] != node.String() {
				b.Error("query0 mismatch query param", k, "need", query0target[k], "got", node.String())
			}
		})
	}
}

func printErr(t testing.TB, tst *testTargets, args ...interface{}) {
	t.Error("\nsrc: "+tst.url+"\n", args)
}
