package analyzer

import (
	"bufio"
	"dependency-graph-visualiser/src/builder"
	"fmt"
	"os"
	"strings"
)

func LoadTestRepository(path string) (builder.Graph, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open test repository: %w", err)
	}
	defer file.Close()

	if graph, err := buildGraphTest(file); err != nil {
		return nil, err
	} else {
	return graph, nil
	}
}

func fixRepositoryDependencies (graph builder.Graph) (builder.Graph) {
	for _, dependencies := range graph {
		for _, dependency := range dependencies {
			if _, exists := graph[dependency]; !exists {
				graph[dependency] = []string{}
			}
		}
	}
	return graph
}

func buildGraphTest(file *os.File) (builder.Graph, error) {
	graph := make(builder.Graph)
	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		left, right, ok := strings.Cut(line, ":")
		name := strings.TrimSpace(left)
		if !ok || !isTestPackage(name) {
			return nil, fmt.Errorf("line %d: expected format A: B C", lineNumber)
		}
		if _, exists := graph[name]; exists {
			return nil, fmt.Errorf("line %d: package %s is duplicated", lineNumber, name)
		}
		for _, dependency := range strings.Fields(right) {
			if !isTestPackage(dependency) {
				return nil, fmt.Errorf("line %d: invalid package %q", lineNumber, dependency)
			}
			graph[name] = append(graph[name], dependency)
		}
		if graph[name] == nil {
			graph[name] = []string{}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read test repository: %w", err)
	}
	graph = fixRepositoryDependencies(graph)
	return graph, nil
}

func isTestPackage(value string) bool { return len(value) == 1 && value[0] >= 'A' && value[0] <= 'Z' }
