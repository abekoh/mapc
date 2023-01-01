package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/abekoh/mapstructor/internal/generator"
	"log"
	"strings"
)

type parameters struct {
	from generator.Param
	to   generator.Param
	out  string
}

func main() {
	var params parameters
	flag.StringVar(&params.out, "out", "", "output file")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("number of arguments must be 2")
	}
	params.from = parseStructParam(args[0])
	params.to = parseStructParam(args[1])
	run(params)
}

func parseStructParam(s string) generator.Param {
	ss := strings.Split(s, ":")
	if len(ss) != 3 {
		log.Fatal("from/to param must be like 'filepath:package:structname'")
	}
	return generator.Param{
		Dir:    ss[0],
		Pkg:    ss[1],
		Struct: ss[2],
	}
}

func run(params parameters) {
	g := generator.Generator{}
	var buf bytes.Buffer
	if err := g.Generate(&buf, params.from, params.to); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", buf.String())
}
