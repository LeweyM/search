package graphdraw

import "testing"

func TestShapeDraw(t *testing.T) {
	expected := `
┌─────┐
│     │
└─────┘
`
	sq := Square{2, 2, 1}
	b := Board{shapes: []Square{sq}}
	result := b.draw()

	if expected != result {
		t.Fatalf("Expected %s, got %s", expected, result)
	}

}
