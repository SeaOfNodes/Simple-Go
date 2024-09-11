package ast

import (
	"go/ast"
	"go/token"
	"strconv"
)

func Num(i int) *ast.BasicLit {
	return &ast.BasicLit{Value: strconv.Itoa(i), Kind: token.INT}
}

func Expr(a any) ast.Expr {
	switch t := a.(type) {
	case int:
		return Num(t)
	case string:
		return ID(t)
	case ast.Expr:
		return t
	}
	return nil
}

func Bin(lhs any, op string, rhs any) *ast.BinaryExpr {
	return &ast.BinaryExpr{X: Expr(lhs), Op: Op(op), Y: Expr(rhs)}
}

func Un(op string, value any) *ast.UnaryExpr {
	return &ast.UnaryExpr{X: Expr(value), Op: Op(op)}
}

func Paren(value any) *ast.ParenExpr {
	return &ast.ParenExpr{X: Expr(value)}
}

func Assign(id string, value any) *ast.AssignStmt {
	return &ast.AssignStmt{Lhs: []ast.Expr{ID(id)}, Rhs: []ast.Expr{Expr(value)}}
}

func Decl(id string, value any) *ast.DeclStmt {
	return &ast.DeclStmt{Decl: &ast.GenDecl{
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names:  []*ast.Ident{ID(id)},
				Values: []ast.Expr{Expr(value)},
			},
		},
	}}
}

func Op(op string) token.Token {
	switch op {
	case "+":
		return token.ADD
	case "-":
		return token.SUB
	case "*":
		return token.MUL
	case "/":
		return token.QUO
	}
	return 0
}

func ID(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func Ret(expr any) *ast.ReturnStmt {
	return &ast.ReturnStmt{
		Results: []ast.Expr{
			Expr(expr),
		},
	}
}

func Block(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}
