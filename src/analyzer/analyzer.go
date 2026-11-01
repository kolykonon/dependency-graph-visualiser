package analyzer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// find package in opened archive
func findPackage(reader io.Reader, packageName string) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 4 * 1024 * 1024)
	
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
		} else if value, ok := strings.CutPrefix(line,"Depends: "); ok {
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
func FindDependencies(packageName string, url string) (string, error){
	path, err := downloadPackageArchive(url, "data")
	if err != nil {
		return "", err
	}
	reader, err := openPackageArchive(path)
	if err != nil {
		return "", err
	}
	dependency, err := findPackage(reader,packageName)
	if err != nil {
		return "", err
	}

	return dependency, nil

}

