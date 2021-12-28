package urlvector

type pse uint64

func (e *pse) encode(lo, hi uint32) {
	*e = pse(lo)<<32 | pse(hi)
}

func (e pse) decode() (lo, hi uint32) {
	lo = uint32(e >> 32)
	hi = uint32(e & 0xffffffff)
	return
}
