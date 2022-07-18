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
			Expressions: []Ast{
				CharacterLiteral{Character: 'a'},
				CharacterLiteral{Character: 'B'},
				CharacterLiteral{Character: 'c'},
			},
		}},
		{name: "modifiers", input: "a+b?c*", expectedResult: &Group{
			Expressions: []Ast{
				ModifierExpression{Expression: CharacterLiteral{Character: 'a'}, Modifier: OneOrManyModifier},
				ModifierExpression{Expression: CharacterLiteral{Character: 'b'}, Modifier: ZeroOrOneModifier},
				ModifierExpression{Expression: CharacterLiteral{Character: 'c'}, Modifier: ZeroOrManyModifier},
			},
		}},
		{name: "branches", input: "a|bc|d", expectedResult: &Branch{[]Ast{
			&Branch{[]Ast{
				&Group{[]Ast{CharacterLiteral{Character: 'a'}}},
				&Group{[]Ast{
					CharacterLiteral{Character: 'b'},
					CharacterLiteral{Character: 'c'},
				}},
			}},
			&Group{[]Ast{CharacterLiteral{Character: 'd'}}},
		}}},
		{name: "groups", input: "(a)(bc)*(d)", expectedResult: &Group{
			Expressions: []Ast{
				&Group{Expressions: []Ast{CharacterLiteral{Character: 'a'}}},
				ModifierExpression{
					Expression: &Group{Expressions: []Ast{
						CharacterLiteral{Character: 'b'},
						CharacterLiteral{Character: 'c'},
					}},
					Modifier: ZeroOrManyModifier,
				},
				&Group{Expressions: []Ast{CharacterLiteral{Character: 'd'}}},
			},
		}},
		{name: "all together", input: "(cat|(dog)(s)?)*", expectedResult: &Group{
			Expressions: []Ast{
				ModifierExpression{
					Expression: &Branch{Expressions: []Ast{
						&Group{Expressions: []Ast{
							CharacterLiteral{Character: 'c'},
							CharacterLiteral{Character: 'a'},
							CharacterLiteral{Character: 't'},
						}},
						&Group{Expressions: []Ast{
							&Group{Expressions: []Ast{
								CharacterLiteral{Character: 'd'},
								CharacterLiteral{Character: 'o'},
								CharacterLiteral{Character: 'g'},
							}},
							ModifierExpression{
								Expression: &Group{Expressions: []Ast{
									CharacterLiteral{Character: 's'},
								}},
								Modifier: ZeroOrOneModifier,
							},
						}},
					}},
					Modifier: ZeroOrManyModifier,
				},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Parser{}

			result := p.Parse(tt.input)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Fatalf("Expected [%+v], got [%+v]", tt.expectedResult, result)
			}
		})
	}
}
