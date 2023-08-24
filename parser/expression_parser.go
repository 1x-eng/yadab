package parser

import . "github.com/1x-eng/yadab/util"

func parseExpressions(tokens []*Token, initialCursor uint, delimiters []Token) (*[]*Expression, uint, bool) {
	cursor := initialCursor

	exps := []*Expression{}
outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for delimiter
		current := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.Equals(current) {
				break outer
			}
		}

		// Look for comma
		if len(exps) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(CommaSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*Token, initialCursor uint, _ Token) (*Expression, uint, bool) {
	cursor := initialCursor

	kinds := []TokenKind{IdentifierKind, NumericKind, StringKind}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &Expression{
				Literal: t,
				Kind:    LiteralKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}
