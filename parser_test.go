package v9

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	type test struct {
		name, input    string
		expectedResult Node
	}

	tests := []test{
		{name: "simple string", input: "aBc", expectedResult: &Group{
			ChildNodes: []Node{
				CharacterLiteral{Character: 'a'},
				CharacterLiteral{Character: 'B'},
				CharacterLiteral{Character: 'c'},
			},
		}},
		{name: "wildcard character", input: "ab.", expectedResult: &Group{
			ChildNodes: []Node{
				CharacterLiteral{Character: 'a'},
				CharacterLiteral{Character: 'b'},
				WildcardLiteral{},
			},
		}},
		{name: "branches", input: "ab|cd|ef", expectedResult: &Branch{ChildNodes: []Node{
			&Group{ChildNodes: []Node{
				CharacterLiteral{Character: 'a'},
				CharacterLiteral{Character: 'b'},
			}},
			&Group{ChildNodes: []Node{
				CharacterLiteral{Character: 'c'},
				CharacterLiteral{Character: 'd'},
			}},
			&Group{ChildNodes: []Node{
				CharacterLiteral{Character: 'e'},
				CharacterLiteral{Character: 'f'},
			}},
		}}},
		{name: "groups", input: "a(b|c)", expectedResult: &Group{ChildNodes: []Node{
			CharacterLiteral{Character: 'a'},
			&Branch{ChildNodes: []Node{
				&Group{ChildNodes: []Node{
					CharacterLiteral{Character: 'b'},
				}},
				&Group{ChildNodes: []Node{
					CharacterLiteral{Character: 'c'},
				}},
			}},
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{}
			tokens := lex(tt.input)

			result := p.Parse(tokens)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Fatalf("Expected:\n%+v\n\nGot:\n%+v\n", tt.expectedResult, result)
			}
		})
	}
}
