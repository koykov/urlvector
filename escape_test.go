package urlvector

import (
	"strconv"
	"strings"
	"testing"

	"github.com/koykov/byteconv"
)

var stages = []struct {
	raw, exp, exp1 string
}{
	{"", "", ""},
	{"abc", "abc", "abc"},
	{"10%", "10%25", "10%25"},
	{
		" ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;",
		"+%3F%26%3D%23%2B%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09%3A%2F%40%24%27%28%29%2A%2C%3B",
		"%20%3F&=%23+%25%21%3C%3E%23%22%7B%7D%7C%5C%5E%5B%5D%60%E2%98%BA%09:%2F@$%27%28%29%2A%2C%3B",
	},
	{
		"one two",
		"one+two",
		"one%20two",
	},
	{
		"Фотки собак",
		"%D0%A4%D0%BE%D1%82%D0%BA%D0%B8+%D1%81%D0%BE%D0%B1%D0%B0%D0%BA",
		"%D0%A4%D0%BE%D1%82%D0%BA%D0%B8%20%D1%81%D0%BE%D0%B1%D0%B0%D0%BA",
	},

	{
		"shortrun(break)shortrun",
		"shortrun%28break%29shortrun",
		"shortrun%28break%29shortrun",
	},

	{
		"longerrunofcharacters(break)anotherlongerrunofcharacters",
		"longerrunofcharacters%28break%29anotherlongerrunofcharacters",
		"longerrunofcharacters%28break%29anotherlongerrunofcharacters",
	},
	{
		strings.Repeat("padded/with+various%characters?that=need$some@escaping+paddedsowebreak/256bytes", 4),
		strings.Repeat("padded%2Fwith%2Bvarious%25characters%3Fthat%3Dneed%24some%40escaping%2Bpaddedsowebreak%2F256bytes", 4),
		strings.Repeat("padded%2Fwith+various%25characters%3Fthat=need$some@escaping+paddedsowebreak%2F256bytes", 4),
	},
}

func TestEscape(t *testing.T) {
	for i, stage := range stages {
		t.Run("query/"+strconv.Itoa(i), func(t *testing.T) {
			var buf []byte
			buf = QueryEscape(buf[:0], byteconv.S2B(stage.raw))
			if r := byteconv.B2S(buf); r != stage.exp {
				t.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp, r)
			}
		})
	}
	for i, stage := range stages {
		t.Run("path/"+strconv.Itoa(i), func(t *testing.T) {
			var buf []byte
			buf = PathEscape(buf[:0], byteconv.S2B(stage.raw))
			if r := byteconv.B2S(buf); r != stage.exp1 {
				t.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp1, r)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	for i, stage := range stages {
		t.Run("query/"+strconv.Itoa(i), func(t *testing.T) {
			var buf []byte
			buf = QueryUnescape(buf[:0], byteconv.S2B(stage.exp))
			if r := byteconv.B2S(buf); r != stage.raw {
				t.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp, r)
			}
		})
	}
	for i, stage := range stages {
		t.Run("path/"+strconv.Itoa(i), func(t *testing.T) {
			if len(stage.exp1) == 0 {
				return
			}
			var buf []byte
			buf = PathUnescape(buf[:0], byteconv.S2B(stage.exp1))
			if r := byteconv.B2S(buf); r != stage.raw {
				t.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.raw, r)
			}
		})
	}
}

func BenchmarkEscape(b *testing.B) {
	for i, stage := range stages {
		b.Run("query/"+strconv.Itoa(i), func(b *testing.B) {
			var buf []byte
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = QueryEscape(buf[:0], byteconv.S2B(stage.raw))
				if r := byteconv.B2S(buf); r != stage.exp {
					b.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp, r)
				}
			}
		})
	}
	for i, stage := range stages {
		b.Run("path/"+strconv.Itoa(i), func(b *testing.B) {
			var buf []byte
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = PathEscape(buf[:0], byteconv.S2B(stage.raw))
				if r := byteconv.B2S(buf); r != stage.exp1 {
					b.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp1, r)
				}
			}
		})
	}
}

func BenchmarkUnescape(b *testing.B) {
	for i, stage := range stages {
		b.Run("query/"+strconv.Itoa(i), func(b *testing.B) {
			var buf []byte
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = QueryUnescape(buf[:0], byteconv.S2B(stage.exp))
				if r := byteconv.B2S(buf); r != stage.raw {
					b.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.exp, r)
				}
			}
		})
	}
	for i, stage := range stages {
		b.Run("path/"+strconv.Itoa(i), func(b *testing.B) {
			var buf []byte
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				buf = PathUnescape(buf[:0], byteconv.S2B(stage.exp1))
				if r := byteconv.B2S(buf); r != stage.raw {
					b.Errorf("escape mismatch:\n\tneed '%s'\n\tgot  '%s'", stage.raw, r)
				}
			}
		})
	}
}
