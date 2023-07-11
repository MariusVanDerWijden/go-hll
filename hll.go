package gohll

import (
	"encoding/binary"
	"math/bits"
)

type Hashable interface {
	Hash() [32]byte
}

type Metrics struct {
	Inserts     int
	Updates     int
	Cardinality float64
	Buckets     []byte
}

type Hll struct {
	buckets []byte
	inserts int
	updates int
}

func NewHll(bucketCount uint64) *Hll {
	return &Hll{
		buckets: make([]byte, bucketCount),
	}
}

func (h *Hll) Add(item Hashable) bool {
	hash := item.Hash()
	// calculate bucket index
	bucket := int(binary.BigEndian.Uint64(hash[:8])) % len(h.buckets)
	// calculate trailing zeros
	trailingZeros := ctz(hash)
	inserted := h.insertBucket(bucket, trailingZeros)
	// update metrics
	if inserted {
		h.inserts++
	}
	h.updates++
	return inserted
}

func (h *Hll) Has(item Hashable) bool {
	hash := item.Hash()
	// calculate bucket index
	bucket := int(binary.BigEndian.Uint64(hash[:8])) % len(h.buckets)
	// calculate trailing zeros
	trailingZeros := ctz(hash)
	return h.buckets[bucket] >= trailingZeros
}

func (h *Hll) Merge(other *Hll) bool {
	h.updates += other.updates
	// inserts are deliberately not added here.
	var inserted bool
	for index, bucket := range other.buckets {
		inserted = inserted || h.insertBucket(index, bucket)
	}
	return inserted
}

func (h *Hll) Stats() *Metrics {
	cardinality := cardinality(h.buckets)
	return &Metrics{
		Inserts:     h.inserts,
		Updates:     h.updates,
		Cardinality: cardinality,
		Buckets:     h.buckets,
	}
}

func (h *Hll) insertBucket(index int, element byte) bool {
	if h.buckets[index] < element {
		h.buckets[index] = element
		return true
	}
	return false
}

// TODO for counting hashes, it might be faster to do this with uint64s
func ctz(n [32]byte) byte {
	tlz := 0
	for i := 32; i > 0; i-- {
		zeros := bits.TrailingZeros8(n[i-1])
		tlz += zeros
		if zeros != 8 {
			break
		}
	}
	return byte(tlz)
}

// TODO this cardinality computation is very basic and has a large error
func cardinality(buckets []byte) float64 {
	sum := float64(0)
	for _, bucket := range buckets {
		sum += float64(1) / float64(bucket)
	}
	return float64(len(buckets)) / sum
}
