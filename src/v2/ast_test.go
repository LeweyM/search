package v2

import "testing"

func TestAstGroupString(t *testing.T) {
	group := Group{
		Expressions: []Ast{
			CharacterLiteral{Character: 'a'},
			CharacterLiteral{Character: 'b'},
			CharacterLiteral{Character: 'c'},
		},
	}

	expected := `
Group {[
----Literal {a}
---- 
----Literal {b}
---- 
----Literal {c}
----]
}`

	if group.String() != expected {
		t.Fatalf("Expected: %s, got: %s", expected, group.String())
	}
}
