package graph

import (
	"strconv"
	"strings"

	"github.com/SeaOfNodes/Simple-Go/chapter02/ir"
)

func Visualize() string {
	nodes := allNodes()
	sb := &strings.Builder{}
	sb.WriteString("digraph chapter02 {\n")

	// To keep the Scopes below the graph and pointing up into the graph we need to group the Nodes in a subgraph cluster, and the scopes into a different subgraph cluster.  THEN we can draw edges between the scopes and nodes.  If we try to cross subgraph cluster borders while still making the subgraphs DOT gets confused.
	sb.WriteString("\trankdir=BT;\n") // Force Nodes before Scopes

	// Preserve node input order
	sb.WriteString("\tordering=\"in\";\n")

	// Merge multiple edges hitting the same node.  Makes common shared nodes much prettier to look at.
	sb.WriteString("\tconcentrate=\"true\";\n")

	// Just the Nodes first, in a cluster no edges
	visualizeNodes(sb, nodes)

	// Walk the Node edges
	visualizeNodeEdges(sb, nodes)

	sb.WriteString("}\n")
	return sb.String()
}

func visualizeNodes(sb *strings.Builder, nodes []ir.Node) {
	sb.WriteString("\tsubgraph cluster_Nodes {\n") // Magic "cluster_" in the subgraph name
	for _, n := range nodes {
		sb.WriteString("\t\t")
		sb.WriteString(ir.UniqueName(n))
		sb.WriteString(" [ ")

		// control nodes have box shape, other nodes are ellipses, i.e. default shape
		if n.IsControl() {
			sb.WriteString("shape=box style=filled fillcolor=yellow ")
		}
		sb.WriteString("label=\"")
		sb.WriteString(n.GraphicLabel())
		sb.WriteString("\" ];\n")
	}
	sb.WriteString("\t}\n") // End Node cluster

}

func visualizeNodeEdges(sb *strings.Builder, nodes []ir.Node) {
	sb.WriteString("\tedge [ fontname=Helvetica, fontsize=8 ];\n")
	for _, n := range nodes {
		// In this chapter we do display the Constant->Start edge;
		for i, def := range ir.Ins(n) {
			if def == nil {
				continue
			}
			// Most edges land here use->def
			sb.WriteString("\t")
			sb.WriteString(ir.UniqueName(n))
			sb.WriteString(" -> ")
			sb.WriteString(ir.UniqueName(def))
			// Number edges, so we can see how they track
			sb.WriteString("[taillabel=")
			sb.WriteString(strconv.Itoa(i))
			if _, ok := n.(*ir.ConstantNode); ok || n == ir.StartNode {
				sb.WriteString(" style=dotted")
			} else if def.IsControl() {
				sb.WriteString(" color=red")
			}
			// control edges are colored red
			sb.WriteString("];\n")
		}
	}

}

func allNodes() []ir.Node {
	var all []ir.Node
	walkNodes(ir.StartNode, func(n ir.Node) bool {
		all = append(all, n)
		return true
	})
	return all
}

func walkNodes(start ir.Node, walkFunc func(ir.Node) bool) {
	walkNodesInternal(start, walkFunc, make(map[ir.Node]struct{}))
}

func walkNodesInternal(start ir.Node, walkFunc func(ir.Node) bool, walked map[ir.Node]struct{}) {
	if _, ok := walked[start]; ok {
		return
	}
	walked[start] = struct{}{}
	if !walkFunc(start) {
		return
	}

	for _, n := range ir.Outs(start) {
		walkNodesInternal(n, walkFunc, walked)
	}
}
