package ast

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	type test struct {
		name, input    string
		expectedResult Ast
	}

	tests := []test{
		{name: "simple string", input: "aBc", expectedResult: &Group{
			expressions: []Ast{
				CharacterLiteral{character: 'a'},
				CharacterLiteral{character: 'B'},
				CharacterLiteral{character: 'c'},
			},
		}},
		{name: "modifiers", input: "a+b?c*", expectedResult: &Group{
			expressions: []Ast{
				ModifierExpression{expression: CharacterLiteral{character: 'a'}, modifier: OneOrMany},
				ModifierExpression{expression: CharacterLiteral{character: 'b'}, modifier: zeroOrOne},
				ModifierExpression{expression: CharacterLiteral{character: 'c'}, modifier: zeroOrMany},
			},
		}},
		{name: "branches", input: "a|bc|d", expectedResult: &Branch{
			expressions: []Ast{
				&Group{expressions: []Ast{CharacterLiteral{character: 'a'}}},
				&Group{expressions: []Ast{
					CharacterLiteral{character: 'b'},
					CharacterLiteral{character: 'c'},
				}},
				&Group{expressions: []Ast{CharacterLiteral{character: 'd'}}},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser{}

			result := p.parse(tt.input)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Fatalf("Expected [%+v], got [%+v]", tt.expectedResult, result)
			}
		})
	}
}
