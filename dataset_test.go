package urlvector

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var dsStages []struct {
	idx int
	ds  string
	src string
}

func init() {
	files, err := filepath.Glob("dataset/*.txt")
	if err != nil {
		return
	}
	for i := 0; i < len(files); i++ {
		file, err := os.Open(files[i])
		if err != nil {
			continue
		}
		fn := filepath.Base(files[i])
		ds := strings.TrimSuffix(fn, filepath.Ext(fn))

		scr := bufio.NewScanner(file)
		var c int
		for scr.Scan() {
			dsStages = append(dsStages, struct {
				idx     int
				ds, src string
			}{idx: c, ds: ds, src: scr.Text()})
			c++
		}
		_ = file.Close()
	}
}

func TestDataset(t *testing.T) {
	vec := NewVector()
	for i := 0; i < len(dsStages); i++ {
		stg := &dsStages[i]
		err := vec.ParseCopyStr(stg.src)
		if err != nil {
			t.Error(err)
		}
		_ = vec.Query()
		vec.Reset()
	}
}

func BenchmarkDataset(b *testing.B) {
	if len(dsStages) == 0 {
		return
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		stg := &dsStages[i%len(dsStages)]
		vec := Acquire()
		_ = vec.ParseCopyStr(stg.src)
		_ = vec.Query()
		Release(vec)
	}
}
