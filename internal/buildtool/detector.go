package buildtool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type DetectionResult struct {
	BuildTool Type
	FilePath  string
	Parser    Parser
}

func Detect(startDir string, fs afero.Fs) (*DetectionResult, error) {
	current := startDir
	for {
		if result := detectInDirectory(current, fs); result != nil {
			return result, nil
		}

		parent := filepath.Dir(current)
		if parent == current || parent == "" {
			break
		}
		current = parent
	}
	return nil, fmt.Errorf("no build file found in %s or any parent directory", startDir)
}

func detectInDirectory(dir string, fs afero.Fs) *DetectionResult {
	pomPath := filepath.Join(dir, "pom.xml")
	if exists(fs, pomPath) {
		parser := GetParser(Maven, fs)
		if parser != nil {
			return &DetectionResult{BuildTool: Maven, FilePath: pomPath, Parser: parser}
		}
	}

	gradleKtsPath := filepath.Join(dir, "build.gradle.kts")
	if exists(fs, gradleKtsPath) {
		parser := GetParser(GradleKotln, fs)
		if parser != nil {
			return &DetectionResult{BuildTool: GradleKotln, FilePath: gradleKtsPath, Parser: parser}
		}
	}

	gradlePath := filepath.Join(dir, "build.gradle")
	if exists(fs, gradlePath) {
		parser := GetParser(Gradle, fs)
		if parser != nil {
			return &DetectionResult{BuildTool: Gradle, FilePath: gradlePath, Parser: parser}
		}
	}

	return nil
}

func exists(fs afero.Fs, path string) bool {
	_, err := fs.Stat(path)
	return err == nil
}

func FindBuildFile(startDir string, fs afero.Fs) (string, Type, error) {
	result, err := Detect(startDir, fs)
	if err != nil {
		return "", "", err
	}
	return result.FilePath, result.BuildTool, nil
}

func DetectFromPath(path string) Type {
	base := filepath.Base(path)
	switch {
	case base == "pom.xml":
		return Maven
	case base == "build.gradle.kts":
		return GradleKotln
	case base == "build.gradle":
		return Gradle
	default:
		return ""
	}
}

func BuildFileNames() []string {
	return []string{"pom.xml", "build.gradle.kts", "build.gradle"}
}

func IsBuildFile(path string) bool {
	base := filepath.Base(path)
	for _, name := range BuildFileNames() {
		if base == name {
			return true
		}
	}
	return false
}

func GetBuildFileName(t Type) string {
	switch t {
	case Maven:
		return "pom.xml"
	case GradleKotln:
		return "build.gradle.kts"
	case Gradle:
		return "build.gradle"
	default:
		return ""
	}
}

func ParseType(s string) Type {
	switch strings.ToLower(s) {
	case "maven":
		return Maven
	case "gradle":
		return Gradle
	case "gradle-kotlin":
		return GradleKotln
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}

func (t Type) DisplayName() string {
	switch t {
	case Maven:
		return "Maven"
	case Gradle:
		return "Gradle (Groovy)"
	case GradleKotln:
		return "Gradle (Kotlin)"
	default:
		return string(t)
	}
}

func (t Type) IsGradle() bool {
	return t == Gradle || t == GradleKotln
}

func DetectWithCwd(fs afero.Fs) (*DetectionResult, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current directory: %w", err)
	}
	return Detect(cwd, fs)
}
