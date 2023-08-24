package lexer

import (
	. "github.com/1x-eng/yadab/util"
)

func lexKeyword(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic
	keywords := []Keyword{
		SelectKeyword,
		InsertKeyword,
		ValuesKeyword,
		TableKeyword,
		CreateKeyword,
		WhereKeyword,
		FromKeyword,
		IntoKeyword,
		TextKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.Pointer = ic.Pointer + uint(len(match))
	cur.Loc.Col = ic.Loc.Col + uint(len(match))

	kind := KeywordKind
	if match == string(IntKeyword) || match == string(TextKeyword) {
		kind = TypeKind
	}

	return &Token{
		Value: match,
		Kind:  kind,
		Loc:   ic.Loc,
	}, cur, true
}
