# Simple-Go
**Simple, implemented in Go.**

Please note this is not an exact translation of `Simple`, but it follows the same chapters. It is very much a work in progress!

## How to run
Every chapter (except chapter01) can be run with custom input:
```sh
cd chapter02
go run ./cmd/compiler -s "return 1+1;"
```
The output will be the sea of nodes graph.

For more usage options run:
```sh
go run ./cmd/compiler -h
```

## Compiler instructions
These are instructions that can be embedded into your code. The compiler executes these once they are reached, so their placement matters.

| Instruction | Description |
| ----------- | ----------- |
| #showGraph | Prints the sea of nodes graph in the current state. |
| #disablePeephole | Disables peephole optimizations from the current state. |

*Compiler instructions are only supported from chapter03 onwards.*
