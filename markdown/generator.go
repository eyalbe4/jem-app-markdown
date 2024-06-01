package markdown

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AppsMarkdownGenerator struct{}

type Version struct {
	VersionName string
	Platforms   []string
}

type App struct {
	Name        string
	Description string
	Versions    []Version
}

func (gen *AppsMarkdownGenerator) Generate(pathToAppsDir string) (markdown string, err error) {
	apps, err := gen.scanApps(pathToAppsDir)
	if err != nil {
		return
	}

	markdown, err = gen.generateMarkdown(apps)
	return
}

func (gen *AppsMarkdownGenerator) scanApps(dir string) ([]App, error) {
	var apps []App
	appsMap := make(map[string]App)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "app.yml" {
			appDir := filepath.Dir(path)
			versionDir := filepath.Base(appDir)
			nameDir := filepath.Base(filepath.Dir(appDir))

			description, platforms, err := gen.getDescriptionAndPlatformsFromAppYaml(path)
			if err != nil {
				return err
			}

			version := Version{
				VersionName: versionDir,
				Platforms:   platforms,
			}

			// Check if the app with the same name already exists
			if app, exists := appsMap[nameDir]; exists {
				// Add the version to the existing app
				app.Versions = append(app.Versions, version)
				// Update the value in the map
				appsMap[nameDir] = app
			} else {
				// Create a new app and add it to the map
				app := App{
					Name:        nameDir,
					Description: description,
					Versions:    []Version{version},
				}
				appsMap[nameDir] = app
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map values to slice
	for _, app := range appsMap {
		apps = append(apps, app)
	}

	// Sort the apps slice by app name
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].Name < apps[j].Name
	})

	return apps, nil
}

func (gen *AppsMarkdownGenerator) getDescriptionAndPlatformsFromAppYaml(filePath string) (string, []string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	// Unmarshal the YAML content into the struct
	type yamlContent struct {
		Description string                 `yaml:"description"`
		Platforms   map[string]interface{} `yaml:"platforms"`
	}
	var y yamlContent
	if err = yaml.Unmarshal(content, &y); err != nil {
		return "", nil, err
	}

	var platforms []string
	for platform := range y.Platforms {
		platforms = append(platforms, platform)
	}

	return y.Description, platforms, nil
}

func (gen *AppsMarkdownGenerator) generateMarkdown(apps []App) (string, error) {
	var sb strings.Builder

	// Define a list of available emojis, to be used for the platforms
	emojis := []string{"🍎", "🍏", "🐧", "🪟", "🚀", "📱", "💻", "🖥️", "📀", "🔧"}
	emojiIndex := 0

	// Map to keep track of platform to emoji assignment
	platformToEmoji := make(map[string]string)

	// Write markdown content
	for _, app := range apps {
		sb.WriteString(fmt.Sprintf("## 🖥️ %s\n", app.Name))
		// Use HTML tags for smaller font size and color the description in light blue
		sb.WriteString(fmt.Sprintf("<small style=\"color:lightblue;\">%s</small>\n", app.Description))
		sb.WriteString("<details>\n")
		sb.WriteString("<summary>📦 Versions</summary>\n")
		sb.WriteString("<ul>\n")
		for _, version := range app.Versions {
			sb.WriteString(fmt.Sprintf("<li>🏷️ %s\n", version.VersionName))
			sb.WriteString("<ul>\n")
			for _, platform := range version.Platforms {
				// Assign an emoji to the platform if not already assigned
				emoji, exists := platformToEmoji[platform]
				if !exists {
					if emojiIndex >= len(emojis) {
						emoji = "🔧" // Default emoji if we run out of unique emojis
					} else {
						emoji = emojis[emojiIndex]
						emojiIndex++
					}
					platformToEmoji[platform] = emoji
				}
				sb.WriteString(fmt.Sprintf("<li>%s %s</li>\n", emoji, platform))
			}
			sb.WriteString("</ul>\n")
			sb.WriteString("</li>\n")
		}
		sb.WriteString("</ul>\n")
		sb.WriteString("</details>\n")
		sb.WriteString("\n")
	}

	return sb.String(), nil
}
