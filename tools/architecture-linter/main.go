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
	// 解析対象のディレクトリを指定
	var violations []string
	dirs := []string{"internal/domain", "internal/usecase", "internal/infrastructure"}

	for _, dir := range dirs {
		fmt.Printf("Checking directory: %s\n", dir)

		// ディレクトリを再帰的に検査
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
					importPath := imp.Path.Value[1 : len(imp.Path.Value)-1] // import pathから""を除去
					if strings.HasPrefix(importPath, "github.com/MoneyForest/go-clean-architecture-boilerplate/internal/infrastructure") {
						if strings.HasPrefix(dir, "internal/domain") ||
							strings.HasPrefix(dir, "internal/usecase") {
							violations = append(violations, fmt.Sprintf("      Violation: %s は %s をimportしてはいけません\n", path, importPath))
						}
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
		fmt.Println("Architecture violations found:")
		for _, violation := range violations {
			fmt.Print(violation)
		}
		os.Exit(1)
	}
}
