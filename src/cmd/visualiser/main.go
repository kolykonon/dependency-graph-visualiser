package main

import (
	"dependency-graph-visualiser/src/analyzer"
	"dependency-graph-visualiser/src/params"
	"fmt"
	"os"
)


func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		fmt.Fprintln(os.Stderr, "Incorrect arguments quantity")
		return
	}
	params, err := params.NewParams(
		args[0],
		args[1],
		args[2],
		args[3],
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Printf("name: %s\nurl: %s\nmode: %s\nfilename: %s\n", params.Name, params.URL, params.Mode, params.Filename)

	dep, err := analyzer.FindDependencies(params.Name,params.URL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Dependencies: %s\n", dep)
	
}