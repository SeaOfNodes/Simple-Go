package simple

import (
	goParser "go/parser"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir"
	"github.com/SeaOfNodes/Simple-Go/chapter02/parser"
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

func Simple(source string) (*ir.ReturnNode, error) {
	p := parser.NewParser(source)
	n, err := p.Parse()
	if err != nil {
		// Enrich syntax errors with source info
		if s, ok := err.(*parser.SyntaxError); ok {
			return nil, &SourceError{s, source, s.Offset}
		}
		return nil, err
	}

	generator := ir.NewGenerator()
	ret, err := generator.Generate(n)
	if err != nil {
		// Enrich ast errors with source info
		if a, ok := err.(*ir.ASTError); ok {
			return nil, &SourceError{a, source, p.PosToOffset(a.Pos)}
		}
		return nil, err
	}
	return ret, nil
}

func GoSimple(source string) (*ir.ReturnNode, error) {
	n, err := goParser.ParseExpr(source)
	if err != nil {
		return nil, err
	}

	generator := ir.NewGenerator()
	return generator.Generate(n)
}
