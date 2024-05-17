package ext_sort

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func Test(t *testing.T) {
	data := [1024]uint64{}
	for i := range data {
		val, err := rand.Int(rand.Reader, big.NewInt(100000))
		if err != nil {
			panic(err)
		}
		data[i] = val.Uint64()
	}
}
