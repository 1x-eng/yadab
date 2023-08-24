package lexer

import (
	. "github.com/1x-eng/yadab/util"
)

func lexCharacterDelimited(source string, ic Cursor, delimiter byte) (*Token, Cursor, bool) {
	cur := ic

	if len(source[cur.Pointer:]) == 0 || source[cur.Pointer] != delimiter {
		return nil, ic, false
	}

	cur.Loc.Col++
	cur.Pointer++

	var value []byte
	for ; cur.Pointer < uint(len(source)); cur.Pointer++ {
		c := source[cur.Pointer]

		if c == delimiter {
			if cur.Pointer+1 >= uint(len(source)) || source[cur.Pointer+1] != delimiter {
				cur.Pointer++ // consume the closing delimiter
				cur.Loc.Col++
				return &Token{
					Value: string(value),
					Loc:   ic.Loc,
					Kind:  StringKind,
				}, cur, true
			} else {
				value = append(value, delimiter)
				cur.Pointer++
				cur.Loc.Col++
			}
		}

		value = append(value, c)
		cur.Loc.Col++
	}

	return nil, ic, false
}

func lexString(source string, ic Cursor) (*Token, Cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}
