package parser

import . "github.com/1x-eng/yadab/util"

func tokenFromKeyword(k Keyword) Token {
	return Token{
		Kind:  KeywordKind,
		Value: string(k),
	}
}

func tokenFromSymbol(s Symbol) Token {
	return Token{
		Kind:  SymbolKind,
		Value: string(s),
	}
}

func expectToken(tokens []*Token, cursor uint, t Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}

	return t.Equals(tokens[cursor])
}
