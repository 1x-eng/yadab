package lexer

import (
	. "github.com/1x-eng/yadab/util"
)

func lexSymbol(source string, ic Cursor) (*Token, Cursor, bool) {
	c := source[ic.Pointer]
	cur := ic
	cur.Pointer++
	cur.Loc.Col++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.Loc.Line++
		cur.Loc.Col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	symbols := []Symbol{
		CommaSymbol,
		LeftParenSymbol,
		RightParenSymbol,
		SemicolonSymbol,
		AsteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	// Use `ic`, not `cur`
	match := longestMatch(source, ic, options)
	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	return &Token{
		Value: match,
		Loc:   ic.Loc,
		Kind:  SymbolKind,
	}, cur, true
}
