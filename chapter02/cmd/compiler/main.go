package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	simple "github.com/SeaOfNodes/Simple-Go/chapter02"
	"github.com/SeaOfNodes/Simple-Go/chapter02/graph"
	"github.com/SeaOfNodes/Simple-Go/chapter02/ir"
)

func main() {
	useGoAST := flag.Bool("a", false, "")
	printString := flag.Bool("s", false, "")
	disablePeephole := flag.Bool("d", false, "")
	flag.Usage = func() {
		fmt.Println("Simple compiler written in Go. Prints graph representation of IR.")
		fmt.Printf("Usage: %s [-a] [-d] [-s] <code>\n", os.Args[0])
		fmt.Println("\t-a\tUse Go AST parser")
		fmt.Println("\t-d\tDisable peephole optimizations")
		fmt.Println("\t-s\tPrint string visualization")
		fmt.Println("\t-h\tPrint this help and exit")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Missing code argument")
		flag.Usage()
		return
	}
	code := flag.Args()[0]

	if *disablePeephole {
		ir.DisablePeephole = true
	}

	var node ir.Node
	var err error
	if *useGoAST {
		node, err = simple.GoSimple(code)
		if err != nil {
			log.Fatalf("Compiler error: %v", err)
		}
	} else {
		node, err = simple.Simple(code)
		if err != nil {
			log.Fatalf("Compiler error: %v", err)
		}
	}

	if *printString {
		fmt.Printf("String:\n\n%s", ir.ToString(node))
	} else {
		fmt.Printf("Graph:\n\n%s", graph.Visualize())
	}
}
