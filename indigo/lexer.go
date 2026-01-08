package indigo

import (
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

var EOF = errors.New("EOF")

type DocumentPosition struct {
	Line   int
	Char   int
	Offset int
}

func (dp *DocumentPosition) advanceRune(r rune) {
	dp.Offset = utf8.RuneLen(r)
	dp.Char += 1
	dp.Line += 1
	if r == '\n' {
		dp.Char = 0
		dp.Line = 0
	}
}

func (dp *DocumentPosition) advanceString(s string) {
	for r := range []rune(s) {
		dp.advanceRune(rune(r))
	}
}

type Token[T any] struct {
	lexeme   string
	position DocumentPosition
	kind     T
}

type Lexer[T any] interface {
	Next() (*Token[T], error)
	Synchronize() error
}

type RegexLexer[T any] struct {
	position DocumentPosition
	rules    map[*regexp.Regexp]T
	skipRule *regexp.Regexp
	source   string
}

func MustNewLexerFactory[T any](rules map[string]T, skipRule string) func(string) Lexer[T] {
	factory, err := NewLexerFactory(rules, skipRule)
	if err != nil {
		panic(err)
	}
	return factory
}

func NewLexerFactory[T any](rules map[string]T, skipRule string) (func(string) Lexer[T], error) {
	regexRules := make(map[*regexp.Regexp]T)
	for regularExpression, kind := range rules {
		compiledRegularExpression, err := regexp.Compile("^" + regularExpression)
		if numSubExpressions := compiledRegularExpression.NumSubexp(); numSubExpressions != 0 {
			return nil, fmt.Errorf("regular expression must have 0 capture groups but %d were found", numSubExpressions)
		}
		if err != nil {
			return nil, err
		}
		regexRules[compiledRegularExpression] = kind
	}

	compiledSkipRule, err := regexp.Compile("^")
	if err != nil {
		return nil, err
	}

	return func(source string) Lexer[T] {
		return &RegexLexer[T]{
			position: DocumentPosition{0, 0, 0},
			skipRule: compiledSkipRule,
			source:   source,
			rules:    regexRules,
		}
	}, nil
}

func (rl *RegexLexer[T]) isEOF() bool {
	return rl.position.Offset >= len(rl.source)
}

func (rl *RegexLexer[T]) currSourceSlice() string {
	return rl.source[rl.position.Offset:]
}

func (rl *RegexLexer[T]) Next() (*Token[T], error) {
	skippableString := rl.skipRule.FindString(rl.currSourceSlice())
	rl.position.advanceString(skippableString)

	if rl.isEOF() {
		return nil, EOF
	}

	for rule, kind := range rl.rules {
		lexeme := rule.FindString(rl.currSourceSlice())

		if lexeme == "" {
			continue
		}

		start_of_lexeme := rl.position
		rl.position.advanceString(lexeme)

		return &Token[T]{
			lexeme:   lexeme,
			position: start_of_lexeme,
			kind:     kind,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (rl *RegexLexer[T]) Synchronize() error {
	skippableStringIndices := rl.skipRule.FindStringIndex(rl.currSourceSlice())

	if skippableStringIndices == nil || skippableStringIndices[0] == 0 {
		return fmt.Errorf("could not locate synchornization point")
	}

	rl.position.advanceString(rl.currSourceSlice()[:skippableStringIndices[0]])

	return nil
}
