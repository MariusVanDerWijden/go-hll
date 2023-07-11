package bloom

import gohll "github.com/MariusVanDerWijden/go-hll"

type Bloom struct {
	hll *gohll.Hll
}

func NewBloom(size uint64) *Bloom {
	return &Bloom{
		hll: gohll.NewHll(size),
	}
}

func (b *Bloom) Add(item gohll.Hashable) bool {
	return b.hll.Add(item)
}

func (b *Bloom) Check(item gohll.Hashable) bool {
	return b.hll.Has(item)
}

func (b *Bloom) Union(other *Bloom) {
	b.hll.Merge(other.hll)
}
