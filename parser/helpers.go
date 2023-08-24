package parser

import (
	"fmt"

	. "github.com/1x-eng/yadab/util"
)

func helpMessage(tokens []*Token, cursor uint, msg string) {
	var c *Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, msg, c.Value)
}

func parseToken(tokens []*Token, initialCursor uint, expectedKind TokenKind) (*Token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	tok := tokens[cursor]

	// Check if the token kind is as expected
	if tok.Kind != expectedKind {
		return nil, initialCursor, false
	}

	cursor++

	return tok, cursor, true
}

func parseColumnDefinitions(tokens []*Token, initialCursor uint, delimiter Token) (*[]*ColumnDefinition, uint, bool) {
	cursor := initialCursor

	cds := []*ColumnDefinition{}
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for a delimiter
		current := tokens[cursor]
		if delimiter.Equals(current) {
			break
		}

		// Look for a comma
		if len(cds) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for a column name
		id, newCursor, ok := parseToken(tokens, cursor, IdentifierKind)
		if !ok {
			helpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		// Look for a column type
		ty, newCursor, ok := parseToken(tokens, cursor, TypeKind)
		if !ok {
			helpMessage(tokens, cursor, "Expected column type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &ColumnDefinition{
			Name:     *id,
			Datatype: *ty,
		})
	}

	return &cds, cursor, true
}
