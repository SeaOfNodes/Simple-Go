package parser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

var EOFError = fmt.Errorf("EOF")
var NANError = errors.New("not a number")

type lexer struct {
	input    []byte
	position int
}

func (l *lexer) IsEOF() bool {
	return l.position >= len(l.input)
}

func (l *lexer) peek() (byte, bool) {
	if l.IsEOF() {
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

func (l *lexer) parseNumberString() (string, int, error) {
	s, pos := l.parseToken(func(b byte) bool { return unicode.IsDigit(rune(b)) })
	if len(s) > 1 && s[0] == '0' {
		return "", pos, fmt.Errorf("integer values cannot start with '0'")
	}
	if len(s) == 0 {
		return "", pos, NANError
	}
	return s, pos, nil
}

// ReadNumber skips whitespaces and retrieves the next number from input. If the next token is not a valid number, an error is returned.
func (l *lexer) ReadNumber() (string, int, error) {
	l.skipWhitespace()
	return l.parseNumberString()
}

// ReadToken skips whitespaces and retrieves the next token from input. Returns the token, the offset of the start of the token, true if the token is a valid identifier and an error if one occurred.
func (l *lexer) ReadToken() (string, int, bool, error) {
	l.skipWhitespace()
	b, ok := l.peek()
	switch {
	case !ok:
		return "", l.position, false, nil
	case isValidIDStart(b):
		id, pos := l.parseID()
		return id, pos, true, nil
	case unicode.IsNumber(rune(b)):
		num, pos, err := l.parseNumberString()
		return num, pos, false, err
	}
	l.position++
	return string(b), l.position - 1, false, nil
}

// ReadOp skips whitespaces and retrieves the next op (+-/*) from input.
func (l *lexer) ReadOp() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || !strings.ContainsRune("+-*/=", rune(b)) {
		return 0, l.position, false
	}

	l.position++
	return b, l.position - 1, true
}

// ReadByte retreives the next non-whitespace byte from input. Returns false if there are no non-whitespace bytes in input.
func (l *lexer) ReadByte() (byte, int, bool) {
	l.skipWhitespace()
	b, ok := l.nextByte()
	if !ok {
		return 0, l.position, false
	}
	return b, l.position - 1, true
}

func (l *lexer) ReadID() (string, int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || !isValidIDStart(b) {
		return "", l.position, false
	}
	id, pos := l.parseID()
	return id, pos, true
}

func (l *lexer) Read(next byte) (int, bool) {
	l.skipWhitespace()
	b, ok := l.peek()
	if !ok || b != next {
		return l.position, false
	}
	l.position++
	return l.position - 1, true
}
