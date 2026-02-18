package detector

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

type Scanner struct {
	fs         afero.Fs
	projectDir string
	sourceRoot string
	testRoot   string
}

type ScanResult struct {
	SourceFiles []*JavaFile
	TestFiles   []*JavaFile
	BasePackage string
	SourceRoot  string
	TestRoot    string
	BuildTool   string
	HasGradle   bool
	HasMaven    bool
}

func NewScanner(fs afero.Fs, projectDir string) *Scanner {
	return &Scanner{
		fs:         fs,
		projectDir: projectDir,
		sourceRoot: "src/main/java",
		testRoot:   "src/test/java",
	}
}

func (s *Scanner) Scan() (*ScanResult, error) {
	result := &ScanResult{
		SourceRoot: s.sourceRoot,
		TestRoot:   s.testRoot,
	}

	s.detectBuildTool(result)

	sourcePath := filepath.Join(s.projectDir, s.sourceRoot)
	if exists, _ := afero.DirExists(s.fs, sourcePath); exists {
		files, err := s.scanDirectory(sourcePath, false)
		if err != nil {
			return nil, err
		}
		result.SourceFiles = files
	}

	testPath := filepath.Join(s.projectDir, s.testRoot)
	if exists, _ := afero.DirExists(s.fs, testPath); exists {
		files, err := s.scanDirectory(testPath, true)
		if err != nil {
			return nil, err
		}
		result.TestFiles = files
	}

	result.BasePackage = s.detectBasePackage(result.SourceFiles)

	return result, nil
}

func (s *Scanner) detectBuildTool(result *ScanResult) {
	pomPath := filepath.Join(s.projectDir, "pom.xml")
	if exists, _ := afero.Exists(s.fs, pomPath); exists {
		result.HasMaven = true
		result.BuildTool = "maven"
	}

	gradlePath := filepath.Join(s.projectDir, "build.gradle")
	gradleKtsPath := filepath.Join(s.projectDir, "build.gradle.kts")
	if exists, _ := afero.Exists(s.fs, gradlePath); exists {
		result.HasGradle = true
		result.BuildTool = "gradle"
	} else if exists, _ := afero.Exists(s.fs, gradleKtsPath); exists {
		result.HasGradle = true
		result.BuildTool = "gradle"
	}
}

func (s *Scanner) scanDirectory(dir string, isTest bool) ([]*JavaFile, error) {
	var files []*JavaFile

	err := afero.Walk(s.fs, dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".java") {
			return nil
		}

		javaFile, err := s.parseJavaFile(path, dir, isTest)
		if err != nil {
			return nil
		}

		files = append(files, javaFile)
		return nil
	})

	return files, err
}

func (s *Scanner) parseJavaFile(path, rootDir string, isTest bool) (*JavaFile, error) {
	file, err := s.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	jf := &JavaFile{
		Path:     path,
		FileType: FileTypeUnknown,
	}

	if isTest {
		jf.FileType = FileTypeTest
	}

	scanner := bufio.NewScanner(file)
	inBlockComment := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "/*") {
			inBlockComment = true
		}
		if strings.Contains(line, "*/") {
			inBlockComment = false
			continue
		}
		if inBlockComment || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "package ") {
			jf.Package = extractPackage(line)
			continue
		}

		if strings.HasPrefix(line, "import ") {
			jf.Imports = append(jf.Imports, extractImport(line))
			continue
		}

		if strings.HasPrefix(line, "@") {
			annotation := extractAnnotation(line)
			if annotation != "" {
				jf.Annotations = append(jf.Annotations, annotation)
			}
			continue
		}

		if isClassDeclaration(line) {
			parseClassDeclaration(line, jf)
			break
		}

		if isInterfaceDeclaration(line) {
			jf.IsInterface = true
			jf.ClassName = extractInterfaceName(line)
			parseInterfaceExtends(line, jf)
			break
		}
	}

	if jf.ClassName == "" {
		jf.ClassName = extractClassNameFromPath(path)
	}

	if !isTest {
		jf.FileType = classifyJavaFile(jf)
	}

	return jf, scanner.Err()
}

func (s *Scanner) detectBasePackage(files []*JavaFile) string {
	if len(files) == 0 {
		return ""
	}

	packageCounts := make(map[string]int)

	for _, f := range files {
		if f.Package == "" {
			continue
		}

		parts := strings.Split(f.Package, ".")
		for i := 1; i <= len(parts); i++ {
			prefix := strings.Join(parts[:i], ".")
			packageCounts[prefix]++
		}
	}

	var basePackage string
	maxDepth := 0

	for pkg, count := range packageCounts {
		if count == len(files) {
			depth := strings.Count(pkg, ".") + 1
			if depth > maxDepth {
				maxDepth = depth
				basePackage = pkg
			}
		}
	}

	if basePackage == "" && len(files) > 0 {
		basePackage = findShortestCommonPrefix(files)
	}

	return basePackage
}

func findShortestCommonPrefix(files []*JavaFile) string {
	if len(files) == 0 {
		return ""
	}

	var packages []string
	for _, f := range files {
		if f.Package != "" {
			packages = append(packages, f.Package)
		}
	}

	if len(packages) == 0 {
		return ""
	}

	prefix := packages[0]
	for _, pkg := range packages[1:] {
		prefix = commonPrefix(prefix, pkg)
	}

	if idx := strings.LastIndex(prefix, "."); idx > 0 {
		prefix = prefix[:idx]
	}

	return prefix
}

func commonPrefix(a, b string) string {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	var common []string
	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		if partsA[i] == partsB[i] {
			common = append(common, partsA[i])
		} else {
			break
		}
	}

	return strings.Join(common, ".")
}

var (
	packageRegex     = regexp.MustCompile(`^package\s+([a-zA-Z0-9_.]+)\s*;`)
	importRegex      = regexp.MustCompile(`^import\s+(?:static\s+)?([a-zA-Z0-9_.]+)\s*;`)
	annotationRegex  = regexp.MustCompile(`^@(\w+)`)
	classRegex       = regexp.MustCompile(`(?:public\s+)?(?:abstract\s+)?class\s+(\w+)`)
	extendsRegex     = regexp.MustCompile(`extends\s+(\w+)`)
	implementsRegex  = regexp.MustCompile(`implements\s+([\w\s,<>]+)`)
	interfaceRegex   = regexp.MustCompile(`(?:public\s+)?interface\s+(\w+)`)
	interfaceExtends = regexp.MustCompile(`extends\s+([\w\s,<>]+)`)
)

func extractPackage(line string) string {
	matches := packageRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractImport(line string) string {
	matches := importRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractAnnotation(line string) string {
	matches := annotationRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func isClassDeclaration(line string) bool {
	return strings.Contains(line, "class ") && !strings.HasPrefix(line, "//")
}

func isInterfaceDeclaration(line string) bool {
	return strings.Contains(line, "interface ") && !strings.HasPrefix(line, "//")
}

func parseClassDeclaration(line string, jf *JavaFile) {
	if strings.Contains(line, "abstract ") {
		jf.IsAbstract = true
	}

	classMatches := classRegex.FindStringSubmatch(line)
	if len(classMatches) > 1 {
		jf.ClassName = classMatches[1]
	}

	extendsMatches := extendsRegex.FindStringSubmatch(line)
	if len(extendsMatches) > 1 {
		jf.ExtendsClass = extendsMatches[1]
	}

	implementsMatches := implementsRegex.FindStringSubmatch(line)
	if len(implementsMatches) > 1 {
		interfaces := strings.Split(implementsMatches[1], ",")
		for _, iface := range interfaces {
			iface = strings.TrimSpace(iface)
			iface = strings.Split(iface, "<")[0]
			if iface != "" {
				jf.ImplementsInterfaces = append(jf.ImplementsInterfaces, iface)
			}
		}
	}
}

func extractInterfaceName(line string) string {
	matches := interfaceRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func parseInterfaceExtends(line string, jf *JavaFile) {
	matches := interfaceExtends.FindStringSubmatch(line)
	if len(matches) > 1 {
		interfaces := strings.Split(matches[1], ",")
		for _, iface := range interfaces {
			iface = strings.TrimSpace(iface)
			iface = strings.Split(iface, "<")[0]
			if iface != "" {
				jf.ImplementsInterfaces = append(jf.ImplementsInterfaces, iface)
			}
		}
	}
}

func extractClassNameFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".java")
}

func classifyJavaFile(jf *JavaFile) JavaFileType {
	for _, ann := range jf.Annotations {
		switch ann {
		case "RestController", "Controller":
			return FileTypeController
		case "Service":
			return FileTypeService
		case "Repository":
			return FileTypeRepository
		case "Entity", "Document", "Table":
			return FileTypeEntity
		case "Mapper":
			return FileTypeMapper
		case "Configuration", "Component", "Bean":
			return FileTypeConfig
		case "ControllerAdvice", "RestControllerAdvice":
			return FileTypeException
		}
	}

	for _, iface := range jf.ImplementsInterfaces {
		if strings.Contains(iface, "Repository") {
			return FileTypeRepository
		}
	}

	className := jf.ClassName

	if strings.HasSuffix(className, "Controller") || strings.HasSuffix(className, "Resource") {
		return FileTypeController
	}
	if strings.HasSuffix(className, "Service") || strings.HasSuffix(className, "ServiceImpl") {
		return FileTypeService
	}
	if strings.HasSuffix(className, "Repository") {
		return FileTypeRepository
	}
	if strings.HasSuffix(className, "Entity") {
		return FileTypeEntity
	}
	if strings.HasSuffix(className, "Mapper") {
		return FileTypeMapper
	}
	if strings.HasSuffix(className, "Exception") {
		return FileTypeException
	}
	if strings.HasSuffix(className, "Config") || strings.HasSuffix(className, "Configuration") {
		return FileTypeConfig
	}

	if strings.HasSuffix(className, "Request") || strings.HasSuffix(className, "Response") ||
		strings.HasSuffix(className, "DTO") || strings.HasSuffix(className, "Dto") {
		return FileTypeDTO
	}

	return FileTypeUnknown
}

func (s *Scanner) GetFilesByType(files []*JavaFile, fileType JavaFileType) []*JavaFile {
	var result []*JavaFile
	for _, f := range files {
		if f.FileType == fileType {
			result = append(result, f)
		}
	}
	return result
}

func (s *Scanner) GetFilesByAnnotation(files []*JavaFile, annotation string) []*JavaFile {
	var result []*JavaFile
	for _, f := range files {
		for _, ann := range f.Annotations {
			if ann == annotation {
				result = append(result, f)
				break
			}
		}
	}
	return result
}

func (s *Scanner) GetFilesExtending(files []*JavaFile, parentClass string) []*JavaFile {
	var result []*JavaFile
	for _, f := range files {
		if f.ExtendsClass == parentClass {
			result = append(result, f)
		}
	}
	return result
}

func (s *Scanner) GroupFilesByPackage(files []*JavaFile) map[string][]*JavaFile {
	groups := make(map[string][]*JavaFile)
	for _, f := range files {
		groups[f.Package] = append(groups[f.Package], f)
	}
	return groups
}

func (s *Scanner) GetUniquePackages(files []*JavaFile) []string {
	seen := make(map[string]bool)
	var packages []string

	for _, f := range files {
		if f.Package != "" && !seen[f.Package] {
			seen[f.Package] = true
			packages = append(packages, f.Package)
		}
	}

	return packages
}
