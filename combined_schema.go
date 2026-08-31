//go:build combined_schema

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func combineGraphqlFiles(file *os.File) error {
	pattern := fmt.Sprintf("pkg/clinical/presentation/graph/*.graphql")

	fileList, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, filename := range fileList {
		content, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		_, err = file.Write(content)
		if err != nil {
			return err
		}

		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fileName := "schema.graphql"

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer file.Close()

	err = combineGraphqlFiles(file)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Files combined successfully.")
}
