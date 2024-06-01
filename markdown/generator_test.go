package markdown

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	}(tempDir)

	// Create a sample app.yml file
	createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description", []string{"darwin_arm64", "darwin_amd64"})
	createTestAppYaml(t, tempDir, "app1", "v2.0", "Test app 1 description", []string{"darwin_arm64"})
	createTestAppYaml(t, tempDir, "app2", "v1.0", "Test app 2 description", []string{"darwin_amd64"})

	generator := AppsMarkdownGenerator{}
	markdown, err := generator.Generate(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(markdown, "## 🖥️ app1") || !strings.Contains(markdown, "## 🖥️ app2") {
		t.Fatalf("Generated markdown is missing expected content: %s", markdown)
	}

	if !strings.Contains(markdown, "🏷️ v1.0") || !strings.Contains(markdown, "🏷️ v2.0") {
		t.Fatalf("Generated markdown is missing expected versions: %s", markdown)
	}

	if !strings.Contains(markdown, "🍎 darwin_arm64") || !strings.Contains(markdown, "🍏 darwin_amd64") {
		t.Fatalf("Generated markdown is missing expected platforms: %s", markdown)
	}
}

func TestScanApps(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("Failed removing all files in path %s: %v", path, err)
		}
	}(tempDir)

	// Create a sample app.yml file
	createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description", []string{"darwin_arm64", "darwin_amd64"})
	createTestAppYaml(t, tempDir, "app1", "v2.0", "Test app 1 description", []string{"darwin_arm64"})
	createTestAppYaml(t, tempDir, "app2", "v1.0", "Test app 2 description", []string{"darwin_amd64"})

	generator := AppsMarkdownGenerator{}
	apps, err := generator.scanApps(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(apps) != 2 {
		t.Fatalf("Expected 2 apps, got %d", len(apps))
	}
}

func TestGetDescriptionAndPlatformsFromAppYaml(t *testing.T) {
	// Setup
	tempDir := setupTestDir(t)
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			t.Fatalf("Failed removing all files in path %s: %v", path, err)
		}
	}(tempDir)

	filePath := createTestAppYaml(t, tempDir, "app1", "v1.0", "Test app 1 description", []string{"darwin_arm64", "darwin_amd64"})

	generator := AppsMarkdownGenerator{}
	description, platforms, err := generator.getDescriptionAndPlatformsFromAppYaml(filePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedDescription := "Test app 1 description"
	expectedPlatforms := []string{"darwin_arm64", "darwin_amd64"}
	if description != expectedDescription {
		t.Fatalf("Expected description '%s', got '%s'", expectedDescription, description)
	}
	if !equalSlices(platforms, expectedPlatforms) {
		t.Fatalf("Expected platforms '%v', got '%v'", expectedPlatforms, platforms)
	}
}

func TestGenerateMarkdown(t *testing.T) {
	apps := []App{
		{Name: "app1", Description: "Test app 1 description", Versions: []Version{
			{VersionName: "v1.0", Platforms: []string{"darwin_arm64", "darwin_amd64"}},
			{VersionName: "v2.0", Platforms: []string{"darwin_arm64"}},
		}},
		{Name: "app2", Description: "Test app 2 description", Versions: []Version{
			{VersionName: "v1.0", Platforms: []string{"darwin_amd64"}},
		}},
	}

	generator := AppsMarkdownGenerator{}
	markdown, err := generator.generateMarkdown(apps)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(markdown, "## 🖥️ app1") || !strings.Contains(markdown, "## 🖥️ app2") {
		t.Fatalf("Generated markdown is missing expected content: %s", markdown)
	}

	if !strings.Contains(markdown, "🏷️ v1.0") || !strings.Contains(markdown, "🏷️ v2.0") {
		t.Fatalf("Generated markdown is missing expected versions: %s", markdown)
	}

	if !strings.Contains(markdown, "🍎 darwin_arm64") || !strings.Contains(markdown, "🍏 darwin_amd64") {
		t.Fatalf("Generated markdown is missing expected platforms: %s", markdown)
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

func createTestAppYaml(t *testing.T, baseDir, appName, version, description string, platforms []string) string {
	appDir := filepath.Join(baseDir, appName, version)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		t.Fatalf("Unable to create app dir: %v", err)
	}

	platformsYaml := ""
	for _, platform := range platforms {
		platformsYaml += fmt.Sprintf("  %s:\n", platform)
	}

	appYamlPath := filepath.Join(appDir, "app.yml")
	content := fmt.Sprintf("description: %s\nplatforms:\n%s", description, platformsYaml)
	if err := os.WriteFile(appYamlPath, []byte(content), 0644); err != nil {
		t.Fatalf("Unable to write app.yml: %v", err)
	}

	return appYamlPath
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
