package analyzer

import (
	"bufio"
	"dependency-graph-visualiser/src/builder"
	"fmt"
	"io"
	"strings"
)

func LoadProductionRepository(url string) (builder.Graph, error) {
	path, err := downloadPackageArchive(url, "data")
	if err != nil {
		return nil, err
	}
	reader, err := openPackageArchive(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	graph := make(builder.Graph)
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var name, depends string
	store := func() {
		if name != "" {
			graph[name] = parseDependencies(depends)
		}
		name, depends = "", ""
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			store()
			continue
		}
		if value, ok := strings.CutPrefix(line, "Package: "); ok {
			name = value
		} else if value, ok := strings.CutPrefix(line, "Depends: "); ok {
			depends = value
		} else if strings.HasPrefix(line, " ") && depends != "" {
			depends += " " + strings.TrimSpace(line)
		}
	}
	store()
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read repository: %w", err)
	}
	return graph, nil
}

// parseDependencies chooses the first package from each alternative group and
// strips version and architecture qualifiers used in Debian metadata.
func parseDependencies(value string) []string {
	var result []string
	for _, group := range strings.Split(value, ",") {
		name := strings.TrimSpace(strings.Split(group, "|")[0])
		fields := strings.Fields(name)
		if len(fields) == 0 {
			continue
		}
		name = fields[0]
		name = strings.Split(name, ":")[0]
		if name != "" {
			result = append(result, name)
		}
	}
	return result
}

// find package in opened archive
func findPackage(reader io.Reader, packageName string) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	var currentName string
	var depends string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if currentName == packageName {
				return depends, nil
			}
			currentName, depends = "", ""
			continue
		}

		if value, ok := strings.CutPrefix(line, "Package: "); ok {
			currentName = value
		} else if value, ok := strings.CutPrefix(line, "Depends: "); ok {
			depends = value
		} else if strings.HasPrefix(line, " ") && depends != "" {
			depends += " " + strings.TrimSpace(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if currentName == packageName {
		return depends, nil
	}

	return "", fmt.Errorf("package %q not found", packageName)
}

// full dep extraction pipeline
func FindDependencies(packageName string, url string) (string, error) {
	path, err := downloadPackageArchive(url, "data")
	if err != nil {
		return "", err
	}
	reader, err := openPackageArchive(path)
	if err != nil {
		return "", err
	}
	dependency, err := findPackage(reader, packageName)
	if err != nil {
		return "", err
	}

	return dependency, nil

}
