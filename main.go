package main

import (
	"flag"
	"fmt"
	"github.com/caos/documentation/pkg/docu"
	"os"
)

func main() {
	var path, struc, md string
	flag.StringVar(&path, "path", "", "The path to the go-file which contains the struct")
	flag.StringVar(&struc, "struct", "", "The name of the struct for which the documentation should be generated")
	flag.StringVar(&md, "output", "", "The path to the folder which should be used for the output")
	flag.Parse()

	if path == "" || struc == "" || md == "" {
		fmt.Println("Please provide all parameters")
	}

	doc := docu.New()
	if err := doc.Parse(path, struc); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := doc.GenerateMarkDown(md, nil); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Finished")
}
