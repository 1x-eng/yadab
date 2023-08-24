package lexer

import (
	"fmt"

	. "github.com/1x-eng/yadab/util"
)

func Lex(source string) ([]*Token, error) {
	tokens := []*Token{}
	cur := Cursor{}

lex:
	for cur.Pointer < uint(len(source)) {
		lexers := []Lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexType, lexIdentifier}
		for _, l := range lexers {
			if token, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				// Omit nil tokens for valid, but empty syntax like newlines
				if token != nil {
					tokens = append(tokens, token)
				}

				continue lex
			}
		}

		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].Value
		}
		return nil, fmt.Errorf("Unable to lex token%s, at %d:%d", hint, cur.Loc.Line, cur.Loc.Col)
	}

	return tokens, nil
}
