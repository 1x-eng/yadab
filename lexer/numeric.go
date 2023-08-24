package lexer

import (
	. "github.com/1x-eng/yadab/util"
)

func lexNumeric(source string, ic Cursor) (*Token, Cursor, bool) {
	cur := ic

	periodFound := false
	expMarkerFound := false

	for ; cur.Pointer < uint(len(source)); cur.Pointer++ {
		c := source[cur.Pointer]
		cur.Loc.Col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExpMarker := c == 'e'

		// Must start with a digit or period
		if cur.Pointer == ic.Pointer {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			if periodFound {
				return nil, ic, false
			}

			periodFound = true
			continue
		}

		if isExpMarker {
			if expMarkerFound {
				return nil, ic, false
			}

			// No periods allowed after expMarker
			periodFound = true
			expMarkerFound = true

			// expMarker must be followed by digits
			if cur.Pointer == uint(len(source)-1) {
				return nil, ic, false
			}

			cNext := source[cur.Pointer+1]
			if cNext == '-' || cNext == '+' {
				cur.Pointer++
				cur.Loc.Col++
			}

			continue
		}

		if !isDigit {
			break
		}
	}

	// No characters accumulated
	if cur.Pointer == ic.Pointer {
		return nil, ic, false
	}

	return &Token{
		Value: source[ic.Pointer:cur.Pointer],
		Loc:   ic.Loc,
		Kind:  NumericKind,
	}, cur, true
}
