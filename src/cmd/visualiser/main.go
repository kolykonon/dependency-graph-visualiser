// Command visualiser builds and renders a package dependency graph.
package main

import (
	"dependency-graph-visualiser/src/analyzer"
	"dependency-graph-visualiser/src/builder"
	"dependency-graph-visualiser/src/params"
	"fmt"
	"os"
	"path/filepath"
)

const argumentCount = 4

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(args []string) error {
	if len(args) != argumentCount {
		return fmt.Errorf("usage: visualiser <package> <repository> <mode> <output.svg>")
	}
	configuration, err := params.NewParams(args[0], args[1], args[2], args[3])
	if err != nil {
		return err
	}
	result, err := buildGraph(configuration)
	if err != nil {
		return err
	}
	return writeOutput(result, configuration.Filename)
}

func buildGraph(configuration *params.Params) (builder.Result, error) {
	graph, err := loadGraph(configuration)
	if err != nil {
		return builder.Result{}, err
	}
	return builder.BuildReachable(graph, configuration.Name)
}

func loadGraph(configuration *params.Params) (builder.Graph, error) {
	if configuration.Mode == params.WorkModeTesting {
		return analyzer.LoadTestRepository(configuration.URL)
	}
	return analyzer.LoadProductionRepository(configuration.URL)
}

func writeOutput(result builder.Result, svgFilename string) error {
	plantUMLFilename := plantUMLFilename(svgFilename)
	if err := result.WritePlantUML(plantUMLFilename); err != nil {
		return err
	}
	if err := result.RenderSVG(svgFilename); err != nil {
		return err
	}
	fmt.Print(result.PlantUML())
	printCycles(result.Cycles)
	fmt.Printf("PlantUML: %s\nGraph image: %s\n", plantUMLFilename, svgFilename)
	return nil
}

func plantUMLFilename(svgFilename string) string {
	baseName := svgFilename[:len(svgFilename)-len(filepath.Ext(svgFilename))]
	return baseName + ".puml"
}

func printCycles(cycles [][]string) {
	if len(cycles) == 0 {
		fmt.Println("Cycles: none")
		return
	}
	for _, cycle := range cycles {
		fmt.Printf("Cycle: %v\n", cycle)
	}
}
