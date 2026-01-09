package indigo

import (
	"fmt"
	"testing"
)

func TestLexerPositions(t *testing.T) {
	lexer := IndigoLexer("( 12\n\n 123")
	token, err := lexer.Next()
	if err != nil {
		t.Errorf("could not lex first token")
	}

	t.Log(token.lexeme)

	t.Logf("%#v", lexer.(*RegexLexer[LexemeKind]).position)

	should := DocumentPosition{Line: 0, Char: 0, Offset: 0}
	if token.position != should {
		t.Errorf("expected %v but found %v", should, token.position)
	}

	token, err = lexer.Next()
	if err != nil {
		t.Errorf("could not lex second token: %s", err)
		return
	}

	should = DocumentPosition{Line: 0, Char: 2, Offset: 2}
	if token.position != should {
		t.Errorf("expected %v but found %v", should, token.position)
	}

	token, err = lexer.Next()
	if err != nil {
		t.Errorf("could not lex third token")
	}

	should = DocumentPosition{Line: 2, Char: 1, Offset: 7}
	if token.position != should {
		t.Errorf("expected %v but found %v", should, token.position)
	}
}

func TestLexerKinds(t *testing.T) {
	lexer := IndigoLexer("( 12\n\n 123")
	var lexerKinds []LexemeKind
	for {
		token, err := lexer.Next()
		if err != nil {
			if err == EOF {
				break
			}
			t.Errorf("%s", err)
		}
		lexerKinds = append(lexerKinds, token.kind)
	}

	expectedKinds := []LexemeKind{LParen, Integer, Integer}

	doesNotMatchError := fmt.Errorf("found   : %v\nexpected: %v", lexerKinds, expectedKinds)

	if len(lexerKinds) != len(expectedKinds) {
		t.Error(doesNotMatchError)
	}

	for i := range lexerKinds {
		if lexerKinds[i] != expectedKinds[i] {
			t.Error(doesNotMatchError)
		}
	}

}
