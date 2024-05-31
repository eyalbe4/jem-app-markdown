package main

import (
	"fmt"
	"jem-apps-chart/markdown"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: You should provide the path the 'apps' directory as an argument")
		os.Exit(1)
	}

	pathToApps := os.Args[1]
	generator := markdown.AppsMarkdownGenerator{}
	content, err := generator.Generate(filepath.Clean(pathToApps))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(content)
}
