package markdown

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}(tempDir)

	// Create a sample app.yml file
	createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description")
	createTestAppYaml(t, tempDir, "app1", "v2.0", "Test app 1 description")
	createTestAppYaml(t, tempDir, "app2", "v1.0", "Test app 2 description")

	generator := AppsMarkdownGenerator{}
	markdown, err := generator.Generate(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(markdown, "## app1") || !strings.Contains(markdown, "## app2") {
		t.Fatalf("Generated markdown is missing expected content: %s", markdown)
	}

	if !strings.Contains(markdown, "v1.0") || !strings.Contains(markdown, "v2.0") {
		t.Fatalf("Generated markdown is missing expected versions: %s", markdown)
	}
}

func TestScanApps(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer os.RemoveAll(tempDir)

	// Create a sample app.yml file
	createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description")
	createTestAppYaml(t, tempDir, "app1", "v2.0", "Test app 1 description")
	createTestAppYaml(t, tempDir, "app2", "v1.0", "Test app 2 description")

	generator := AppsMarkdownGenerator{}
	apps, err := generator.scanApps(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(apps) != 2 {
		t.Fatalf("Expected 2 apps, got %d", len(apps))
	}
}

func TestGetDescriptionFromAppYaml(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description")

	generator := AppsMarkdownGenerator{}
	description, err := generator.getDescriptionFromAppYaml(filePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDescription := "Test app 1 description"
	if description != expectedDescription {
		t.Fatalf("Expected description '%s', got '%s'", expectedDescription, description)
	}
}

func TestGenerateMarkdown(t *testing.T) {
	apps := []App{
		{Name: "app1", Description: "Test app 1 description", Versions: []string{"v1.0", "v2.0"}},
		{Name: "app2", Description: "Test app 2 description", Versions: []string{"v1.0"}},
	}

	generator := AppsMarkdownGenerator{}
	markdown, err := generator.generateMarkdown(apps)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(markdown, "## app1") || !strings.Contains(markdown, "## app2") {
		t.Fatalf("Generated markdown is missing expected content: %s", markdown)
	}

	if !strings.Contains(markdown, "v1.0") || !strings.Contains(markdown, "v2.0") {
		t.Fatalf("Generated markdown is missing expected versions: %s", markdown)
	}
}

// Helper functions

func setupTestDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Unable to create temp dir: %v", err)
	}
	return tempDir
}

func createTestAppYaml(t *testing.T, baseDir, appName, version, description string) string {
	appDir := filepath.Join(baseDir, appName, version)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("Unable to create app dir: %v", err)
	}

	appYamlPath := filepath.Join(appDir, "app.yml")
	content := fmt.Sprintf("description: %s", description)
	if err := ioutil.WriteFile(appYamlPath, []byte(content), 0644); err != nil {
		t.Fatalf("Unable to write app.yml: %v", err)
	}

	return appYamlPath
}
