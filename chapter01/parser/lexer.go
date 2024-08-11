package parser

import (
	"fmt"
	"strings"
	"unicode"
)

var EOFError = fmt.Errorf("EOF")

type lexer struct {
	input    []byte
	position int
}

func (l *lexer) isEOF() bool {
	return l.position >= len(l.input)
}

func (l *lexer) peek() (byte, bool) {
	if l.isEOF() {
		return 0, false
	}
	return l.input[l.position], true
}

func (l *lexer) nextByte() (byte, bool) {
	c, ok := l.peek()
	l.position++
	return c, ok
}

func (l *lexer) isWhitespace() bool {
	c, ok := l.peek()
	if !ok {
		return false
	}
	// Includes space, tab, newline, CR
	return c <= ' '
}

// skipWhitespace forwards the current position until the next non-whitespace character
func (l *lexer) skipWhitespace() {
	for l.isWhitespace() {
		l.position++
	}
}

// isValidIDStart returns true when b can be the start of an identifier
func isValidIDStart(b byte) bool {
	return unicode.IsLetter(rune(b)) || b == '_'
}

// isValidIDByte returns true when b is a valid identifier character
func isValidIDByte(b byte) bool {
	return unicode.In(rune(b), unicode.Letter, unicode.Digit) || b == '_'
}

func (l *lexer) parseToken(isInToken func(b byte) bool) (string, int) {
	start := l.position
	for {
		b, ok := l.nextByte()
		if !ok || !isInToken(b) {
			break
		}
	}
	// Read one byte after last byte in token, so return to the end of token
	l.position--
	s := string(l.input[start:l.position])
	return s, start
}

func (l *lexer) parseID() (string, int) {
	return l.parseToken(isValidIDByte)
}

func (l *lexer) parsePunctuation() (string, int) {
	return l.parseToken(func(b byte) bool { return unicode.IsPunct(rune(b)) })
}

func (l *lexer) parseNumberString() (string, int, error) {
	s, pos := l.parseToken(func(b byte) bool { return unicode.IsDigit(rune(b)) })
	if len(s) > 1 && s[0] == '0' {
		return "", pos, fmt.Errorf("integer values cannot start with '0'")
	}
	if len(s) == 0 {
		return "", pos, fmt.Errorf("not a number")
	}
	return s, pos, nil
}

// ReadNumber skips whitespaces and retrieves the next number from input. If the next token is not a valid number, an error is returned.
func (l *lexer) ReadNumber() (string, int, error) {
	l.skipWhitespace()
	return l.parseNumberString()
}

// ReadToken skips whitespaces and retrieves the next token from input.
func (l *lexer) ReadToken() (string, int, error) {
	l.skipWhitespace()
	b, ok := l.peek()
	switch {
	case !ok:
		return "", l.position, nil
	case isValidIDStart(b):
		id, pos := l.parseID()
		return id, pos, nil
	case unicode.IsNumber(rune(b)):
		return l.parseNumberString()
	case unicode.IsPunct(rune(b)):
		p, pos := l.parsePunctuation()
		return p, pos, nil
	}
	l.position++
	return string(b), l.position - 1, nil
}

// ReadOp skips whitespaces and retrieves the next op (+-/*) from input.
func (l *lexer) ReadOp() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || !strings.ContainsRune("+-*/", rune(b)) {
		return 0, 0, false
	}

	l.position++
	return b, l.position - 1, true
}

// ReadByte retreives the next non-whitespace byte from input. Returns false if there are no non-whitespace bytes in input.
func (l *lexer) ReadByte() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.nextByte()
	if !ok {
		return 0, 0, false
	}
	return b, l.position - 1, true
}
