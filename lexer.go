package main

import (
	"fmt"
	"strings"
)

type location struct {
	line uint
	col  uint
}

type keyword string

const (
	selectKeyword keyword = "select"
	fromKeyword   keyword = "from"
	asKeyword     keyword = "as"
	tableKeyword  keyword = "table"
	createKeyword keyword = "create"
	insertKeyword keyword = "insert"
	intoKeyword   keyword = "into"
	valuesKeyword keyword = "values"
	intKeyword    keyword = "int"
	textKeyword   keyword = "text"
	whereKeyword  keyword = "where"
)

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftParenSymbol  symbol = "("
	rightParenSymbol symbol = ")"
	eqSymbol         symbol = "="
	neqSymbol        symbol = "<>"
	neqSymbol2       symbol = "!="
	concatSymbol     symbol = "||"
	plusSymbol       symbol = "+"
	ltSymbol         symbol = "<"
	lteSymbol        symbol = "<="
	gtSymbol         symbol = ">"
	gteSymbol        symbol = ">="
)

type tokenKind uint

const (
	keywordKind tokenKind = iota
	symbolKind
	identifierKind
	stringKind
	numericKind
)

type token struct {
	value string
	kind  tokenKind
	loc   location
}

type cursor struct {
	pointer uint
	loc     location
}

func (t *token) equals(other *token) bool {
	return t.value == other.value && t.kind == other.kind
}

type lexer func(string, cursor) (*token, cursor, bool)

// lexNumeric is a Go function that performs lexical analysis on a numeric value in a given source string.
//
// It takes two parameters:
// - source: a string representing the source code to be analyzed.
// - ic: a cursor representing the initial position in the source string.
//
// It returns three values:
// - token: a pointer to a token struct representing the analyzed token.
// - cursor: a cursor representing the updated position in the source string after analysis.
// - bool: a boolean value indicating whether the analysis was successful or not.
func lexNumeric(source string, ic cursor) (*token, cursor, bool) {
	cur := ic
	periodFound := false
	expMarkerFound := false

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		cur.loc.col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExpMarker := c == 'e'

		if cur.pointer == ic.pointer {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}
			periodFound = isPeriod
			continue
		}

		if isPeriod && periodFound {
			return nil, ic, false
		}

		if isExpMarker && expMarkerFound {
			return nil, ic, false
		}

		if isPeriod {
			periodFound = true
		}

		if isExpMarker {
			periodFound = true
			expMarkerFound = true

			if cur.pointer == uint(len(source)-1) {
				return nil, ic, false
			}

			cNext := source[cur.pointer+1]
			if cNext == '-' || cNext == '+' {
				cur.pointer++
				cur.loc.col++
			}
		}

		if !isDigit {
			break
		}
	}

	if cur.pointer == ic.pointer {
		return nil, ic, false
	}

	return &token{
		value: source[ic.pointer:cur.pointer],
		loc:   ic.loc,
		kind:  numericKind,
	}, cur, true
}

// lexCharacterDelimited scans a source string starting from the given cursor position
// and looks for a character delimiter.
//
// It returns a token, the updated cursor position, and a boolean indicating whether
// a delimiter was found.
//
// Parameters:
// - source: The source string to scan.
// - ic: The initial cursor position.
// - delimiter: The character delimiter to look for.
//
// Returns:
// - The token representing the delimited value.
// - The updated cursor position.
// - A boolean indicating whether a delimiter was found.
func lexCharacterDelimited(source string, ic cursor, delimiter byte) (*token, cursor, bool) {
	cur := ic

	if len(source[cur.pointer:]) == 0 || source[cur.pointer] != delimiter {
		return nil, ic, false
	}

	cur.loc.col++
	cur.pointer++

	var value []byte
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter && (cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter) {
			return &token{
				value: string(value),
				loc:   ic.loc,
				kind:  stringKind,
			}, cur, true
		}

		value = append(value, c)
		cur.loc.col++
	}

	return nil, ic, false
}

// lexString is a function that lexes a string literal in the given source code.
//
// It takes the source code as a string and the initial cursor position as input.
// It returns a token, a new cursor position, and a boolean value indicating whether lexing was successful.
func lexString(source string, ic cursor) (*token, cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}

// lexSymbol is a Go function that performs lexical analysis on a given source string and cursor position.
//
// Parameters:
// - source: the source string to analyze.
// - ic: the initial cursor position.
//
// Returns:
// - token: the token found during analysis.
// - cursor: the updated cursor position.
// - bool: a boolean indicating whether the token is valid or not.
func lexSymbol(source string, ic cursor) (*token, cursor, bool) {
	c := source[ic.pointer]
	cur := ic

	// Will get overwritten later if not an ignored syntax
	cur.pointer++
	cur.loc.col++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.loc.line++
		cur.loc.col = 0
		fallthrough
	case '\t', ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	symbols := []symbol{
		commaSymbol,
		leftParenSymbol,
		rightParenSymbol,
		semicolonSymbol,
		asteriskSymbol,
	}

	options := make([]string, len(symbols))
	for i, s := range symbols {
		options[i] = string(s)
	}

	// Use `ic`, not `cur`
	match := longestMatch(source, ic, options)

	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		value: match,
		loc:   ic.loc,
		kind:  symbolKind,
	}, cur, true
}

// lexKeyword lexes a keyword from the source string starting at the given cursor position.
//
// Parameters:
// - source: The source string to lex.
// - ic: The initial cursor position.
//
// Returns:
// - token: The lexed token.
// - cursor: The updated cursor position.
// - bool: Indicates if a keyword was successfully lexed.
func lexKeyword(source string, ic cursor) (*token, cursor, bool) {
	cur := ic
	keywords := []keyword{
		selectKeyword,
		insertKeyword,
		valuesKeyword,
		tableKeyword,
		createKeyword,
		whereKeyword,
		fromKeyword,
		intoKeyword,
		textKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.col = ic.loc.col + uint(len(match))

	return &token{
		value: match,
		kind:  keywordKind,
		loc:   ic.loc,
	}, cur, true
}

// longestMatch finds the longest match in a given source string, based on a list of options.
//
// Parameters:
// - source: the source string to search within.
// - ic: the initial cursor position.
// - options: a list of strings representing the options to match against.
//
// Returns:
// - string: the longest match found from the options list.
func longestMatch(source string, ic cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.pointer < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.pointer]))...)
		cur.pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.pointer-ic.pointer]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

// lexIdentifier is a Go function that analyzes a source string and a cursor position to identify an identifier token.
//
// It takes two parameters:
// - source: a string representing the source code to be analyzed.
// - ic: a cursor object representing the current position in the source code.
//
// It returns three values:
// - token: a pointer to a token struct representing the identified token.
// - cursor: a cursor object representing the updated position in the source code.
// - bool: a boolean value indicating whether an identifier token was successfully identified.
func lexIdentifier(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	c := source[cur.pointer]
	if !isAlphabetical(c) {
		return nil, ic, false
	}
	cur.pointer++
	cur.loc.col++

	value := []byte{c}
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c = source[cur.pointer]

		if !isAlphabetical(c) && !isNumeric(c) && !isSpecialCharacter(c) {
			break
		}

		value = append(value, c)
		cur.loc.col++
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &token{
		value: strings.ToLower(string(value)),
		loc:   ic.loc,
		kind:  identifierKind,
	}, cur, true
}

func isAlphabetical(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func isSpecialCharacter(c byte) bool {
	return c == '$' || c == '_'
}

func lex(source string) ([]*token, error) {
	tokens := []*token{}
	cur := cursor{}

	for cur.pointer < uint(len(source)) {
		lexers := []lexer{lexString, lexNumeric, lexSymbol, lexKeyword, lexIdentifier}
		found := false
		for _, l := range lexers {
			if token, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				// Omit nil tokens for valid, but empty syntax like newlines
				if token != nil {
					tokens = append(tokens, token)
				}

				found = true
				break
			}
		}

		if !found {
			hint := ""
			if len(tokens) > 0 {
				hint = " after " + tokens[len(tokens)-1].value
			}
			return nil, fmt.Errorf("could not lex unrecognized input at line %d, column %d%s", cur.loc.line+1, cur.loc.col+1, hint)
		}
	}

	return tokens, nil
}
