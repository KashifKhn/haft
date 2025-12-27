package routes

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	methodStyles = map[string]lipgloss.Style{
		"GET":    lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true),
		"POST":   lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true),
		"PUT":    lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true),
		"PATCH":  lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true),
		"DELETE": lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true),
	}
	pathStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	controllerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	headerStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
)

type Route struct {
	Method     string
	Path       string
	Controller string
	Handler    string
	File       string
	Line       int
}

func NewCommand() *cobra.Command {
	var jsonOutput bool
	var showFiles bool

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List all REST endpoints",
		Long: `Scan the project and list all REST API endpoints.

Parses Java source files to find Spring MVC annotations like
@GetMapping, @PostMapping, @RequestMapping, etc.`,
		Example: `  # List all routes
  haft routes

  # Show file locations
  haft routes --files

  # Output as JSON
  haft routes --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoutes(jsonOutput, showFiles)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	cmd.Flags().BoolVarP(&showFiles, "files", "f", false, "Show file locations")

	return cmd
}

func runRoutes(jsonOutput bool, showFiles bool) error {
	fs := afero.NewOsFs()
	_, err := buildtool.DetectWithCwd(fs)
	if err != nil {
		return fmt.Errorf("not a Spring Boot project: %w", err)
	}

	srcDir := findSourceDir()

	routes, err := scanForRoutes(srcDir)
	if err != nil {
		return fmt.Errorf("failed to scan routes: %w", err)
	}

	if len(routes) == 0 {
		fmt.Println("No routes found.")
		return nil
	}

	sortRoutes(routes)

	if jsonOutput {
		return printRoutesJSON(routes)
	}

	return printRoutesFormatted(routes, showFiles)
}

func findSourceDir() string {
	possibleDirs := []string{
		"src/main/java",
		"src/main/kotlin",
	}

	for _, dir := range possibleDirs {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	return "src/main/java"
}

func scanForRoutes(srcDir string) ([]Route, error) {
	var routes []Route

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || (!strings.HasSuffix(path, ".java") && !strings.HasSuffix(path, ".kt")) {
			return nil
		}

		fileRoutes, err := parseFileForRoutes(path)
		if err != nil {
			return nil
		}

		routes = append(routes, fileRoutes...)
		return nil
	})

	return routes, err
}

func parseFileForRoutes(filePath string) ([]Route, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var routes []Route
	var classPath string
	var controllerName string
	scanner := bufio.NewScanner(file)
	lineNum := 0

	requestMappingRe := regexp.MustCompile(`@RequestMapping\s*\(\s*(?:value\s*=\s*)?["']([^"']+)["']`)
	classMappingRe := regexp.MustCompile(`@RequestMapping\s*\(\s*(?:value\s*=\s*)?["']([^"']+)["']`)
	classNameRe := regexp.MustCompile(`class\s+(\w+)`)

	getMappingRe := regexp.MustCompile(`@GetMapping\s*(?:\(\s*(?:value\s*=\s*)?["']?([^"'\)]*?)["']?\s*\))?`)
	postMappingRe := regexp.MustCompile(`@PostMapping\s*(?:\(\s*(?:value\s*=\s*)?["']?([^"'\)]*?)["']?\s*\))?`)
	putMappingRe := regexp.MustCompile(`@PutMapping\s*(?:\(\s*(?:value\s*=\s*)?["']?([^"'\)]*?)["']?\s*\))?`)
	patchMappingRe := regexp.MustCompile(`@PatchMapping\s*(?:\(\s*(?:value\s*=\s*)?["']?([^"'\)]*?)["']?\s*\))?`)
	deleteMappingRe := regexp.MustCompile(`@DeleteMapping\s*(?:\(\s*(?:value\s*=\s*)?["']?([^"'\)]*?)["']?\s*\))?`)

	methodRe := regexp.MustCompile(`(?:public|private|protected)\s+\S+\s+(\w+)\s*\(`)

	var pendingMethod string
	var pendingPath string
	var pendingLine int

	isController := false
	fileContent, _ := os.ReadFile(filePath)
	fileStr := string(fileContent)

	if strings.Contains(fileStr, "@RestController") || strings.Contains(fileStr, "@Controller") {
		isController = true
	}

	if !isController {
		return nil, nil
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if matches := classNameRe.FindStringSubmatch(line); matches != nil && controllerName == "" {
			controllerName = matches[1]
		}

		if matches := classMappingRe.FindStringSubmatch(line); matches != nil && classPath == "" {
			classPath = matches[1]
		}

		if matches := getMappingRe.FindStringSubmatch(line); matches != nil {
			pendingMethod = "GET"
			pendingPath = cleanPath(matches[1])
			pendingLine = lineNum
		} else if matches := postMappingRe.FindStringSubmatch(line); matches != nil {
			pendingMethod = "POST"
			pendingPath = cleanPath(matches[1])
			pendingLine = lineNum
		} else if matches := putMappingRe.FindStringSubmatch(line); matches != nil {
			pendingMethod = "PUT"
			pendingPath = cleanPath(matches[1])
			pendingLine = lineNum
		} else if matches := patchMappingRe.FindStringSubmatch(line); matches != nil {
			pendingMethod = "PATCH"
			pendingPath = cleanPath(matches[1])
			pendingLine = lineNum
		} else if matches := deleteMappingRe.FindStringSubmatch(line); matches != nil {
			pendingMethod = "DELETE"
			pendingPath = cleanPath(matches[1])
			pendingLine = lineNum
		}

		if matches := requestMappingRe.FindStringSubmatch(line); matches != nil {
			methodMatch := regexp.MustCompile(`method\s*=\s*RequestMethod\.(\w+)`)
			if mm := methodMatch.FindStringSubmatch(line); mm != nil {
				pendingMethod = mm[1]
				pendingPath = cleanPath(matches[1])
				pendingLine = lineNum
			}
		}

		if pendingMethod != "" {
			if matches := methodRe.FindStringSubmatch(line); matches != nil {
				fullPath := joinPaths(classPath, pendingPath)
				routes = append(routes, Route{
					Method:     pendingMethod,
					Path:       fullPath,
					Controller: controllerName,
					Handler:    matches[1],
					File:       filePath,
					Line:       pendingLine,
				})
				pendingMethod = ""
				pendingPath = ""
			}
		}
	}

	return routes, nil
}

func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, "\"'")
	if path == "" {
		return ""
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

func joinPaths(basePath, subPath string) string {
	if basePath == "" && subPath == "" {
		return "/"
	}
	if basePath == "" {
		if subPath == "" {
			return "/"
		}
		return subPath
	}
	if subPath == "" {
		return basePath
	}

	basePath = strings.TrimSuffix(basePath, "/")
	if !strings.HasPrefix(subPath, "/") {
		subPath = "/" + subPath
	}

	return basePath + subPath
}

func sortRoutes(routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}
		methodOrder := map[string]int{"GET": 1, "POST": 2, "PUT": 3, "PATCH": 4, "DELETE": 5}
		return methodOrder[routes[i].Method] < methodOrder[routes[j].Method]
	})
}

func printRoutesFormatted(routes []Route, showFiles bool) error {
	fmt.Println()
	fmt.Println(headerStyle.Render("  API Routes"))
	fmt.Println(strings.Repeat("─", 70))

	for _, r := range routes {
		methodStyle := methodStyles[r.Method]
		if methodStyle.Value() == "" {
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		}

		method := methodStyle.Render(fmt.Sprintf("%-7s", r.Method))
		path := pathStyle.Render(r.Path)
		handler := controllerStyle.Render(fmt.Sprintf("%s.%s", r.Controller, r.Handler))

		if showFiles {
			fmt.Printf("  %s %s\n", method, path)
			fmt.Printf("          %s (%s:%d)\n", handler, r.File, r.Line)
		} else {
			fmt.Printf("  %s %-35s %s\n", method, path, handler)
		}
	}

	fmt.Println()
	fmt.Printf("  Total: %d routes\n\n", len(routes))

	return nil
}

func printRoutesJSON(routes []Route) error {
	fmt.Println("[")
	for i, r := range routes {
		comma := ","
		if i == len(routes)-1 {
			comma = ""
		}
		fmt.Printf(`  {"method": "%s", "path": "%s", "controller": "%s", "handler": "%s", "file": "%s", "line": %d}%s
`, r.Method, r.Path, r.Controller, r.Handler, r.File, r.Line, comma)
	}
	fmt.Println("]")
	return nil
}
