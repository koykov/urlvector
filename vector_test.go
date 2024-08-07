package urlvector

import (
	"bytes"
	"fmt"
	"strings"
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
		{
			"https://www.msn.com/ru-ru/lifestyle/travel/на-фото-памятные-достопримечательности-из-разных-уголков-планеты/ss-AAGnFe0#image=4?ocid=ems.msn.dl.090919.TowerOfPisaItaly",
			testTarget{hash: "#image=4?ocid=ems.msn.dl.090919.TowerOfPisaItaly"},
		},
		{
			"https://lamimyde.pro/bB3.VC0DPE2_lGjHPIXJB-zLJMmN9O0_PQUR5SETR-TVJWNXRYE_Va5bTcVdN-RfegEh5iq_SkTlBmNna-mpMqxrTsn_puVvMwUxp-FzbAFBVCK_RE0FpG5HY-jJJKGLaM1_lOtPRQnRV-aTQU1VFWr_SYkZRaKba-ldpeXfUgT_BiOjVk0l0-0nWoWppqC_as0t1uUvZ-GxxyPzRAG_dC3DTEkFR-JHMIVJpKt_WMmN1ONPR-ER5SpTTUX_pWRXeYEZ9-XbUcTdAem_cgnhJipjZ-Dl0mznYo2_Mq3rZsWtQ-zvMwjxUyx_NAGBYC2DZ-DFEG3HMIj_UK3LOMDNg-5PMQDRlSj_NUDVFWiXY-zZUa2bNcy_Zeyfcg3hJ-jjPkXlNmj_complqwrd-CtZu2vdwD_0yxzNAjBI-0DMEjFcG1_NIzJUK1LL-TNVOkPZQG_YS1TYUTVV-hXZYWZUa5_NczdgexfZ-WhRimjMkz_FmjnNoWpJ-krMsTtQu5_ZwWxIy5zY-TBlCjD?bd_vs=1.0.2&bd_t0=1624275754931&bd_a0=a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15,a16,a17,a18,a19,a20,a21,a22,a23,a24,a25,a26,a27,a28,a29,a30,a31,a32,a33,a34,a35,a36,a37,a38,a39,a40,a41,a42,a43,a44,a45,a46,a47,a48,a49,a50,a51,a52,a53,a54,a55,a56,a57,sm0&bd_a1=59&bd_n0=Mozilla/5.0 (Linux; Android 9; SAMSUNG SM-G950F/G950FXXUCDUD1) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.0 Chrome/87.0.4280.141 Mobile Safari/537.36&bd_n1=[\\\"it-IT\\\",\\\"it\\\",\\\"en-US\\\",\\\"en\\\"]&bd_n2=8&bd_n3=20030107&bd_n4=Google Inc.&bd_o1=false&bd_o4=&bd_o5=false&sseq=3&dseq=3&ce=lnk",
			testTarget{
				host:  "lamimyde.pro",
				path:  "/bB3.VC0DPE2_lGjHPIXJB-zLJMmN9O0_PQUR5SETR-TVJWNXRYE_Va5bTcVdN-RfegEh5iq_SkTlBmNna-mpMqxrTsn_puVvMwUxp-FzbAFBVCK_RE0FpG5HY-jJJKGLaM1_lOtPRQnRV-aTQU1VFWr_SYkZRaKba-ldpeXfUgT_BiOjVk0l0-0nWoWppqC_as0t1uUvZ-GxxyPzRAG_dC3DTEkFR-JHMIVJpKt_WMmN1ONPR-ER5SpTTUX_pWRXeYEZ9-XbUcTdAem_cgnhJipjZ-Dl0mznYo2_Mq3rZsWtQ-zvMwjxUyx_NAGBYC2DZ-DFEG3HMIj_UK3LOMDNg-5PMQDRlSj_NUDVFWiXY-zZUa2bNcy_Zeyfcg3hJ-jjPkXlNmj_complqwrd-CtZu2vdwD_0yxzNAjBI-0DMEjFcG1_NIzJUK1LL-TNVOkPZQG_YS1TYUTVV-hXZYWZUa5_NczdgexfZ-WhRimjMkz_FmjnNoWpJ-krMsTtQu5_ZwWxIy5zY-TBlCjD",
				query: "?bd_vs=1.0.2&bd_t0=1624275754931&bd_a0=a0,a1,a2,a3,a4,a5,a6,a7,a8,a9,a10,a11,a12,a13,a14,a15,a16,a17,a18,a19,a20,a21,a22,a23,a24,a25,a26,a27,a28,a29,a30,a31,a32,a33,a34,a35,a36,a37,a38,a39,a40,a41,a42,a43,a44,a45,a46,a47,a48,a49,a50,a51,a52,a53,a54,a55,a56,a57,sm0&bd_a1=59&bd_n0=Mozilla/5.0 (Linux; Android 9; SAMSUNG SM-G950F/G950FXXUCDUD1) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.0 Chrome/87.0.4280.141 Mobile Safari/537.36&bd_n1=[\\\"it-IT\\\",\\\"it\\\",\\\"en-US\\\",\\\"en\\\"]&bd_n2=8&bd_n3=20030107&bd_n4=Google Inc.&bd_o1=false&bd_o4=&bd_o5=false&sseq=3&dseq=3&ce=lnk",
			},
		},
		{
			"http://sia.ru?section=484&action=show_news&id=414088",
			testTarget{
				host:  "sia.ru",
				query: "?section=484&action=show_news&id=414088",
			},
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
	query0sorted = []byte("?blockID=319385&clientNotice=true&cur=RUB&domain=ultra-software-base.ru&format=json&imgX=360&imgY=240&ip=5.5.5.5&language=ru&limit=1&page=https%3A%2F%2Fultra-software-base.ru%2Fsystem%2Fgoogle-chrome.html%3Fyclid%3D212247430717539672&subage=102&subage_dt=2021-01-29&ua=Mozilla%2F5.0+%28Windows+NT+6.1%3B+Win64%3B+x64%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Chrome%2F89.0.4389.105+YaBrowser%2F21.3.3.230+Yowser%2F2.5+Safari%2F537.36&uid=4f5d0edc-3a3e-48d0-9872-0b48a7998ac6&v=default")
	// 16, 1, 11, 4, 15, 2, 6, 7, 10, 13, 5, 3, 14, 8, 12, 9

	query1 = []byte("http://x.com/1?x&y=1&z")
	query2 = []byte("http://x.com/x/y/z?arr[]=1&arr[]=2&arr[]=3&b=x&arr1[]=a&arr1[]=b&arr1[]=c")

	query3     = []byte("http://x.com/a/b/c?x=1&y=qwerty&z=foo")
	query3repl = []byte("?foo=x&bar=y&a[]=1&a[]=2&a[]=c&b[]=qwe&b[]=rty&z")
	query3new  = []byte("https://foo:bar@google.com:8080/search?q=keys#results")

	queryRemove       = `https://example.ru/services/context/?q=foobar&utm_source=google_adwords_Poisk&utm_medium=cpc&utm_geo=Москва&utm_campaign=Kontekstnaya_reklama_Moskva_Poisk&arg=qwerty&utm_content=Nastroyka_kontekstnoy_reklami&utm_term=настройка контекстной рекламы&gclid=Cj0KCQjwn_LrBRD4ARIsAFEQFKvL-kzn-g2cbcezy6GTobnpRIGtWQgnzZjR2c06qukAWMaLKm0fG40aAtTlEALw_wc`
	queryRemoveExpect = `https://example.ru/services/context/?arg=qwerty&gclid=Cj0KCQjwn_LrBRD4ARIsAFEQFKvL-kzn-g2cbcezy6GTobnpRIGtWQgnzZjR2c06qukAWMaLKm0fG40aAtTlEALw_wc&q=foobar`

	vec = NewVector()
)

func printErr(t testing.TB, tst *testTargets, args ...any) {
	t.Error("\nsrc: "+tst.url+"\n", args)
}

func TestVector(t *testing.T) {
	for i, tst := range cases {
		t.Run(fmt.Sprintf("parse/%d", i), func(t *testing.T) {
			vec.Reset()
			err := vec.ParseStr(tst.url)
			if err != nil {
				if err != tst.target.err {
					t.Error(i, err)
				}
				return
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
		})
	}

	t.Run("parse query", func(t *testing.T) {
		var query *vector.Node

		vec.Reset()
		_ = vec.ParseCopy(query0)
		query = vec.Query()
		query.Each(func(_ int, node *vector.Node) {
			k := node.KeyString()
			if query0target[k] != node.String() {
				t.Error("query0 mismatch query param", k, "need", query0target[k], "got", node.String())
			}
		})

		vec.Reset()
		_ = vec.Parse(query1)
		query = vec.Query()
		if !query.Exists("x") || query.Get("x").String() != "" {
			t.Error("query1 mismatch query param x")
		}
		if !query.Exists("z") || query.Get("z").String() != "" {
			t.Error("query1 mismatch query param z")
		}

		vec.Reset()
		_ = vec.Parse(query2)
		query = vec.Query()
		query.Each(func(_ int, node *vector.Node) {
			switch {
			case node.KeyString() == "b":
				if node.String() != "x" {
					t.Error("query 2 mismatch query param", node.KeyString(), "need", "x", "got", node.String())
				}
			case node.KeyString() == "arr[]":
				if node.Limit() != 3 {
					t.Error("query 2 unexpected length of param arr[]", "need", 3, "got", node.Limit())
				}
				if node.At(1).String() != "2" {
					t.Error("query 2 mismatch query param arr[1]", "need", "2", "got", node.At(1).String())
				}
			case node.KeyString() == "arr1[]":
				if node.Limit() != 3 {
					t.Error("query 2 unexpected length of param arr1[]", "need", 3, "got", node.Limit())
				}
				if node.At(0).String() != "a" {
					t.Error("query 2 mismatch query param arr1[0]", "need", "a", "got", node.At(0).String())
				}
			}
		})
	})

	t.Run("remove if", func(t *testing.T) {
		vec.Reset()
		_ = vec.ParseStr(queryRemove)
		vec.Query().RemoveIf(func(_ int, node *vector.Node) bool {
			return strings.HasPrefix(node.KeyString(), "utm")
		})
		if vec.QuerySort().String() != queryRemoveExpect {
			t.FailNow()
		}
	})

	t.Run("set", func(t *testing.T) {
		vec.Reset()
		_ = vec.Parse(url0)
		vec.SetHostnameString("x.com")
		if h := vec.HostnameString(); h != "x.com" {
			t.Error("query 2 mismatch query param arr1[0]", "need", "x.com", "got", h)
		}
		vec.SetPort(9999)
		if p := vec.Port(); p != 9999 {
			t.Error("query 2 mismatch query param port", "need", 9999, "got", p)
		}
	})

	t.Run("forget query params", func(t *testing.T) {
		vec.Reset()
		_ = vec.Parse(query3)
		if y := vec.Query().GetString("y"); y != "qwerty" {
			t.Error("query 3 mismatch query param y", "need", "qwerty", "got", y)
		}
		vec.SetQueryBytes(query3repl)
		vec.Query().Each(func(_ int, node *vector.Node) {
			switch {
			case node.KeyString() == "foo":
				if node.String() != "x" {
					t.Error("query 3 (forget) mismatch query param", node.KeyString(), "need", "x", "got", node.String())
				}
			case node.KeyString() == "bar":
				if node.String() != "y" {
					t.Error("query 3 (forget) mismatch query param", node.KeyString(), "need", "y", "got", node.String())
				}
			case node.KeyString() == "a[]":
				if node.Limit() != 3 {
					t.Error("query 2 (forget) unexpected length of param a[]", "need", 3, "got", node.Limit())
				}
				if node.At(0).String() != "1" {
					t.Error("query 2 (forget) mismatch query param a[0]", "need", "1", "got", node.At(0).String())
				}
			}
		})
	})

	t.Run("serialize", func(t *testing.T) {
		vec.Reset()
		_ = vec.Parse(query3)

		vec.SetSchemeString("https").
			SetUsernameString("foo").
			SetPasswordString("bar").
			SetHostnameString("google.com").
			SetPort(8080).
			SetPathString("search").
			SetQueryString("q=keys").
			SetHashString("results")
		if n := vec.Bytes(); !bytes.Equal(n, query3new) {
			t.Error("url assembly failed", "need", string(query3new), "got", string(n))
		}
	})

	t.Run("query sort", func(t *testing.T) {
		vec.Reset()
		_ = vec.ParseCopy(query0)
		mod := vec.QuerySort().QueryBytes()
		if !bytes.Equal(mod, query0sorted) {
			t.Error("query 0 sort failed", "\nneed", string(query0sorted), "\n got", string(mod))
		}
	})
}

func BenchmarkVector(b *testing.B) {
	b.Run("parse", func(b *testing.B) {
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
	})

	b.Run("parse query/0", func(b *testing.B) {
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
	})
	b.Run("parse query/1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.ParseCopy(query1)
			query := vec.Query()
			if !query.Exists("x") || query.Get("x").String() != "" {
				b.Error("query1 mismatch query param x")
			}
			if !query.Exists("z") || query.Get("z").String() != "" {
				b.Error("query1 mismatch query param z")
			}
		}
	})
	b.Run("parse query/2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.ParseCopy(query2)
			query := vec.Query()
			query.Each(func(_ int, node *vector.Node) {
				switch {
				case node.KeyString() == "b":
					if node.String() != "x" {
						b.Error("query 2 mismatch query param", node.KeyString(), "need", "x", "got", node.String())
					}
				case node.KeyString() == "arr[]":
					if node.Limit() != 3 {
						b.Error("query 2 unexpected length of param arr[]", "need", 3, "got", node.Limit())
					}
					if node.At(1).String() != "2" {
						b.Error("query 2 mismatch query param arr[1]", "need", "2", "got", node.At(1).String())
					}
				case node.KeyString() == "arr1[]":
					if node.Limit() != 3 {
						b.Error("query 2 unexpected length of param arr1[]", "need", 3, "got", node.Limit())
					}
					if node.At(0).String() != "a" {
						b.Error("query 2 mismatch query param arr1[0]", "need", "a", "got", node.At(0).String())
					}
				}
			})
		}
	})

	b.Run("set no copy", func(b *testing.B) {
		benchSet(b, false)
	})
	b.Run("set copy", func(b *testing.B) {
		benchSet(b, true)
	})

	b.Run("forget query params", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.Parse(query3)
			if y := vec.Query().GetString("y"); y != "qwerty" {
				b.Error("query 3 mismatch query param y", "need", "qwerty", "got", y)
			}
			vec.SetQueryBytes(query3repl)
			vec.Query().Each(func(_ int, node *vector.Node) {
				switch {
				case node.KeyString() == "foo":
					if node.String() != "x" {
						b.Error("query 3 (forget) mismatch query param", node.KeyString(), "need", "x", "got", node.String())
					}
				case node.KeyString() == "bar":
					if node.String() != "y" {
						b.Error("query 3 (forget) mismatch query param", node.KeyString(), "need", "y", "got", node.String())
					}
				case node.KeyString() == "a[]":
					if node.Limit() != 3 {
						b.Error("query 2 (forget) unexpected length of param a[]", "need", 3, "got", node.Limit())
					}
					if node.At(0).String() != "1" {
						b.Error("query 2 (forget) mismatch query param a[0]", "need", "1", "got", node.At(0).String())
					}
				}
			})
		}
	})

	b.Run("serialize", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.Parse(query3)

			vec.SetSchemeString("https").
				SetUsernameString("foo").
				SetPasswordString("bar").
				SetHostnameString("google.com").
				SetPort(8080).
				SetPathString("search").
				SetQueryString("q=keys").
				SetHashString("results")
			if n := vec.Bytes(); !bytes.Equal(n, query3new) {
				b.Error("url assembly failed", "need", string(query3new), "got", string(n))
			}
		}
	})

	b.Run("query sort", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.ParseCopy(query0)
			mod := vec.QuerySort().QueryBytes()
			if !bytes.Equal(mod, query0sorted) {
				b.Error("query 0 sort failed", "need", string(query0sorted), "got", string(mod))
			}
		}
	})

	b.Run("remove if", func(b *testing.B) {
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			vec.Reset()
			_ = vec.ParseStr(queryRemove)
			vec.Query().RemoveIf(func(_ int, node *vector.Node) bool {
				return strings.HasPrefix(node.KeyString(), "utm")
			})
			if vec.QuerySort().String() != queryRemoveExpect {
				b.FailNow()
			}
		}
	})
}

func benchSet(b *testing.B, cpy bool) {
	tst := testTargets{
		url: "http://marquis_warren:major@x.com/h8?x=1#anc",
		target: testTarget{
			scheme:   "http",
			username: "marquis_warren",
			password: "major",
			hostname: "x.com",
			path:     "/h8",
			query:    "?x=1",
			hash:     "#anc",
		},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		vec.Reset()
		if cpy {
			_ = vec.ParseCopy(url0)
		} else {
			_ = vec.Parse(url0)
		}
		vec.SetSchemeString("http").
			SetUsernameString("marquis_warren").
			SetPasswordString("major").
			SetHostnameString("x.com").
			SetPort(9999).
			SetPathString("/h8").
			SetQueryString("?x=1").
			SetHashString("#anc")

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
		if vec.Port() != 9999 {
			printErr(b, &tst, "port mismatch", vec.Port(), "vs", 9999)
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
