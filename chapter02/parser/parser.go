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
	return p.parseBinary()
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

// parseUnary parses the next unary operation(s). This is a recursive functiont that stops once a none-unary operation is met.
func (p *Parser) parseUnary() (ast.Expr, error) {
	op, opPos, hasOp := p.parseOp()
	if !hasOp {
		return p.parsePrimary()
	}

	value, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	return &ast.UnaryExpr{X: value, Op: op, OpPos: opPos}, nil
}

// parseBinary parses unary and binary operations recursively. Stops once there are no more operations to parse, or an illegal sequence is encountered.
func (p *Parser) parseBinary() (ast.Expr, error) {
	lhs, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	return p.parseRhs(lhs)
}

// parseRhs recursively parses the next operation and right hand side, with the given left hand side. Stops once there are no more operations to parse.
func (p *Parser) parseRhs(lhs ast.Expr) (ast.Expr, error) {
	op, opPos, hasOp := p.parseOp()
	if !hasOp {
		return lhs, nil
	}

	rhs, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	return p.parseRhs(p.withPrecedence(lhs, op, opPos, rhs))
}

func (p *Parser) parseOp() (token.Token, token.Pos, bool) {
	op, opOffset, ok := p.lexer.ReadOp()
	if !ok {
		return 0, 0, false
	}
	return opToToken(op), p.offsetToPos(opOffset), true
}

// withPrecedence returns the given lhs, op and rhs as a binary expression while taking into consideration operation precedence. This function assumes rhs is a unary/primary expression.
func (p *Parser) withPrecedence(lhs ast.Expr, op token.Token, opPos token.Pos, rhs ast.Expr) (b *ast.BinaryExpr) {
	binLhs, ok := lhs.(*ast.BinaryExpr)
	if !ok {
		return &ast.BinaryExpr{X: lhs, Y: rhs, Op: op, OpPos: opPos}
	}

	// lExpr lOp mExpr rOp rExpr
	lExpr, mExpr, rExpr := binLhs.X, binLhs.Y, rhs
	lOp, rOp := binLhs.Op, op
	lOpPos, rOpPos := binLhs.OpPos, opPos

	if lOp.Precedence() >= rOp.Precedence() {
		lhs := &ast.BinaryExpr{X: lExpr, Y: mExpr, Op: lOp, OpPos: lOpPos}
		return &ast.BinaryExpr{X: lhs, Y: rExpr, Op: rOp, OpPos: rOpPos}
	}

	rhs = &ast.BinaryExpr{X: mExpr, Y: rExpr, Op: rOp, OpPos: rOpPos}
	return &ast.BinaryExpr{X: lExpr, Y: rhs, Op: lOp, OpPos: lOpPos}
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
