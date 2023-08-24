package lexer

import (
	. "github.com/1x-eng/yadab/util"
)

func lexType(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic
	types := []Keyword{
		IntKeyword,
		TextKeyword,
	}

	var options []string
	for _, t := range types {
		options = append(options, string(t))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Kind:  TypeKind,
		Loc:   ic.Loc,
	}, cur, true
}
