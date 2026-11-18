package main

import (
	"dependency-graph-visualiser/src/analyzer"
	"dependency-graph-visualiser/src/builder"
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
	configuration, err := params.NewParams(
		args[0],
		args[1],
		args[2],
		args[3],
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	var graph builder.Graph
	if configuration.Mode == params.WorkModeTesting {
		graph, err = analyzer.LoadTestRepository(configuration.URL)
	} else {
		graph, err = analyzer.LoadProductionRepository(configuration.URL)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	result, err := builder.BuildReachable(graph, configuration.Name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Print(result.DOT())
	if len(result.Cycles) == 0 {
		fmt.Println("Cycles: none")
	}
	for _, cycle := range result.Cycles {
		fmt.Printf("Cycle: %v\n", cycle)
	}
	if err := result.RenderPNG(configuration.Filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("Graph image: %s\n", configuration.Filename)
}
