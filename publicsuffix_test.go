package urlvector

import "testing"

func TestPublicSuffixDB(t *testing.T) {
	loadFn := func(tb testing.TB, dbFile string) {
		var (
			psdb PublicSuffixDB
			err  error
		)
		if err = psdb.Load(dbFile); err != nil {
			t.Error(err)
		}
	}
	fetchFn := func(tb testing.TB, dbURL string) {
		var (
			psdb PublicSuffixDB
			err  error
		)
		if err = psdb.Fetch(dbURL); err != nil {
			t.Error(err)
		}
	}
	t.Run("load small", func(t *testing.T) { loadFn(t, "testdata/small.psdb") })
	t.Run("load full", func(t *testing.T) { loadFn(t, "testdata/full.psdb") })
	t.Run("fetch small", func(t *testing.T) {
		fetchFn(t, "https://raw.githubusercontent.com/koykov/urlvector/master/testdata/small.psdb")
	})
	t.Run("fetch full", func(t *testing.T) {
		fetchFn(t, "https://raw.githubusercontent.com/koykov/urlvector/master/testdata/full.psdb")
	})
}
