package finite_state_machine

import (
	"search/src/ast"
	"testing"
)

func BenchmarkCompiler(b *testing.B) {
	// ~3200 ns/op = old compiler

	for i := 0; i < b.N; i++ {
		parser := ast.Parser{}
		tree := parser.Parse("abc*.(cat|dog)hello(ad(dc))")
		CompileNEW(tree) // 6400 ns/op
	}
}
