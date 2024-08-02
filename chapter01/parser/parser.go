package parser

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/pkg/errors"
)

var syntaxError = func(msgFormat string, args ...any) error {
	return errors.New(fmt.Sprintf("Syntax error: "+msgFormat, args...))
}

type Parser struct {
	lexer lexer
	file  *token.File
}

func NewParser(source string) *Parser {
	return &Parser{lexer: lexer{input: []byte(source)}, file: token.NewFileSet().AddFile("", 1, len(source))}
}

func (p *Parser) Parse() (*ast.ReturnStmt, error) {
	n, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if b, ok := p.lexer.ReadByte(); ok {
		return nil, syntaxError("unexpected %c", b)
	}
	return n.(*ast.ReturnStmt), nil
}

func (p *Parser) parseStatement() (ast.Node, error) {
	t, start, _, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	switch t {
	case "return":
		return p.parseReturn(p.file.Pos(start))
	}
	return nil, syntaxError("expected a statement got %s", t)
}

func (p *Parser) parseReturn(pos token.Pos) (*ast.ReturnStmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	fmt.Printf("expr: %+v\n", expr)

	// Make sure there is a semicolon after expression
	token, ok := p.lexer.ReadByte()
	if !ok {
		return nil, syntaxError("expected ; after expression")
	}
	if token != ';' {
		return nil, syntaxError("expected ; got %c", token)
	}

	return &ast.ReturnStmt{Return: pos, Results: []ast.Expr{expr}}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	num, start, _, err := p.lexer.ReadNumber()
	if err != nil {
		return nil, syntaxError(err.Error())
	}
	return &ast.BasicLit{ValuePos: p.file.Pos(start), Kind: token.INT, Value: num}, nil
}
