package urlvector

import (
	"sync"

	"github.com/koykov/vector"
)

// Pool represents URL vectors pool.
type Pool struct {
	p sync.Pool
}

var (
	// p is a default instance of the pool.
	// Just call urlvector.Acquire() and urlvector.Release().
	p Pool
)

// Get old vector from the pool or create new one.
func (p *Pool) Get() *Vector {
	v := p.p.Get()
	if v != nil {
		if vec, ok := v.(*Vector); ok {
			vec.Helper = helper
			return vec
		}
	}
	return NewVector()
}

// Put vector back to the pool.
func (p *Pool) Put(vec *Vector) {
	vec.Reset()
	p.p.Put(vec)
}

// Acquire returns vector from default pool instance.
func Acquire() *Vector {
	return p.Get()
}

// Release puts vector back to default pool instance.
func Release(vec *Vector) {
	p.Put(vec)
}

// ReleaseNC puts vector back to pool with enforced no-clear flag.
func ReleaseNC(vec *Vector) {
	vec.SetBit(vector.FlagNoClear, true)
	p.Put(vec)
}

var _, _, _ = Acquire, Release, ReleaseNC
