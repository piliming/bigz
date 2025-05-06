package uint512

import (
	"testing"
)

func BenchmarkLsh(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand512slice(K)

	b.Run("Lsh_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xx[i%K].Lsh(uint(i % 520))
		}
	})

	//// Native: 64 - 64
	//b.Run("Lsh1_512", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		_ = xx[i%K].Lsh1(uint(i % 520))
	//	}
	//})

}

func BenchmarkRsh(b *testing.B) {
	const K = 1024 // should be power of 2
	xx := rand512slice(K)

	b.Run("Rsh_512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = xx[i%K].Rsh(uint(i % 520))
		}
	})

	// Native: 64 - 64
	//b.Run("Rsh1_512", func(b *testing.B) {
	//	for i := 0; i < b.N; i++ {
	//		_ = xx[i%K].Rsh1(uint(i % 520))
	//	}
	//})

}
