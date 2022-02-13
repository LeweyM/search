package finite_state_machine

import "testing"

func BenchmarkCompiler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Compile("abc*.(cat|dog)hello(ad(dc))") // ~3200 ns/op
	}
}
