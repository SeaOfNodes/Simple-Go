package ir

import (
	"fmt"
	"strings"
)

func Visualize(generator *Generator) string {
	nodes := allNodes()
	gb := &graphBuilder{}
	gb.StartBlock("digraph chapter04 {")

	// To keep the Scopes below the graph and pointing up into the graph we need to group the Nodes in a subgraph cluster, and the scopes into a different subgraph cluster.  THEN we can draw edges between the scopes and nodes.  If we try to cross subgraph cluster borders while still making the subgraphs DOT gets confused.
	gb.AppendLine("rankdir=BT;") // Force Nodes before Scopes

	// Preserve node input order
	gb.AppendLine("ordering=\"in\";")

	// Merge multiple edges hitting the same node.  Makes common shared nodes much prettier to look at.
	gb.AppendLine("concentrate=\"true\";")

	// Just the Nodes first, in a cluster no edges
	visualizeNodes(gb, nodes)

	visualizeScopes(gb, generator.Scope)

	// Walk the Node edges
	visualizeNodeEdges(gb, nodes)

	visualizeScopeEdges(gb, generator.Scope)

	gb.EndBlock("}")
	return gb.String()
}

func quoteName(n Node) string {
	name := UniqueName(n)
	if name[0] == '$' {
		return "\"" + name + "\""
	}
	return name
}

func visualizeNodes(gb *graphBuilder, nodes []Node) {
	gb.StartBlock("subgraph cluster_Nodes {") // Magic "cluster_" in the subgraph name
	defer gb.EndBlock("}")
	for _, n := range nodes {
		visualizeNode(gb, n)
	}
}

func visualizeNode(gb *graphBuilder, n Node) {
	switch n.(type) {
	case *ScopeNode, *ProjNode:
		return
	}

	gb.Append("%s [ ", quoteName(n))
	defer gb.AppendLine(" ];")

	if m, ok := n.(MultiNode); ok {
		visualizeMultiNode(gb, m)
		return
	}

	if n.IsControl() {
		gb.Append("shape=box style=filled fillcolor=yellow ")
	}
	gb.Append("label=\"%s\"", n.GraphicLabel())
}

func visualizeMultiNode(gb *graphBuilder, m MultiNode) {
	gb.StartBlock("shape=plaintext label=<")
	defer gb.EndBlock(">")
	gb.StartBlock("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\">")
	defer gb.EndBlock("</TABLE>")
	gb.AppendLine("<TR><TD BGCOLOR=\"yellow\">%s</TD></TR>", m.GraphicLabel())

	projTable := false
	for _, o := range Outs(m) {
		if p, ok := o.(*ProjNode); ok {
			if !projTable {
				projTable = true
				gb.StartBlock("<TR><TD><TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"><TR>")
				defer gb.EndBlock("</TR></TABLE></TD></TR>")
			}
			gb.Append("<TD Port=\"p%d\"", p.i)
			if p.IsControl() {
				gb.Append(" BGCOLOR=\"yellow\"")
			}
			gb.AppendLine(">%s</TD>", p.GraphicLabel())
		}
	}
}

func visualizeScopes(gb *graphBuilder, scope *ScopeNode) {
	gb.AppendLine("node [shape=plaintext];") // Magic "cluster_" in the subgraph name
	for level, table := range scope.Scopes {
		scopeName := fmt.Sprintf("%s_%d", UniqueName(scope), level)
		visualizeScope(gb, scopeName, level, table)
	}
	for range scope.Scopes {
		gb.EndBlock("}")
	}
}

func visualizeScope(gb *graphBuilder, scopeName string, level int, table map[string]int) {
	gb.StartBlock("subgraph cluster_%s {", scopeName)
	gb.StartBlock("%s [label=<", scopeName)
	gb.StartBlock("<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">")

	gb.Append("<TR><TD BGCOLOR=\"cyan\">%d</TD>", level)
	for name := range table {
		gb.Append("<TD PORT=\"%s_%s\">%s</TD>", scopeName, name, name)
	}
	gb.AppendLine("</TR>")

	gb.EndBlock("</TABLE>>")
	gb.EndBlock("];")
}

func visualizeNodeEdges(gb *graphBuilder, nodes []Node) {
	gb.AppendLine("edge [ fontname=Helvetica, fontsize=8 ];")
	for _, n := range nodes {
		switch n.(type) {
		case *ScopeNode, *ConstantNode, *ProjNode:
			continue
		}

		// In this chapter we do display the Constant->Start edge;
		for i, def := range Ins(n) {
			if def == nil {
				continue
			}
			// Most edges land here use->def
			gb.Append("%s -> %s", quoteName(n), quoteName(def))
			// Number edges, so we can see how they track
			gb.Append("[taillabel=%d", i)
			if def.IsControl() {
				gb.Append(" color=red")
			}
			// control edges are colored red
			gb.AppendLine("];")
		}
	}
}

func visualizeScopeEdges(gb *graphBuilder, scope *ScopeNode) {
	gb.AppendLine("edge [style=dashed color=cornflowerblue];")
	for level, table := range scope.Scopes {
		scopeName := fmt.Sprintf("%s_%d", UniqueName(scope), level)
		for name, index := range table {
			n := In(scope, index)
			if n == nil {
				continue
			}
			gb.AppendLine("%s:\"%s_%s\"->%s;", scopeName, scopeName, name, quoteName(n))
		}
	}
}

func allNodes() []Node {
	var all []Node
	walkNodes(StartNode, func(n Node) bool {
		all = append(all, n)
		return true
	})
	return all
}

func walkNodes(start Node, walkFunc func(Node) bool) {
	walkNodesInternal(start, walkFunc, make(map[Node]struct{}))
}

func walkNodesInternal(start Node, walkFunc func(Node) bool, walked map[Node]struct{}) {
	if _, ok := walked[start]; ok {
		return
	}
	walked[start] = struct{}{}
	if !walkFunc(start) {
		return
	}

	for _, n := range Outs(start) {
		walkNodesInternal(n, walkFunc, walked)
	}
}

type graphBuilder struct {
	indent  string
	builder strings.Builder
	newLine bool
}

func (g *graphBuilder) Append(format string, args ...any) {
	if g.newLine {
		g.builder.WriteString("\n")
		g.builder.WriteString(g.indent)
		g.newLine = false
	}
	s := fmt.Sprintf(format, args...)
	g.builder.WriteString(s)
}

func (g *graphBuilder) AppendLine(format string, args ...any) {
	g.Append(format, args...)
	g.newLine = true
}

func (g *graphBuilder) StartBlock(format string, args ...any) {
	g.AppendLine(format, args...)
	g.indent += "\t"
}

func (g *graphBuilder) EndBlock(format string, args ...any) {
	g.indent = g.indent[1:]
	g.AppendLine(format, args...)
}

func (g *graphBuilder) String() string {
	return g.builder.String()
}
