package reading

import (
	"fmt"
	"strconv"

	"github.com/joshua-zingale/indigo/indigo/interfaces"
	"github.com/joshua-zingale/indigo/indigo/internal"
)

type StandardReader struct {
	lexer   Lexer[LexemeKind]
	curr    *Token[LexemeKind]
	currErr error
}

type LexemeKind int

const (
	LParen LexemeKind = iota
	RParen
	Float
	Integer
	Name
)

var IndigoLexer = MustNewLexerFactory([]LexerRule[LexemeKind]{
	{`\(`, LParen},
	{`\)`, RParen},
	{`-?\d*\.\d+`, Float},
	{`-?\d+`, Integer},
	{`[^\d\s\(\)\-][^\(\)\s]*`, Name},
}, `\s+`)

func NewStandardReader(source string) interfaces.IndigoReader {

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
	case Float:
		return sr.readFloat()
	case Integer:
		return sr.readInteger()
	case LParen:
		cons, err := sr.readList()
		if err != nil {
			return nil, fmt.Errorf("invalid list: %e", err)
		}
		return cons, nil
	case Name:
		return sr.readSymbol(), nil
	}

	panic("unreachable")
}

func (sr *StandardReader) readSymbol() interfaces.Symbol {
	symbol := interfaces.Symbol(sr.curr.lexeme)
	sr.Next()
	return symbol
}

func (sr *StandardReader) readFloat() (float64, error) {
	float, err := strconv.ParseFloat(sr.curr.lexeme, 64)
	if err != nil {
		return 0, err
	}
	sr.Next()
	return float, nil
}

func (sr *StandardReader) readInteger() (int, error) {
	integer, err := strconv.ParseInt(sr.curr.lexeme, 10, 32)
	if err != nil {
		return 0, err
	}
	sr.Next()
	return int(integer), nil
}

func (sr *StandardReader) readList() (interfaces.List, error) {
	err := sr.Next()
	if err != nil {
		panic(err)
	}

	var elements internal.List

	for sr.curr.kind != RParen {
		value, err := sr.Read()

		if err == EOF {
			return nil, fmt.Errorf("unbalanced parentheses: missing right parenthesis")
		} else if err != nil {
			return nil, err
		}

		elements = append(elements, value)
	}

	return elements, nil
}
