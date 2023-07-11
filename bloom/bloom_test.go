package bloom

import (
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"testing"

	bloomfilter "github.com/holiman/bloomfilter/v2"
)

var rand = mrand.New(mrand.NewSource(1))

type TestElement [32]byte

func (t TestElement) Hash() [32]byte {
	return t
}

func (f TestElement) Write(p []byte) (n int, err error) { panic("not implemented") }
func (f TestElement) Sum(b []byte) []byte               { panic("not implemented") }
func (f TestElement) Reset()                            { panic("not implemented") }
func (f TestElement) BlockSize() int                    { panic("not implemented") }
func (f TestElement) Size() int                         { return 8 }
func (f TestElement) Sum64() uint64                     { return binary.BigEndian.Uint64(f[:8]) }

func TestBloomCollision(t *testing.T) {
	size := uint64(2048)
	hllBloom := NewBloom(size)
	bloom, err := bloomfilter.New(size*8, 4)
	if err != nil {
		t.Fatal(err)
	}

	var (
		hll int
		b   int
	)

	for i := 0; i < 10000; i++ {
		var input [32]byte
		rand.Read(input[:])

		elem := TestElement(input)
		hllAdded := hllBloom.Check(elem)
		bloomAdded := bloom.Contains(elem)
		if hllAdded {
			hll++
		}
		if bloomAdded {
			b++
		}
		bloom.Add(elem)
		hllBloom.Add(elem)
	}

	fmt.Printf("%v %v \n", hll, b)
	panic("asf")
}

func BenchmarkHll(b *testing.B) {
	size := uint64(2048)
	hllBloom := NewBloom(size)

	for i := 0; i < b.N; i++ {
		var input [32]byte
		rand.Read(input[:])

		elem := TestElement(input)
		hllBloom.Check(elem)
		hllBloom.Add(elem)
	}
}

func BenchmarkBloom(b *testing.B) {
	size := uint64(2048)
	bloom, _ := bloomfilter.New(size*8, 4)

	for i := 0; i < b.N; i++ {
		var input [32]byte
		rand.Read(input[:])

		elem := TestElement(input)
		bloom.Contains(elem)

		bloom.Add(elem)
	}
}
