package parser

import (
	"fmt"

	"github.com/SeaOfNodes/Simple-Go/chapter01/node"
	"github.com/pkg/errors"
)

var syntaxError = func(msgFormat string, args ...any) error {
	return errors.New(fmt.Sprintf("Syntax error: "+msgFormat, args...))
}

type parser struct {
	lexer lexer
}

func NewParser(source string) *parser {
	return &parser{lexer: lexer{input: []byte(source)}}
}

func (p *parser) Parse() (*node.ReturnNode, error) {
	n, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	if b, ok := p.lexer.ReadByte(); ok {
		return nil, syntaxError("unexpected %c", b)
	}
	return n.(*node.ReturnNode), nil
}

func (p *parser) parseStatement() (node.Node, error) {
	token, err := p.lexer.ReadToken()
	if err != nil {
		return nil, err
	}

	switch token {
	case "return":
		return p.parseReturn()
	}
	return nil, syntaxError("expected a statement got %s", token)
}

func (p *parser) parseReturn() (*node.ReturnNode, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	// Make sure there is a semicolon after expression
	token, ok := p.lexer.ReadByte()
	if !ok {
		return nil, syntaxError("expected ; after expression")
	}
	if token != ';' {
		return nil, syntaxError("expected ; got %c", token)
	}

	return node.NewReturnNode(node.StartNode, expr), nil
}

func (p *parser) parseExpr() (node.Node, error) {
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (node.Node, error) {
	num, err := p.lexer.ReadNumber()
	if err != nil {
		return nil, syntaxError(err.Error())
	}
	return node.NewConstantNode(num), nil
}
