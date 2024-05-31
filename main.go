package main

import (
	"fmt"
	"jem-apps-chart/markdown"
	"os"
	"path/filepath"
)

func main() {
	pathToApps := filepath.Join("/", "Users", "eyalb", "dev", "jem-apps", "apps")
	generator := markdown.AppsMarkdownGenerator{}
	content, err := generator.Generate(pathToApps)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(content)
}
