package lexer

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strconv"
)

type Token struct {
	Position token.Position
	GoToken  token.Token
	Literal  string
}

type Lexer struct {
	Scanner scanner.Scanner
	FileSet *token.FileSet

	Tokens   []Token
	Position int

	Error error
}

func (t *Token) String() string {
	return fmt.Sprintf("%s\t%s\t%q\n", t.Position, t.GoToken, t.Literal)
}

func (l *Lexer) Curr() Token {
	if l.Position == len(l.Tokens) {
		pos, tok, lit := l.Scanner.Scan()
		l.Tokens = append(l.Tokens, Token{Position: l.FileSet.Position(pos), GoToken: tok, Literal: lit})
		//debug.Printf("[lexer]: %s", l.Tokens[l.Position])
	}
	return l.Tokens[l.Position]
}

func (l *Lexer) Next() Token {
	l.Position++
	tok := l.Curr()
	return tok
}

func (l *Lexer) Prev() Token {
	if l.Position == 0 {
		panic("no previous token")
	}
	return l.Tokens[l.Position-1]
}

func (l *Lexer) ParseToken(expectedTok token.Token) bool {
	if l.Error != nil {
		return false
	}

	if expectedTok != token.COMMENT {
		for l.Curr().GoToken == token.COMMENT {
			l.Next()
		}
	}

	tok := l.Curr()
	if tok.GoToken == expectedTok {
		l.Next()
		return true
	}

	l.Error = fmt.Errorf("%s:%d:%d: expected %q, got %q (%q)", tok.Position.Filename, tok.Position.Line, tok.Position.Column, expectedTok, tok.GoToken, tok.Literal)
	return false
}

func (l *Lexer) ParseIdent(ident *string) bool {
	if l.ParseToken(token.IDENT) {
		*ident = l.Prev().Literal
		return true
	}
	return false
}

func (l *Lexer) ParseIdentList(idents *[]string) bool {
	var ident string

	for l.ParseIdent(&ident) {
		*idents = append(*idents, ident)
		if !l.ParseToken(token.COMMA) {
			l.Error = nil
			return true
		}
	}

	return len(*idents) != 0
}

func (l *Lexer) ParseIntLit(n *int) bool {
	if l.ParseToken(token.INT) {
		var err error
		*n, err = strconv.Atoi(l.Prev().Literal)
		if err != nil {
			l.Error = fmt.Errorf("failed to parse int value: %v", err)
		}
		return err == nil
	}
	return false
}

func (l *Lexer) ParseStringLit(s *string) bool {
	if l.ParseToken(token.STRING) {
		*s = l.Prev().Literal
		if ((*s)[0] == '"') || ((*s)[0] == '`') {
			*s = (*s)[1:]
		}
		if ((*s)[len(*s)-1] == '"') || ((*s)[len(*s)-1] == '`') {
			*s = (*s)[:len(*s)-1]
		}
		return true
	}
	return false
}
