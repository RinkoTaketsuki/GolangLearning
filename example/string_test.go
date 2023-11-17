package example

import (
	"testing"
)

func BenchmarkStringAssign(b *testing.B) {
	b.StopTimer()
	sli := make([]string, 10000000)
	str := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"
	b.StartTimer()
	for i := range sli {
		sli[i] = str
	}
}

func BenchmarkByteSliceCopy(b *testing.B) {
	b.StopTimer()
	sli := make([][78]byte, 10000000)
	str := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"
	b.StartTimer()
	for i := range sli {
		copy(sli[i][:], []byte(str))
	}
}
