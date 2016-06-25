package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/alexflint/go-mold"
	"github.com/kr/pretty"
)

func main() {
	var args struct {
		File string `arg:"positional,required"`
		Type string `arg:"positional"`
	}
	arg.MustParse(&args)

	f, err := os.Open(args.File)
	if err != nil {
		log.Fatal(err)
	}

	types, err := mold.LoadTypes(f)
	if err != nil {
		log.Fatal(err)
	}

	if args.Type == "" {
		for name, t := range types {
			fmt.Printf("%20s := %-30v %s\n", name, t, t.Kind())
		}
	} else {
		pretty.Println(types[args.Type])
	}
}
