package indigo

import (
	"fmt"
	"strconv"
)

type IndigoReader interface {
	// Returns the next-parsed list
	Read() (any, error)
}

type StandardReader struct {
	lexer   Lexer[LexemeKind]
	curr    *Token[LexemeKind]
	currErr error
}

type LexemeKind int

const (
	LParen LexemeKind = iota
	RParen
	Integer
)

type ProductionKind int

const (
	ProdInteger ProductionKind = iota
	ProdList
)

var IndigoLexer = MustNewLexerFactory(map[string]LexemeKind{
	`\(`:  LParen,
	`\)`:  RParen,
	`\d+`: Integer,
}, `\s+`)

func NewStandardReader(source string) IndigoReader {

	lexer := IndigoLexer(source)

	curr, currErr := lexer.Next()

	return &StandardReader{
		lexer:   lexer,
		curr:    curr,
		currErr: currErr,
	}
}

func (sr *StandardReader) Next() error {
	sr.curr, sr.currErr = sr.lexer.Next()
	return sr.currErr
}

func (sr *StandardReader) Read() (any, error) {

	if sr.currErr != nil {
		return nil, sr.currErr
	}

	switch sr.curr.kind {
	case Integer:
		integer, err := strconv.ParseInt(sr.curr.lexeme, 10, 64)
		if err != nil {
			return nil, err
		}
		sr.Next()
		return integer, nil
	case LParen:
		cons, err := sr.readList()
		if err != nil {
			return nil, fmt.Errorf("invalid list: %e", err)
		}
		return cons, nil
	}

	panic("unreachable")
}

func (sr *StandardReader) readList() (*Cons, error) {
	err := sr.Next()
	if err != nil {
		panic(err)
	}

	var head *Cons = nil

	tail := head

	for sr.curr.kind != RParen {
		value, err := sr.Read()

		if err == EOF {
			return nil, fmt.Errorf("unbalanced parentheses: missing right parenthesis")
		} else if err != nil {
			return nil, err
		}

		*tail = NewCons(value, nil)
		tail = tail.cdr.(*Cons)
	}

	return head, nil
}
