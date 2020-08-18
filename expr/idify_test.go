package expr

import (
	"encoding/base32"
	"hash/fnv"
	"testing"
)

func BenchmarkIdify(b *testing.B) {
	sample := []string{"0", "medium", "a longer string", "a super long duper long very very long string"}
	for n := 0; n < b.N; n++ {
		for _, s := range sample {
			idify(s)
		}
	}
}

func BenchmarkBase32(b *testing.B) {
	var h = fnv.New32a()
	sample := []string{"0", "medium", "a longer string", "a super long duper long very very long string"}
	for n := 0; n < b.N; n++ {
		for _, s := range sample {
			h.Reset()
			h.Write([]byte(s))
			b := h.Sum(nil)
			base32.StdEncoding.EncodeToString(b)
		}
	}
}
