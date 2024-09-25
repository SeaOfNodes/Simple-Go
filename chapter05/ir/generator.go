package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/pkg/errors"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
)

type instruction struct {
	ast.Stmt
	id string
}

// ShowGraphInst instructs the compiler to print the graph state
var ShowGraphInst = &instruction{id: "showGraph"}

// DisablePeepholeInst instructs the compiler to disable peephole optimizations
var DisablePeepholeInst = &instruction{id: "disablePeephole"}

type ASTError struct {
	error
	Pos token.Pos
}

func astError(pos token.Pos, n ast.Node) *ASTError {
	internal := errors.Errorf("Unsupported AST: %#v", n)
	return &ASTError{error: internal, Pos: pos}
}

type Generator struct {
	Scope *ScopeNode
}

func NewGenerator(arg types.Type) *Generator {
	StartNode = newStartNode(types.NewTuple(types.Control, arg))
	return &Generator{Scope: NewScopeNode()}
}

func (g *Generator) Generate(n ast.Node) (*ReturnNode, error) {
	var retNode *ReturnNode
	var err error
	ast.Inspect(n, func(n ast.Node) bool {
		if err != nil {
			return false
		}

		// New scope for the initial control and arguments
		g.Scope.Push()
		defer g.Scope.Pop()
		var control Node
		control, err = peephole(NewProjNode(StartNode, 0, Control))
		if err != nil {
			return false
		}
		g.Scope.Define(Control, control)
		var arg0 Node
		arg0, err = peephole(NewProjNode(StartNode, 1, Arg0))
		if err != nil {
			return false
		}
		g.Scope.Define(Arg0, arg0)

		if block, ok := n.(*ast.BlockStmt); ok {
			var res Node
			res, err = g.generateBlock(block)
			retNode, _ = res.(*ReturnNode)
			return false
		}

		return true
	})
	return retNode, err
}

func (g *Generator) generateBlock(b *ast.BlockStmt) (Node, error) {
	g.Scope.Push()
	defer g.Scope.Pop()

	var res Node
	for _, stmt := range b.List {
		n, err := g.generateStatement(stmt)
		if err != nil {
			return nil, err
		}
		if n != nil {
			res = n
		}
	}

	return res, nil
}

func (g *Generator) generateStatement(s ast.Stmt) (Node, error) {
	switch t := s.(type) {
	case *ast.ReturnStmt:
		return g.generateReturn(t)
	case *ast.DeclStmt:
		spec, ok := t.Decl.(*ast.GenDecl).Specs[0].(*ast.ValueSpec)
		if !ok {
			return nil, astError(s.Pos(), s)
		}
		return g.generateDecl(spec)
	case *ast.BlockStmt:
		return g.generateBlock(t)
	case *ast.AssignStmt:
		return g.generateAssign(t)
	case *ast.IfStmt:
		return g.generateIf(t)
	case *instruction:
		switch s {
		case ShowGraphInst:
			fmt.Println(Visualize(g))
		case DisablePeepholeInst:
			DisablePeephole = true
		}
		return nil, nil
	}
	return nil, astError(s.Pos(), s)
}

func (g *Generator) generateIf(i *ast.IfStmt) (Node, error) {
	expr, err := g.generateExpr(i.Cond)
	if err != nil {
		return nil, err
	}
	switch expr.(type) {
	case *BoolNode, *ConstantNode:
	default:
		return nil, computeError(i.Cond, "if condition must have a bool value")
	}
	ifNode := NewIfNode(g.Scope.Control(), expr)
	pin(ifNode)
	n, err := peephole(ifNode)
	if err != nil {
		return nil, err
	}
	ifNode = n.(*IfNode)
	trueNode, err := peephole(NewProjNode(ifNode, 0, "True"))
	if err != nil {
		return nil, err
	}
	falseNode, err := peephole(NewProjNode(ifNode, 1, "False"))
	if err != nil {
		return nil, err
	}
	falseScope := g.Scope.Clone()

	// Parse true side
	g.Scope.SetControl(trueNode)
	_, err = g.generateStatement(i.Body)
	if err != nil {
		return nil, err
	}
	trueScope := g.Scope

	g.Scope = falseScope
	g.Scope.SetControl(falseNode)
	if i.Else != nil {
		_, err = g.generateStatement(i.Else)
		if err != nil {
			return nil, err
		}
		falseScope = g.Scope
	}

	g.Scope = trueScope
	region, err := trueScope.Merge(falseScope)
	if err != nil {
		return nil, err
	}
	err = g.Scope.SetControl(region)
	if err != nil {
		return nil, err
	}
	return region, nil
}

func (g *Generator) generateAssign(a *ast.AssignStmt) (Node, error) {
	id, ok := a.Lhs[0].(*ast.Ident)
	if !ok {
		return nil, astError(a.Pos(), a)
	}
	expr, err := g.generateExpr(a.Rhs[0])
	if err != nil {
		return nil, err
	}

	exists, err := g.Scope.Update(id.Name, expr)
	if !exists {
		return nil, computeError(id, "unknown identifier")
	}
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (g *Generator) generateDecl(v *ast.ValueSpec) (Node, error) {
	name := v.Names[0].Name
	value, err := g.generateExpr(v.Values[0])
	if err != nil {
		return nil, err
	}
	err = g.Scope.Define(name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (g *Generator) generateReturn(r *ast.ReturnStmt) (*ReturnNode, error) {
	expr, err := g.generateExpr(r.Results[0])
	if err != nil {
		return nil, err
	}
	n, err := peephole(NewReturnNode(g.Scope.Control(), expr))
	if err != nil {
		return nil, err
	}
	g.Scope.SetControl(nil)
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
		case token.EQL:
			return peephole(NewBoolNode(lhs, EQ, rhs))
		case token.GEQ:
			lhs, rhs = rhs, lhs
			fallthrough
		case token.LEQ:
			return peephole(NewBoolNode(lhs, LE, rhs))
		case token.GTR:
			lhs, rhs = rhs, lhs
			fallthrough
		case token.LSS:
			return peephole(NewBoolNode(lhs, LT, rhs))
		case token.NEQ:
			eq, err := peephole(NewBoolNode(lhs, EQ, rhs))
			if err != nil {
				return nil, err
			}
			return peephole(NewNotNode(eq))
		}
	case *ast.ParenExpr:
		return g.generateExpr(t.X)
	case *ast.UnaryExpr:
		value, err := g.generateExpr(t.X)
		if err != nil {
			return nil, err
		}
		switch t.Op {
		case token.SUB:
			return peephole(NewMinusNode(value))
		case token.NOT:
			return peephole(NewNotNode(value))
		}
	case *ast.BasicLit:
		num, err := strconv.Atoi(t.Value)
		if err != nil {
			return nil, err
		}
		return peephole(NewConstantNode(types.NewInt(num)))
	case *ast.Ident:
		n, ok := g.Scope.Lookup(t.Name)
		if !ok {
			return nil, computeError(e, "unknown identifier")
		}
		return n, nil
	}
	return nil, astError(e.Pos(), e)
}
