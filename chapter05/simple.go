package simple

import (
	goParser "go/parser"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter04/ir"
	"github.com/SeaOfNodes/Simple-Go/chapter04/ir/types"
	"github.com/SeaOfNodes/Simple-Go/chapter04/parser"
)

type SourceError struct {
	internal error
	source   string
	offset   int
}

func (s *SourceError) Error() string {
	msg := "\n" + s.source + "\n"
	msg += strings.Repeat(" ", s.offset) + "^\n"
	msg += s.internal.Error()
	return msg
}

func getArgType(arg any) types.Type {
	switch t := arg.(type) {
	case int:
		return types.NewInt(t)
	default:
		return types.Bottom
	}
}

func Simple(source string, arg any) (*ir.ReturnNode, *ir.Generator, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		// Enrich syntax errors with source info
		if s, ok := err.(*parser.SyntaxError); ok {
			return nil, nil, &SourceError{s, source, s.Offset}
		}
		return nil, nil, err
	}

	generator := ir.NewGenerator(getArgType(arg))
	ret, err := generator.Generate(n)
	if err != nil {
		// Enrich ast errors with source info
		if a, ok := err.(*ir.ASTError); ok {
			return nil, nil, &SourceError{a, source, p.PosToOffset(a.Pos)}
		}
		return nil, nil, err
	}
	return ret, generator, nil
}

func GoSimple(source string, arg any) (*ir.ReturnNode, *ir.Generator, error) {
	n, err := goParser.ParseExpr(source)
	if err != nil {
		return nil, nil, err
	}

	generator := ir.NewGenerator(getArgType(arg))
	retNode, err := generator.Generate(n)
	if err != nil {
		return nil, nil, err
	}
	return retNode, generator, nil
}
