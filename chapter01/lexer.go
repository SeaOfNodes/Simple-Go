package parser

import (
	"fmt"
	"strconv"
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

func (l *lexer) parseToken(isInToken func(b byte) bool) string {
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
	return s
}

func (l *lexer) parseID() string {
	return l.parseToken(isValidIDByte)
}

func (l *lexer) parsePunctuation() string {
	return l.parseToken(func(b byte) bool { return unicode.IsPunct(rune(b)) })
}

func (l *lexer) parseNumberString() (string, error) {
	s := l.parseToken(func(b byte) bool { return unicode.IsDigit(rune(b)) })
	if len(s) > 1 && s[0] == '0' {
		return "", fmt.Errorf("integer values cannot start with '0'")
	}
	return s, nil
}

func (l *lexer) parseNumber() (int, error) {
	s, err := l.parseNumberString()
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number: %s", s)
	}
	return i, nil
}

// ReadNumber skips whitespaces and retrieves the next number from input. If the next token is not a valid number, an error is returned.
func (l *lexer) ReadNumber() (int, error) {
	l.skipWhitespace()
	return l.parseNumber()
}

// ReadToken skips whitespaces and retrieves the next token from input.
func (l *lexer) ReadToken() (string, error) {
	l.skipWhitespace()
	b, ok := l.peek()
	switch {
	case !ok:
		return "", nil
	case isValidIDStart(b):
		return l.parseID(), nil
	case unicode.IsNumber(rune(b)):
		return l.parseNumberString()
	case unicode.IsPunct(rune(b)):
		return l.parsePunctuation(), nil

	}
	return string(b), nil
}

// ReadByte retreives the next non-whitespace byte from input. Returns false if there are no non-whitespace bytes in input.
func (l *lexer) ReadByte() (byte, bool) {
	l.skipWhitespace()
	return l.nextByte()
}

// TODO: Remove?
func (l *lexer) isMatch(s string) bool {
	l.skipWhitespace()
	sLen := len(s)
	if l.position+sLen >= len(l.input) {
		return false
	}
	if s != string(l.input[l.position:l.position+sLen]) {
		return false
	}
	l.position += sLen
	return true
}
