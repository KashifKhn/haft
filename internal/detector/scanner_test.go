package detector

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestFS() afero.Fs {
	return afero.NewMemMapFs()
}

func createJavaFile(fs afero.Fs, path, content string) error {
	dir := filepath.Dir(path)
	if err := fs.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return afero.WriteFile(fs, path, []byte(content), 0644)
}

func TestScannerNewScanner(t *testing.T) {
	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	assert.NotNil(t, scanner)
	assert.Equal(t, "/project", scanner.projectDir)
	assert.Equal(t, "src/main/java", scanner.sourceRoot)
	assert.Equal(t, "src/test/java", scanner.testRoot)
}

func TestScannerDetectBuildTool(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		expectMvn  bool
		expectGrd  bool
		expectTool string
	}{
		{
			name:       "maven project",
			files:      []string{"/project/pom.xml"},
			expectMvn:  true,
			expectGrd:  false,
			expectTool: "maven",
		},
		{
			name:       "gradle project",
			files:      []string{"/project/build.gradle"},
			expectMvn:  false,
			expectGrd:  true,
			expectTool: "gradle",
		},
		{
			name:       "gradle kts project",
			files:      []string{"/project/build.gradle.kts"},
			expectMvn:  false,
			expectGrd:  true,
			expectTool: "gradle",
		},
		{
			name:       "both maven and gradle",
			files:      []string{"/project/pom.xml", "/project/build.gradle"},
			expectMvn:  true,
			expectGrd:  true,
			expectTool: "gradle",
		},
		{
			name:       "no build tool",
			files:      []string{},
			expectMvn:  false,
			expectGrd:  false,
			expectTool: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := createTestFS()
			for _, f := range tt.files {
				require.NoError(t, afero.WriteFile(fs, f, []byte{}, 0644))
			}

			scanner := NewScanner(fs, "/project")
			result := &ScanResult{}
			scanner.detectBuildTool(result)

			assert.Equal(t, tt.expectMvn, result.HasMaven)
			assert.Equal(t, tt.expectGrd, result.HasGradle)
			assert.Equal(t, tt.expectTool, result.BuildTool)
		})
	}
}

func TestScannerParseJavaFilePackage(t *testing.T) {
	fs := createTestFS()
	content := `package com.example.app;

import java.util.List;

public class MyClass {
}`

	path := "/project/src/main/java/com/example/app/MyClass.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.Equal(t, "com.example.app", jf.Package)
	assert.Equal(t, "MyClass", jf.ClassName)
}

func TestScannerParseJavaFileAnnotations(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectedAnn  []string
		expectedType JavaFileType
	}{
		{
			name: "controller",
			content: `package com.example;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class UserController {
}`,
			expectedAnn:  []string{"RestController", "RequestMapping"},
			expectedType: FileTypeController,
		},
		{
			name: "service",
			content: `package com.example;

@Service
public class UserService {
}`,
			expectedAnn:  []string{"Service"},
			expectedType: FileTypeService,
		},
		{
			name: "repository",
			content: `package com.example;

@Repository
public interface UserRepository {
}`,
			expectedAnn:  []string{"Repository"},
			expectedType: FileTypeRepository,
		},
		{
			name: "entity",
			content: `package com.example;

@Entity
@Table(name = "users")
public class User {
}`,
			expectedAnn:  []string{"Entity", "Table"},
			expectedType: FileTypeEntity,
		},
		{
			name: "mapper",
			content: `package com.example;

@Mapper
public interface UserMapper {
}`,
			expectedAnn:  []string{"Mapper"},
			expectedType: FileTypeMapper,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := createTestFS()
			path := "/project/src/main/java/Test.java"
			require.NoError(t, createJavaFile(fs, path, tt.content))

			scanner := NewScanner(fs, "/project")
			jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedAnn, jf.Annotations)
			assert.Equal(t, tt.expectedType, jf.FileType)
		})
	}
}

func TestScannerParseJavaFileInheritance(t *testing.T) {
	fs := createTestFS()
	content := `package com.example;

public class User extends BaseEntity implements Auditable, Serializable {
}`

	path := "/project/src/main/java/User.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.Equal(t, "User", jf.ClassName)
	assert.Equal(t, "BaseEntity", jf.ExtendsClass)
	assert.Contains(t, jf.ImplementsInterfaces, "Auditable")
	assert.Contains(t, jf.ImplementsInterfaces, "Serializable")
}

func TestScannerParseJavaFileAbstractClass(t *testing.T) {
	fs := createTestFS()
	content := `package com.example;

public abstract class BaseEntity {
}`

	path := "/project/src/main/java/BaseEntity.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.True(t, jf.IsAbstract)
	assert.Equal(t, "BaseEntity", jf.ClassName)
}

func TestScannerParseJavaFileInterface(t *testing.T) {
	fs := createTestFS()
	content := `package com.example;

public interface UserRepository extends JpaRepository<User, Long> {
}`

	path := "/project/src/main/java/UserRepository.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.True(t, jf.IsInterface)
	assert.Equal(t, "UserRepository", jf.ClassName)
	assert.Contains(t, jf.ImplementsInterfaces, "JpaRepository")
}

func TestScannerParseJavaFileImports(t *testing.T) {
	fs := createTestFS()
	content := `package com.example;

import java.util.List;
import java.util.UUID;
import org.springframework.stereotype.Service;
import static org.junit.Assert.assertEquals;

@Service
public class MyService {
}`

	path := "/project/src/main/java/MyService.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.Contains(t, jf.Imports, "java.util.List")
	assert.Contains(t, jf.Imports, "java.util.UUID")
	assert.Contains(t, jf.Imports, "org.springframework.stereotype.Service")
	assert.Contains(t, jf.Imports, "org.junit.Assert.assertEquals")
}

func TestScannerParseJavaFileWithComments(t *testing.T) {
	fs := createTestFS()
	content := `package com.example;

// This is a comment
/* Block comment
   spanning multiple lines
*/
/**
 * Javadoc comment
 */
@Service
public class MyService {
}`

	path := "/project/src/main/java/MyService.java"
	require.NoError(t, createJavaFile(fs, path, content))

	scanner := NewScanner(fs, "/project")
	jf, err := scanner.parseJavaFile(path, "/project/src/main/java", false)

	require.NoError(t, err)
	assert.Equal(t, "MyService", jf.ClassName)
	assert.Contains(t, jf.Annotations, "Service")
}

func TestScannerClassifyByName(t *testing.T) {
	tests := []struct {
		className    string
		expectedType JavaFileType
	}{
		{"UserController", FileTypeController},
		{"UserResource", FileTypeController},
		{"UserService", FileTypeService},
		{"UserServiceImpl", FileTypeService},
		{"UserRepository", FileTypeRepository},
		{"UserEntity", FileTypeEntity},
		{"UserMapper", FileTypeMapper},
		{"UserException", FileTypeException},
		{"UserConfig", FileTypeConfig},
		{"UserConfiguration", FileTypeConfig},
		{"UserRequest", FileTypeDTO},
		{"UserResponse", FileTypeDTO},
		{"UserDTO", FileTypeDTO},
		{"UserDto", FileTypeDTO},
		{"SomeRandomClass", FileTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.className, func(t *testing.T) {
			jf := &JavaFile{ClassName: tt.className}
			result := classifyJavaFile(jf)
			assert.Equal(t, tt.expectedType, result)
		})
	}
}

func TestScannerDetectBasePackage(t *testing.T) {
	tests := []struct {
		name     string
		packages []string
		expected string
	}{
		{
			name:     "single package",
			packages: []string{"com.example.app"},
			expected: "com.example.app",
		},
		{
			name: "common prefix",
			packages: []string{
				"com.example.app.controller",
				"com.example.app.service",
				"com.example.app.entity",
			},
			expected: "com.example.app",
		},
		{
			name: "feature packages",
			packages: []string{
				"com.example.app.user.controller",
				"com.example.app.user.service",
				"com.example.app.auth.controller",
				"com.example.app.auth.service",
			},
			expected: "com.example.app",
		},
		{
			name:     "empty",
			packages: []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var files []*JavaFile
			for _, pkg := range tt.packages {
				files = append(files, &JavaFile{Package: pkg})
			}

			fs := createTestFS()
			scanner := NewScanner(fs, "/project")
			result := scanner.detectBasePackage(files)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestScannerFullScan(t *testing.T) {
	fs := createTestFS()

	files := map[string]string{
		"/project/src/main/java/com/example/app/controller/UserController.java": `package com.example.app.controller;

@RestController
public class UserController {
}`,
		"/project/src/main/java/com/example/app/service/UserService.java": `package com.example.app.service;

@Service
public class UserService {
}`,
		"/project/src/main/java/com/example/app/entity/User.java": `package com.example.app.entity;

@Entity
public class User extends BaseEntity {
}`,
		"/project/src/main/java/com/example/app/repository/UserRepository.java": `package com.example.app.repository;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
}`,
		"/project/src/test/java/com/example/app/service/UserServiceTest.java": `package com.example.app.service;

@SpringBootTest
public class UserServiceTest {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	scanner := NewScanner(fs, "/project")
	result, err := scanner.Scan()

	require.NoError(t, err)
	assert.Equal(t, 4, len(result.SourceFiles))
	assert.Equal(t, 1, len(result.TestFiles))
	assert.Equal(t, "com.example.app", result.BasePackage)
	assert.True(t, result.HasMaven)
	assert.Equal(t, "maven", result.BuildTool)
}

func TestScannerGetFilesByType(t *testing.T) {
	files := []*JavaFile{
		{ClassName: "UserController", FileType: FileTypeController},
		{ClassName: "AuthController", FileType: FileTypeController},
		{ClassName: "UserService", FileType: FileTypeService},
		{ClassName: "User", FileType: FileTypeEntity},
	}

	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	controllers := scanner.GetFilesByType(files, FileTypeController)
	assert.Equal(t, 2, len(controllers))

	services := scanner.GetFilesByType(files, FileTypeService)
	assert.Equal(t, 1, len(services))

	repos := scanner.GetFilesByType(files, FileTypeRepository)
	assert.Equal(t, 0, len(repos))
}

func TestScannerGetFilesByAnnotation(t *testing.T) {
	files := []*JavaFile{
		{ClassName: "UserController", Annotations: []string{"RestController", "RequestMapping"}},
		{ClassName: "AuthController", Annotations: []string{"RestController"}},
		{ClassName: "UserService", Annotations: []string{"Service"}},
	}

	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	restControllers := scanner.GetFilesByAnnotation(files, "RestController")
	assert.Equal(t, 2, len(restControllers))

	services := scanner.GetFilesByAnnotation(files, "Service")
	assert.Equal(t, 1, len(services))
}

func TestScannerGetFilesExtending(t *testing.T) {
	files := []*JavaFile{
		{ClassName: "User", ExtendsClass: "BaseEntity"},
		{ClassName: "Product", ExtendsClass: "BaseEntity"},
		{ClassName: "BaseEntity", ExtendsClass: ""},
		{ClassName: "Service", ExtendsClass: "Object"},
	}

	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	extending := scanner.GetFilesExtending(files, "BaseEntity")
	assert.Equal(t, 2, len(extending))
}

func TestScannerGroupFilesByPackage(t *testing.T) {
	files := []*JavaFile{
		{ClassName: "UserController", Package: "com.example.controller"},
		{ClassName: "AuthController", Package: "com.example.controller"},
		{ClassName: "UserService", Package: "com.example.service"},
	}

	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	groups := scanner.GroupFilesByPackage(files)

	assert.Equal(t, 2, len(groups))
	assert.Equal(t, 2, len(groups["com.example.controller"]))
	assert.Equal(t, 1, len(groups["com.example.service"]))
}

func TestScannerGetUniquePackages(t *testing.T) {
	files := []*JavaFile{
		{Package: "com.example.controller"},
		{Package: "com.example.controller"},
		{Package: "com.example.service"},
		{Package: ""},
	}

	fs := createTestFS()
	scanner := NewScanner(fs, "/project")

	packages := scanner.GetUniquePackages(files)

	assert.Equal(t, 2, len(packages))
	assert.Contains(t, packages, "com.example.controller")
	assert.Contains(t, packages, "com.example.service")
}

func TestScannerEmptyProject(t *testing.T) {
	fs := createTestFS()
	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))

	scanner := NewScanner(fs, "/project")
	result, err := scanner.Scan()

	require.NoError(t, err)
	assert.Equal(t, 0, len(result.SourceFiles))
	assert.Equal(t, "", result.BasePackage)
}

func TestExtractPackageRegex(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"package com.example.app;", "com.example.app"},
		{"package com.example;", "com.example"},
		{"package myapp;", "myapp"},
		{"  package com.example.app;  ", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := extractPackage(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractImportRegex(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"import java.util.List;", "java.util.List"},
		{"import static org.junit.Assert.assertEquals;", "org.junit.Assert.assertEquals"},
		{"import com.example.MyClass;", "com.example.MyClass"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := extractImport(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractAnnotationRegex(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"@Service", "Service"},
		{"@RestController", "RestController"},
		{"@RequestMapping(\"/api\")", "RequestMapping"},
		{"@Entity(name = \"user\")", "Entity"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := extractAnnotation(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}
