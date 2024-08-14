package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/pkg/errors"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir/types"
)

type ASTError struct {
	error
	Pos token.Pos
}

func astError(pos token.Pos, n ast.Node) *ASTError {
	internal := errors.New(fmt.Sprintf("Unsupported AST: %v", n))
	return &ASTError{error: internal, Pos: pos}
}

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(n ast.Node) (*ReturnNode, error) {
	var retNode *ReturnNode
	var err error
	ast.Inspect(n, func(n ast.Node) bool {
		if err != nil {
			return false
		}

		if ret, ok := n.(*ast.ReturnStmt); ok {
			retNode, err = g.generateReturn(ret)
			if err != nil {
				return false
			}
		}
		return true
	})
	return retNode, err
}

func (g *Generator) generateReturn(r *ast.ReturnStmt) (*ReturnNode, error) {
	expr, err := g.generateExpr(r.Results[0])
	if err != nil {
		return nil, err
	}
	n, err := peephole(NewReturnNode(StartNode, expr))
	if err != nil {
		return nil, err
	}
	return n.(*ReturnNode), nil
}

func (g *Generator) generateExpr(e ast.Expr) (Node, error) {
	switch t := e.(type) {
	case *ast.BinaryExpr:
		lhs, err := g.generateExpr(t.X)
		if err != nil {
			return nil, err
		}
		rhs, err := g.generateExpr(t.Y)
		if err != nil {
			return nil, err
		}
		switch t.Op {
		case token.ADD:
			return peephole(NewAddNode(lhs, rhs))
		case token.SUB:
			return peephole(NewSubNode(lhs, rhs))
		case token.MUL:
			return peephole(NewMulNode(lhs, rhs))
		case token.QUO:
			return peephole(NewDivNode(lhs, rhs))
		}
	case *ast.UnaryExpr:
		value, err := g.generateExpr(t.X)
		if err != nil {
			return nil, err
		}
		if t.Op == token.SUB {
			return peephole(NewMinusNode(value))
		}
	case *ast.BasicLit:
		num, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, err
		}
		return peephole(NewConstantNode(types.NewIntType(num)))
	}
	return nil, astError(e.Pos(), e)
}
