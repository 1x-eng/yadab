package util

type TokenKind uint

type ExpressionKind uint

type Keyword string

type Symbol string

type AstKind uint

type Location struct {
	Line uint
	Col  uint
}

type Token struct {
	Value string
	Kind  TokenKind
	Loc   Location
}

type SelectStatement struct {
	Item []*Expression
	From Token
}

type Expression struct {
	Literal *Token
	Kind    ExpressionKind
}

type CreateTableStatement struct {
	Name Token
	Cols *[]*ColumnDefinition
}

type InsertStatement struct {
	Table  Token
	Values *[]*Expression
}

type ColumnDefinition struct {
	Name     Token
	Datatype Token
}

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

type Ast struct {
	Statements []*Statement
}

type Cursor struct {
	Pointer uint
	Loc     Location
}

type Lexer func(string, Cursor) (*Token, Cursor, bool)
