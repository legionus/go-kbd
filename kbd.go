package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/legion/go-kbd/lexer"
)

func main() {
	flag.Parse()

	input := os.Stdin
	if len(flag.Args()) == 1 {
		s, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer input.Close()
		input = s
	} else if len(flag.Args()) != 0 {
		log.Fatalf("Usage: %s [<filename>]", os.Args[0])
	}

	lex := lexer.NewLexer(input)

	nodes := []lexer.Node{}
	for {
		n, err := lex.Get()
		if err != nil {
			log.Fatal(err)
		}
		if n == nil {
			break
		}
		nodes = append(nodes, n)

		b, err := json.MarshalIndent(n, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(b))
	}
/*
	b, err := json.MarshalIndent(nodes, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(b))
*/
}
