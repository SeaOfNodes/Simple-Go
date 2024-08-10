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
	disablePeephole := flag.Bool("d", false, "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-a] [-d] <code>\n", os.Args[0])
		fmt.Println("\t-a\tUse Go AST parser")
		fmt.Println("\t-d\tDisable peephole optimizations")
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

	if *useGoAST {
		_, err := simple.GoSimple(code)
		if err != nil {
			log.Fatalf("Compiler error: %v", err)
		}
	} else {
		_, err := simple.Simple(code)
		if err != nil {
			log.Fatalf("Compiler error: %v", err)
		}
	}

	fmt.Printf("Graph:\n\n%s", graph.Visualize())
}
