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
		{name: "simple string", input: "abc", expectedResult: Group{
			expressions: []Ast{
				CharacterLiteral{character: 'a'},
				CharacterLiteral{character: 'b'},
				CharacterLiteral{character: 'c'},
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
