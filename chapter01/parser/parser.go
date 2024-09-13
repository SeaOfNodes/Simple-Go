package parser

import (
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type SyntaxError struct {
	error
	Offset int
}

func syntaxError(offset int, msgFormat string, args ...any) *SyntaxError {
	internal := errors.Errorf("Syntax error: "+msgFormat, args...)
	return &SyntaxError{error: internal, Offset: offset}
}

type Parser struct {
	lexer  lexer
	file   *token.File
	fset   *token.FileSet
	source string
}

func NewParser(source string) *Parser {
	fset := token.NewFileSet()
	return &Parser{source: source, lexer: lexer{input: []byte(source)}, fset: fset, file: fset.AddFile("", 1, len(source))}
}

func (p *Parser) Parse() (*ast.ReturnStmt, error) {
	n, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if b, offset, ok := p.lexer.ReadByte(); ok {
		return nil, syntaxError(offset, "unexpected %c", b)
	}
	return n.(*ast.ReturnStmt), nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	t, offset, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	switch t {
	case "return":
		return p.parseReturn(p.offsetToPos(offset))
	}
	return nil, syntaxError(offset, "expected a statement got %s", t)
}

func (p *Parser) parseReturn(pos token.Pos) (*ast.ReturnStmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Make sure there is a semicolon after expression
	token, offset, ok := p.lexer.ReadByte()
	if !ok {
		return nil, syntaxError(offset, "expected ; after expression")
	}
	if token != ';' {
		return nil, syntaxError(offset, "expected ; got %c", token)
	}

	return &ast.ReturnStmt{Return: pos, Results: []ast.Expr{expr}}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parsePrimary()
}

func opToToken(op byte) token.Token {
	switch op {
	case '+':
		return token.ADD
	case '-':
		return token.SUB
	case '*':
		return token.MUL
	case '/':
		return token.QUO
	}
	return token.ILLEGAL
}

func (p *Parser) offsetToPos(offset int) token.Pos {
	return p.file.Pos(offset)
}

func (p *Parser) PosToOffset(pos token.Pos) int {
	return p.file.Offset(pos)
}

func (p *Parser) parsePrimary() (*ast.BasicLit, error) {
	num, offset, err := p.lexer.ReadNumber()
	if err != nil {
		return nil, syntaxError(offset, err.Error())
	}
	return &ast.BasicLit{ValuePos: p.offsetToPos(offset), Kind: token.INT, Value: num}, nil
}

// string creates a string representation of the node n. Used for debugging.
func (p *Parser) string(n ast.Node) string {
	sb := &strings.Builder{}
	printer.Fprint(sb, p.fset, n)
	return sb.String()
}
