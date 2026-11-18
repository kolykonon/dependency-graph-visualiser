package builder

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type Graph map[string][]string

type Result struct {
	Graph  Graph
	Cycles [][]string
}

func BuildReachable(source Graph, root string) (Result, error) {
	if _, ok := source[root]; !ok {
		return Result{}, fmt.Errorf("package %q is unreachable from the repository description", root)
	}
	result := Result{Graph: make(Graph)}
	colour := make(map[string]uint8) // 0 white, 1 gray, 2 black
	type frame struct {
		name string
		next int
	}
	stack := []frame{{name: root}}
	path := []string{root}
	colour[root] = 1
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		deps := source[top.name]
		if top.next == 0 {
			result.Graph[top.name] = append([]string(nil), deps...)
		}
		if top.next == len(deps) {
			colour[top.name] = 2
			stack = stack[:len(stack)-1]
			path = path[:len(path)-1]
			continue
		}
		dependency := deps[top.next]
		top.next++
		switch colour[dependency] {
		case 0:
			colour[dependency] = 1
			stack = append(stack, frame{name: dependency})
			path = append(path, dependency)
		case 1:
			for i, packageName := range path {
				if packageName == dependency {
					cycle := append([]string(nil), path[i:]...)
					result.Cycles = append(result.Cycles, append(cycle, dependency))
					break
				}
			}
		}
	}
	return result, nil
}

func (r Result) DOT() string {
	lines := []string{"digraph dependencies {", "  rankdir=LR;", "  node [shape=box];"}
	nodes := make([]string, 0, len(r.Graph))
	for node := range r.Graph {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	for _, node := range nodes {
		if len(r.Graph[node]) == 0 {
			lines = append(lines, fmt.Sprintf("  %q;", node))
		}
		for _, dependency := range r.Graph[node] {
			lines = append(lines, fmt.Sprintf("  %q -> %q;", node, dependency))
		}
	}
	return strings.Join(append(lines, "}"), "\n") + "\n"
}

func (r Result) RenderPNG(filename string) error {
	command := exec.Command("dot", "-Tpng", "-o", filename)
	command.Stdin = bytes.NewBufferString(r.DOT())
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("render graph with Graphviz: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
