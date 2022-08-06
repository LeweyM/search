package graphdraw

import "testing"

func TestGraphDraw(t *testing.T) {
	expected := `
┌─────┐
│  A  │
└─────┘
`
	result := Draw(node{label: "A"})

	if expected != result {
		t.Fatalf("Expected %s, got %s", expected, result)
	}

}

func TestGraphDraw2(t *testing.T) {
	expected := `
┌───────┐
│  BIG  │
└───────┘
`
	result := Draw(node{label: "BIG"})

	if expected != result {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}

func TestGraphDraw3(t *testing.T) {
	expected := `
┌─────┐      ┌─────┐  
│  A  │─────>│  B  │
└─────┘      └─────┘
`
	a := node{label: "A"}
	b := node{label: "B"}
	a.to(&b)
	result := Draw(a)

	if expected != result {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}
