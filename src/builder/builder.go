// Package builder constructs dependency graphs and their representations.
package builder

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	whiteVertex = iota
	grayVertex
	blackVertex
	svgBaseWidth    = 220
	svgNodeSpacing  = 150
	svgNodeCenter   = 110
	svgNodeTop      = 80
	svgNodeWidth    = 80
	svgNodeHeight   = 34
	svgNodeTextY    = 102
	svgEdgeStartY   = 95
	svgEdgeControlY = 35
)

// Graph maps each package to its direct dependencies.
type Graph map[string][]string

// Result contains the reachable graph and its directed cycles.
type Result struct {
	Graph  Graph
	Cycles [][]string
}

type frame struct {
	name string
	next int
}

// BuildReachable builds the transitive graph from root using iterative DFS.
func BuildReachable(source Graph, root string) (Result, error) {
	if _, ok := source[root]; !ok {
		return Result{}, fmt.Errorf("package %q is unreachable from the repository", root)
	}
	result := Result{Graph: make(Graph)}
	colour := map[string]uint8{root: grayVertex}
	stack := []frame{{name: root}}
	path := []string{root}
	for len(stack) > 0 {
		stack, path = processTop(source, &result, colour, stack, path)
	}
	return result, nil
}

func processTop(
	source Graph,
	result *Result,
	colour map[string]uint8,
	stack []frame,
	path []string,
) ([]frame, []string) {
	top := &stack[len(stack)-1]
	dependencies := source[top.name]
	if top.next == whiteVertex {
		result.Graph[top.name] = append([]string(nil), dependencies...)
	}
	if top.next == len(dependencies) {
		colour[top.name] = blackVertex
		return stack[:len(stack)-1], path[:len(path)-1]
	}
	dependency := dependencies[top.next]
	top.next++
	return visitDependency(result, colour, stack, path, dependency)
}

func visitDependency(
	result *Result,
	colour map[string]uint8,
	stack []frame,
	path []string,
	dependency string,
) ([]frame, []string) {
	switch colour[dependency] {
	case whiteVertex:
		colour[dependency] = grayVertex
		return append(stack, frame{name: dependency}), append(path, dependency)
	case grayVertex:
		result.Cycles = appendCycle(result.Cycles, path, dependency)
	}
	return stack, path
}

func appendCycle(cycles [][]string, path []string, dependency string) [][]string {
	for index, packageName := range path {
		if packageName == dependency {
			cycle := append([]string(nil), path[index:]...)
			return append(cycles, append(cycle, dependency))
		}
	}
	return cycles
}

// PlantUML returns a PlantUML diagram for the reachable graph.
func (r Result) PlantUML() string {
	nodes := sortedNodes(r.Graph)
	aliases := make(map[string]string, len(nodes))
	lines := []string{"@startuml", "left to right direction", "skinparam shadowing false"}
	for index, node := range nodes {
		aliases[node] = fmt.Sprintf("P%d", index)
		lines = append(lines, fmt.Sprintf("rectangle %q as %s", node, aliases[node]))
	}
	for _, node := range nodes {
		for _, dependency := range r.Graph[node] {
			edge := fmt.Sprintf("%s --> %s", aliases[node], aliases[dependency])
			lines = append(lines, edge)
		}
	}
	return strings.Join(append(lines, "@enduml"), "\n") + "\n"
}

// WritePlantUML saves the PlantUML source representation.
func (r Result) WritePlantUML(filename string) error {
	if err := os.WriteFile(filename, []byte(r.PlantUML()), 0644); err != nil {
		return fmt.Errorf("write PlantUML file: %w", err)
	}
	return nil
}

// RenderSVG writes a standalone SVG without a local PlantUML installation.
func (r Result) RenderSVG(filename string) error {
	nodes := sortedNodes(r.Graph)
	width := svgBaseWidth + len(nodes)*svgNodeSpacing
	var svg strings.Builder
	writeSVGHeader(&svg, width)
	position := svgPositions(nodes)
	writeSVGEdges(&svg, r.Graph, nodes, position)
	writeSVGNodes(&svg, nodes, position)
	svg.WriteString("</svg>\n")
	if err := os.WriteFile(filename, []byte(svg.String()), 0644); err != nil {
		return fmt.Errorf("write SVG file: %w", err)
	}
	return nil
}

func sortedNodes(graph Graph) []string {
	nodes := make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	return nodes
}

func writeSVGHeader(svg *strings.Builder, width int) {
	fmt.Fprintf(svg, "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" ", width, svgBaseWidth)
	fmt.Fprintf(svg, "viewBox=\"0 0 %d %d\">\n", width, svgBaseWidth)
	svg.WriteString("<defs><marker id=\"arrow\" markerWidth=\"8\" markerHeight=\"8\" ")
	svg.WriteString("refX=\"7\" refY=\"4\" orient=\"auto\">")
	svg.WriteString("<path d=\"M0,0 L8,4 L0,8 z\" fill=\"#356A9A\"/></marker></defs>\n")
	svg.WriteString("<rect width=\"100%\" height=\"100%\" fill=\"white\"/>\n")
}

func svgPositions(nodes []string) map[string]int {
	position := make(map[string]int, len(nodes))
	for index, node := range nodes {
		position[node] = svgNodeCenter + index*svgNodeSpacing
	}
	return position
}

func writeSVGEdges(svg *strings.Builder, graph Graph, nodes []string, position map[string]int) {
	for _, node := range nodes {
		for _, dependency := range graph[node] {
			controlX := (position[node] + position[dependency]) / 2
			fmt.Fprintf(
				svg,
				"<path d=\"M%d %d Q%d %d %d %d\" fill=\"none\" ",
				position[node], svgEdgeStartY, controlX, svgEdgeControlY,
				position[dependency], svgEdgeStartY,
			)
			svg.WriteString("stroke=\"#356A9A\" stroke-width=\"2\" marker-end=\"url(#arrow)\"/>\n")
		}
	}
}

func writeSVGNodes(svg *strings.Builder, nodes []string, position map[string]int) {
	for _, node := range nodes {
		fmt.Fprintf(
			svg,
			"<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" ",
			position[node]-svgNodeWidth/2,
			svgNodeTop,
			svgNodeWidth,
			svgNodeHeight,
		)
		svg.WriteString("rx=\"6\" fill=\"#F7FBFF\" stroke=\"#356A9A\"/>")
		fmt.Fprintf(svg, "<text x=\"%d\" y=\"%d\" text-anchor=\"middle\" ", position[node], svgNodeTextY)
		fmt.Fprintf(svg, "font-family=\"sans-serif\" font-size=\"14\">%s</text>\n", xmlEscape(node))
	}
}

func xmlEscape(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;").Replace(value)
}
