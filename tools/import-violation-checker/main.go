package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var violations []string
	dirs := []string{"internal/domain", "internal/usecase", "internal/infrastructure"}

	for _, dir := range dirs {
		fmt.Printf("Checking directory: %s\n", dir)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				fset := token.NewFileSet()
				file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
				if err != nil {
					return err
				}

				fmt.Printf("  Checking file: %s\n", path)

				for _, imp := range file.Imports {
					importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
					violation := checkImportViolation(dir, path, importPath)
					if violation != "" {
						violations = append(violations, violation)
					}
				}
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(violations) > 0 {
		fmt.Println("Clean Architecture import violations found:")
		for _, violation := range violations {
			fmt.Print(violation)
		}
		os.Exit(1)
	}
}

func checkImportViolation(dir, path, importPath string) string {
	switch {
	case strings.HasPrefix(dir, "internal/domain"):
		switch {
		case strings.Contains(importPath, "internal/usecase"):
			return fmt.Sprintf("Domain layer should not import usecase layer: %s imports %s\n", path, importPath)
		case strings.Contains(importPath, "internal/infrastructure"):
			return fmt.Sprintf("Domain layer should not import infrastructure layer: %s imports %s\n", path, importPath)
		}
		return ""
	case strings.HasPrefix(dir, "internal/usecase"):
		switch {
		case strings.Contains(importPath, "internal/infrastructure"):
			return fmt.Sprintf("Usecase layer should not import infrastructure layer: %s imports %s\n", path, importPath)
		}
		return ""
	default:
		return ""
	}
}
