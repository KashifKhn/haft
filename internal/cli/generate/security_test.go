package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KashifKhn/haft/internal/detector"
	"github.com/KashifKhn/haft/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityCommand(t *testing.T) {
	cmd := newSecurityCommand()

	assert.Equal(t, "security", cmd.Use)
	assert.Contains(t, cmd.Aliases, "sec")
	assert.Contains(t, cmd.Aliases, "auth")
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
}

func TestSecurityCommandFlags(t *testing.T) {
	cmd := newSecurityCommand()

	pkgFlag := cmd.Flags().Lookup("package")
	assert.NotNil(t, pkgFlag)
	assert.Equal(t, "p", pkgFlag.Shorthand)

	jwtFlag := cmd.Flags().Lookup("jwt")
	assert.NotNil(t, jwtFlag)

	sessionFlag := cmd.Flags().Lookup("session")
	assert.NotNil(t, sessionFlag)

	oauth2Flag := cmd.Flags().Lookup("oauth2")
	assert.NotNil(t, oauth2Flag)

	allFlag := cmd.Flags().Lookup("all")
	assert.NotNil(t, allFlag)

	noInteractiveFlag := cmd.Flags().Lookup("no-interactive")
	assert.NotNil(t, noInteractiveFlag)

	skipEntitiesFlag := cmd.Flags().Lookup("skip-entities")
	assert.NotNil(t, skipEntitiesFlag)

	refreshFlag := cmd.Flags().Lookup("refresh")
	assert.NotNil(t, refreshFlag)
}

func TestSecurityTypeConstants(t *testing.T) {
	assert.Equal(t, SecurityType("jwt"), SecurityJWT)
	assert.Equal(t, SecurityType("session"), SecuritySession)
	assert.Equal(t, SecurityType("oauth2"), SecurityOAuth2)
}

func TestSecurityDependenciesDefinition(t *testing.T) {
	assert.GreaterOrEqual(t, len(securityDependencies), 6)

	springSecurityFound := false
	jwtApiFound := false
	oauth2Found := false

	for _, dep := range securityDependencies {
		if dep.ArtifactId == "spring-boot-starter-security" {
			springSecurityFound = true
			assert.Equal(t, "org.springframework.boot", dep.GroupId)
			assert.Nil(t, dep.Required)
		}
		if dep.ArtifactId == "jjwt-api" {
			jwtApiFound = true
			assert.Equal(t, "io.jsonwebtoken", dep.GroupId)
			assert.Contains(t, dep.Required, SecurityJWT)
		}
		if dep.ArtifactId == "spring-boot-starter-oauth2-client" {
			oauth2Found = true
			assert.Equal(t, "org.springframework.boot", dep.GroupId)
			assert.Contains(t, dep.Required, SecurityOAuth2)
		}
	}

	assert.True(t, springSecurityFound, "spring-boot-starter-security not found")
	assert.True(t, jwtApiFound, "jjwt-api not found")
	assert.True(t, oauth2Found, "spring-boot-starter-oauth2-client not found")
}

func TestUserEntityNamesDefinition(t *testing.T) {
	assert.GreaterOrEqual(t, len(userEntityNames), 5)
	assert.Contains(t, userEntityNames, "User")
	assert.Contains(t, userEntityNames, "AppUser")
	assert.Contains(t, userEntityNames, "Account")
}

func TestRoleEntityNamesDefinition(t *testing.T) {
	assert.GreaterOrEqual(t, len(roleEntityNames), 3)
	assert.Contains(t, roleEntityNames, "Role")
	assert.Contains(t, roleEntityNames, "Authority")
}

func TestGetSecurityPackageAllArchitectures(t *testing.T) {
	tests := []struct {
		name        string
		arch        detector.ArchitectureType
		basePackage string
		expected    string
	}{
		{
			name:        "layered",
			arch:        detector.ArchLayered,
			basePackage: "com.example.app",
			expected:    "com.example.app.security",
		},
		{
			name:        "feature",
			arch:        detector.ArchFeature,
			basePackage: "com.example.app",
			expected:    "com.example.app.common.security",
		},
		{
			name:        "hexagonal",
			arch:        detector.ArchHexagonal,
			basePackage: "com.example.app",
			expected:    "com.example.app.infrastructure.security",
		},
		{
			name:        "clean",
			arch:        detector.ArchClean,
			basePackage: "com.example.app",
			expected:    "com.example.app.infrastructure.security",
		},
		{
			name:        "flat",
			arch:        detector.ArchFlat,
			basePackage: "com.example.app",
			expected:    "com.example.app.security",
		},
		{
			name:        "modular",
			arch:        detector.ArchModular,
			basePackage: "com.example.app",
			expected:    "com.example.app.security",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.arch,
				BasePackage:  tt.basePackage,
			}
			result := getSecurityPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUserEntityPackageAllArchitectures(t *testing.T) {
	tests := []struct {
		name        string
		arch        detector.ArchitectureType
		basePackage string
		expected    string
	}{
		{
			name:        "layered",
			arch:        detector.ArchLayered,
			basePackage: "com.example.app",
			expected:    "com.example.app.entity",
		},
		{
			name:        "feature",
			arch:        detector.ArchFeature,
			basePackage: "com.example.app",
			expected:    "com.example.app.user",
		},
		{
			name:        "hexagonal",
			arch:        detector.ArchHexagonal,
			basePackage: "com.example.app",
			expected:    "com.example.app.domain.model",
		},
		{
			name:        "clean",
			arch:        detector.ArchClean,
			basePackage: "com.example.app",
			expected:    "com.example.app.domain.entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.arch,
				BasePackage:  tt.basePackage,
			}
			result := getUserEntityPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUserRepositoryPackageAllArchitectures(t *testing.T) {
	tests := []struct {
		name        string
		arch        detector.ArchitectureType
		basePackage string
		expected    string
	}{
		{
			name:        "layered",
			arch:        detector.ArchLayered,
			basePackage: "com.example.app",
			expected:    "com.example.app.repository",
		},
		{
			name:        "feature",
			arch:        detector.ArchFeature,
			basePackage: "com.example.app",
			expected:    "com.example.app.user",
		},
		{
			name:        "hexagonal",
			arch:        detector.ArchHexagonal,
			basePackage: "com.example.app",
			expected:    "com.example.app.adapter.persistence",
		},
		{
			name:        "clean",
			arch:        detector.ArchClean,
			basePackage: "com.example.app",
			expected:    "com.example.app.infrastructure.persistence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &detector.ProjectProfile{
				Architecture: tt.arch,
				BasePackage:  tt.basePackage,
			}
			result := getUserRepositoryPackage(profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildSecurityTemplateData(t *testing.T) {
	profile := &detector.ProjectProfile{
		BasePackage:     "com.example.myapp",
		Architecture:    detector.ArchLayered,
		Lombok:          detector.LombokProfile{Detected: true},
		HasValidation:   true,
		ValidationStyle: detector.ValidationJakarta,
		IDType:          "Long",
	}

	cfg := securityConfig{
		UserEntityName:   "User",
		GenerateEntities: true,
	}

	data := buildSecurityTemplateData(
		profile,
		"com.example.myapp.security",
		"com.example.myapp.entity",
		"com.example.myapp.repository",
		cfg,
	)

	assert.Equal(t, "com.example.myapp", data["BasePackage"])
	assert.Equal(t, "com.example.myapp.security", data["SecurityPackage"])
	assert.Equal(t, "com.example.myapp.entity", data["UserEntityPackage"])
	assert.Equal(t, "com.example.myapp.repository", data["UserRepositoryPackage"])
	assert.Equal(t, "User", data["UserEntityName"])
	assert.Equal(t, true, data["HasLombok"])
	assert.Equal(t, true, data["HasValidation"])
	assert.Equal(t, "jakarta.validation", data["ValidationImport"])
	assert.Equal(t, "Long", data["IdType"])
	assert.Equal(t, "IDENTITY", data["IdStrategy"])
	assert.NotEmpty(t, data["DefaultJwtSecret"])
}

func TestBuildSecurityTemplateDataWithUUID(t *testing.T) {
	profile := &detector.ProjectProfile{
		BasePackage:     "com.example.myapp",
		Architecture:    detector.ArchLayered,
		Lombok:          detector.LombokProfile{Detected: false},
		HasValidation:   false,
		ValidationStyle: detector.ValidationNone,
		IDType:          "UUID",
	}

	cfg := securityConfig{
		UserEntityName:   "AppUser",
		GenerateEntities: false,
	}

	data := buildSecurityTemplateData(
		profile,
		"com.example.myapp.security",
		"com.example.myapp.entity",
		"com.example.myapp.repository",
		cfg,
	)

	assert.Equal(t, false, data["HasLombok"])
	assert.Equal(t, false, data["HasValidation"])
	assert.Equal(t, "UUID", data["IdType"])
	assert.Equal(t, "UUID", data["IdStrategy"])
	assert.Equal(t, "AppUser", data["UserEntityName"])
}

func TestBuildSecurityTemplateDataJavaxValidation(t *testing.T) {
	profile := &detector.ProjectProfile{
		BasePackage:     "com.example.myapp",
		Architecture:    detector.ArchLayered,
		HasValidation:   true,
		ValidationStyle: detector.ValidationJavax,
	}

	cfg := securityConfig{UserEntityName: "User"}

	data := buildSecurityTemplateData(
		profile,
		"com.example.myapp.security",
		"com.example.myapp.entity",
		"com.example.myapp.repository",
		cfg,
	)

	assert.Equal(t, "javax.validation", data["ValidationImport"])
}

func TestContainsAnySecurityType(t *testing.T) {
	tests := []struct {
		name     string
		haystack []SecurityType
		needles  []SecurityType
		expected bool
	}{
		{
			name:     "match found",
			haystack: []SecurityType{SecurityJWT, SecuritySession},
			needles:  []SecurityType{SecurityJWT},
			expected: true,
		},
		{
			name:     "no match",
			haystack: []SecurityType{SecuritySession},
			needles:  []SecurityType{SecurityJWT},
			expected: false,
		},
		{
			name:     "empty haystack",
			haystack: []SecurityType{},
			needles:  []SecurityType{SecurityJWT},
			expected: false,
		},
		{
			name:     "empty needles",
			haystack: []SecurityType{SecurityJWT},
			needles:  []SecurityType{},
			expected: false,
		},
		{
			name:     "multiple needles one match",
			haystack: []SecurityType{SecurityOAuth2},
			needles:  []SecurityType{SecurityJWT, SecurityOAuth2},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsAnySecurityType(tt.haystack, tt.needles)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatMissingDeps(t *testing.T) {
	deps := []securityDependency{
		{Name: "Spring Security"},
		{Name: "JJWT API"},
		{Name: "JJWT Impl"},
	}

	result := formatMissingDeps(deps)
	assert.Equal(t, "Spring Security, JJWT API, JJWT Impl", result)
}

func TestFormatMissingDepsEmpty(t *testing.T) {
	deps := []securityDependency{}
	result := formatMissingDeps(deps)
	assert.Equal(t, "", result)
}

func TestFormatMissingDepsSingle(t *testing.T) {
	deps := []securityDependency{
		{Name: "Spring Security"},
	}
	result := formatMissingDeps(deps)
	assert.Equal(t, "Spring Security", result)
}

func TestGetSecurityTemplatesJWT(t *testing.T) {
	templates := getSecurityTemplates(SecurityJWT)

	assert.GreaterOrEqual(t, len(templates), 9)

	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		templateNames[tmpl.fileName] = true
	}

	assert.True(t, templateNames["SecurityConfig.java"])
	assert.True(t, templateNames["JwtUtil.java"])
	assert.True(t, templateNames["JwtAuthenticationFilter.java"])
	assert.True(t, templateNames["AuthenticationController.java"])
	assert.True(t, templateNames["AuthRequest.java"])
	assert.True(t, templateNames["AuthResponse.java"])
	assert.True(t, templateNames["RegisterRequest.java"])
	assert.True(t, templateNames["RefreshTokenRequest.java"])
	assert.True(t, templateNames["CustomUserDetailsService.java"])
}

func TestGetSecurityTemplatesSession(t *testing.T) {
	templates := getSecurityTemplates(SecuritySession)

	assert.GreaterOrEqual(t, len(templates), 4)

	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		templateNames[tmpl.fileName] = true
	}

	assert.True(t, templateNames["SecurityConfig.java"])
	assert.True(t, templateNames["CustomUserDetailsService.java"])
	assert.True(t, templateNames["AuthController.java"])
	assert.True(t, templateNames["RegisterRequest.java"])
}

func TestGetSecurityTemplatesOAuth2(t *testing.T) {
	templates := getSecurityTemplates(SecurityOAuth2)

	assert.GreaterOrEqual(t, len(templates), 4)

	templateNames := make(map[string]bool)
	for _, tmpl := range templates {
		templateNames[tmpl.fileName] = true
	}

	assert.True(t, templateNames["SecurityConfig.java"])
	assert.True(t, templateNames["OAuth2UserService.java"])
	assert.True(t, templateNames["OAuth2SuccessHandler.java"])
	assert.True(t, templateNames["OAuth2UserPrincipal.java"])
}

func TestGetSecurityTemplatesUnknown(t *testing.T) {
	templates := getSecurityTemplates(SecurityType("unknown"))
	assert.Nil(t, templates)
}

func TestGenerateJwtSecret(t *testing.T) {
	secret1 := generateJwtSecret()
	secret2 := generateJwtSecret()

	assert.NotEmpty(t, secret1)
	assert.NotEmpty(t, secret2)
	assert.NotEqual(t, secret1, secret2)
	assert.GreaterOrEqual(t, len(secret1), 32)
}

func TestSecurityConfigStruct(t *testing.T) {
	cfg := securityConfig{
		BasePackage:      "com.example.app",
		SecurityTypes:    []SecurityType{SecurityJWT, SecuritySession},
		GenerateEntities: true,
		UserEntityName:   "User",
	}

	assert.Equal(t, "com.example.app", cfg.BasePackage)
	assert.Len(t, cfg.SecurityTypes, 2)
	assert.True(t, cfg.GenerateEntities)
	assert.Equal(t, "User", cfg.UserEntityName)
}

func TestSecurityDependencyStruct(t *testing.T) {
	dep := securityDependency{
		Name:       "Test Dep",
		GroupId:    "com.test",
		ArtifactId: "test-artifact",
		Version:    "1.0.0",
		Required:   []SecurityType{SecurityJWT},
	}

	assert.Equal(t, "Test Dep", dep.Name)
	assert.Equal(t, "com.test", dep.GroupId)
	assert.Equal(t, "test-artifact", dep.ArtifactId)
	assert.Equal(t, "1.0.0", dep.Version)
	assert.Contains(t, dep.Required, SecurityJWT)
}

func TestSecurityMultiSelectWrapperInit(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := securityMultiSelectWrapper{model: model}

	cmd := wrapper.Init()
	assert.Nil(t, cmd)
}

func TestSecurityMultiSelectWrapperView(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test Label",
		Items: items,
	})
	wrapper := securityMultiSelectWrapper{model: model}

	view := wrapper.View()
	assert.NotEmpty(t, view)
}

func TestSecurityMultiSelectWrapperUpdateCtrlC(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := securityMultiSelectWrapper{model: model}

	newModel, cmd := wrapper.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestSecurityMultiSelectWrapperUpdateRegularKey(t *testing.T) {
	items := []components.MultiSelectItem{
		{Label: "Test", Value: "test", Selected: false},
	}
	model := components.NewMultiSelect(components.MultiSelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := securityMultiSelectWrapper{model: model}

	newModel, _ := wrapper.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.NotNil(t, newModel)
}

func TestSelectWrapperInit(t *testing.T) {
	items := []components.SelectItem{
		{Label: "Test", Value: "test"},
	}
	model := components.NewSelect(components.SelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := selectWrapper{model: model}

	cmd := wrapper.Init()
	assert.Nil(t, cmd)
}

func TestSelectWrapperView(t *testing.T) {
	items := []components.SelectItem{
		{Label: "Test", Value: "test"},
	}
	model := components.NewSelect(components.SelectConfig{
		Label: "Test Label",
		Items: items,
	})
	wrapper := selectWrapper{model: model}

	view := wrapper.View()
	assert.NotEmpty(t, view)
}

func TestSelectWrapperUpdateCtrlC(t *testing.T) {
	items := []components.SelectItem{
		{Label: "Test", Value: "test"},
	}
	model := components.NewSelect(components.SelectConfig{
		Label: "Test",
		Items: items,
	})
	wrapper := selectWrapper{model: model}

	newModel, cmd := wrapper.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.NotNil(t, newModel)
	assert.NotNil(t, cmd)
}

func TestDetectUserEntityNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-detect-user-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	fs := afero.NewOsFs()
	entityName, _ := detectUserEntity(tmpDir, fs, "com.example.demo")
	assert.Empty(t, entityName)
}

func TestDetectUserEntityFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-detect-user-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	entityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "entity")
	require.NoError(t, os.MkdirAll(entityPath, 0755))

	userFile := filepath.Join(entityPath, "User.java")
	require.NoError(t, os.WriteFile(userFile, []byte("public class User {}"), 0644))

	fs := afero.NewOsFs()
	entityName, _ := detectUserEntity(tmpDir, fs, "com.example.demo")
	assert.Equal(t, "User", entityName)
}

func TestDetectUserEntityAppUser(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-detect-user-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	entityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "model")
	require.NoError(t, os.MkdirAll(entityPath, 0755))

	userFile := filepath.Join(entityPath, "AppUser.java")
	require.NoError(t, os.WriteFile(userFile, []byte("public class AppUser {}"), 0644))

	fs := afero.NewOsFs()
	entityName, _ := detectUserEntity(tmpDir, fs, "com.example.demo")
	assert.Equal(t, "AppUser", entityName)
}

func TestDetectUserEntityNoSourcePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-detect-user-nosrc-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	fs := afero.NewOsFs()
	entityName, _ := detectUserEntity(tmpDir, fs, "com.example.demo")
	assert.Empty(t, entityName)
}

func TestGenerateSecurityNoSourcePath(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-nosrc-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	cfg := securityConfig{
		SecurityTypes:  []SecurityType{SecurityJWT},
		UserEntityName: "User",
	}

	err = generateSecurity(profile, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find src/main/java directory")
}

func TestGenerateSecurityJWTIntegration(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-jwt-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:     "com.example.demo",
		Architecture:    detector.ArchLayered,
		Lombok:          detector.LombokProfile{Detected: true},
		HasValidation:   true,
		ValidationStyle: detector.ValidationJakarta,
		IDType:          "Long",
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityJWT},
		UserEntityName:   "User",
		GenerateEntities: true,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "security")
	assert.DirExists(t, securityPath)

	assert.FileExists(t, filepath.Join(securityPath, "SecurityConfig.java"))
	assert.FileExists(t, filepath.Join(securityPath, "JwtUtil.java"))
	assert.FileExists(t, filepath.Join(securityPath, "JwtAuthenticationFilter.java"))
	assert.FileExists(t, filepath.Join(securityPath, "AuthenticationController.java"))
	assert.FileExists(t, filepath.Join(securityPath, "AuthRequest.java"))
	assert.FileExists(t, filepath.Join(securityPath, "AuthResponse.java"))
	assert.FileExists(t, filepath.Join(securityPath, "RegisterRequest.java"))
	assert.FileExists(t, filepath.Join(securityPath, "RefreshTokenRequest.java"))
	assert.FileExists(t, filepath.Join(securityPath, "CustomUserDetailsService.java"))

	entityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "entity")
	assert.FileExists(t, filepath.Join(entityPath, "User.java"))
	assert.FileExists(t, filepath.Join(entityPath, "Role.java"))

	repoPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "repository")
	assert.FileExists(t, filepath.Join(repoPath, "UserRepository.java"))
	assert.FileExists(t, filepath.Join(repoPath, "RoleRepository.java"))
}

func TestGenerateSecuritySessionIntegration(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-session-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
		Lombok:       detector.LombokProfile{Detected: false},
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecuritySession},
		UserEntityName:   "User",
		GenerateEntities: false,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "security")
	assert.DirExists(t, securityPath)

	assert.FileExists(t, filepath.Join(securityPath, "SecurityConfig.java"))
	assert.FileExists(t, filepath.Join(securityPath, "CustomUserDetailsService.java"))
	assert.FileExists(t, filepath.Join(securityPath, "AuthController.java"))
	assert.FileExists(t, filepath.Join(securityPath, "RegisterRequest.java"))
}

func TestGenerateSecurityOAuth2Integration(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-oauth2-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityOAuth2},
		UserEntityName:   "User",
		GenerateEntities: false,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "security")
	assert.DirExists(t, securityPath)

	assert.FileExists(t, filepath.Join(securityPath, "SecurityConfig.java"))
	assert.FileExists(t, filepath.Join(securityPath, "OAuth2UserService.java"))
	assert.FileExists(t, filepath.Join(securityPath, "OAuth2SuccessHandler.java"))
	assert.FileExists(t, filepath.Join(securityPath, "OAuth2UserPrincipal.java"))
}

func TestGenerateSecuritySkipsExistingFiles(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-skip-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "security")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	existingFile := filepath.Join(srcPath, "SecurityConfig.java")
	require.NoError(t, os.WriteFile(existingFile, []byte("existing content"), 0644))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityJWT},
		UserEntityName:   "User",
		GenerateEntities: false,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	content, err := os.ReadFile(existingFile)
	require.NoError(t, err)
	assert.Equal(t, "existing content", string(content))
}

func TestGenerateSecurityFeatureArchitecture(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-feature-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchFeature,
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityJWT},
		UserEntityName:   "User",
		GenerateEntities: true,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "common", "security")
	assert.DirExists(t, securityPath)
	assert.FileExists(t, filepath.Join(securityPath, "SecurityConfig.java"))

	userPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "user")
	assert.FileExists(t, filepath.Join(userPath, "User.java"))
	assert.FileExists(t, filepath.Join(userPath, "UserRepository.java"))
}

func TestGenerateSecurityHexagonalArchitecture(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-hexagonal-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchHexagonal,
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityJWT},
		UserEntityName:   "User",
		GenerateEntities: true,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "infrastructure", "security")
	assert.DirExists(t, securityPath)

	entityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "domain", "model")
	assert.FileExists(t, filepath.Join(entityPath, "User.java"))

	repoPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "adapter", "persistence")
	assert.FileExists(t, filepath.Join(repoPath, "UserRepository.java"))
}

func TestRunSecurityNoInteractiveWithoutType(t *testing.T) {
	cmd := newSecurityCommand()
	cmd.SetArgs([]string{"--no-interactive", "--package", "com.example.app"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify authentication type")
}

func TestRunSecurityMissingPackage(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-nopkg-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	require.NoError(t, os.Chdir(tmpDir))

	cmd := newSecurityCommand()
	cmd.SetArgs([]string{"--no-interactive", "--jwt"})

	err = cmd.Execute()
	assert.Error(t, err)
}

func TestTemplateInfoStruct(t *testing.T) {
	info := templateInfo{
		template:    "security/jwt/SecurityConfig.java.tmpl",
		fileName:    "SecurityConfig.java",
		conditional: "",
	}

	assert.Equal(t, "security/jwt/SecurityConfig.java.tmpl", info.template)
	assert.Equal(t, "SecurityConfig.java", info.fileName)
	assert.Equal(t, "", info.conditional)
}

func TestTemplateInfoWithConditional(t *testing.T) {
	info := templateInfo{
		template:    "some/template.tmpl",
		fileName:    "SomeFile.java",
		conditional: "HasFeature",
	}

	assert.Equal(t, "HasFeature", info.conditional)
}

func TestCheckSecurityDependenciesNoBuildFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-check-deps-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	fs := afero.NewOsFs()
	_, err = checkSecurityDependencies(tmpDir, fs, []SecurityType{SecurityJWT})
	assert.Error(t, err)
}

func TestCheckSecurityDependenciesWithPom(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-check-deps-pom-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
    </dependencies>
</project>`

	pomPath := filepath.Join(tmpDir, "pom.xml")
	require.NoError(t, os.WriteFile(pomPath, []byte(pomContent), 0644))

	fs := afero.NewOsFs()
	missing, err := checkSecurityDependencies(tmpDir, fs, []SecurityType{SecurityJWT})
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(missing), 4)

	hasSecurityDep := false
	hasJwtDep := false
	for _, dep := range missing {
		if dep.ArtifactId == "spring-boot-starter-security" {
			hasSecurityDep = true
		}
		if dep.ArtifactId == "jjwt-api" {
			hasJwtDep = true
		}
	}
	assert.True(t, hasSecurityDep)
	assert.True(t, hasJwtDep)
}

func TestCheckSecurityDependenciesSecurityExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-check-deps-exists-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-security</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-data-jpa</artifactId>
        </dependency>
    </dependencies>
</project>`

	pomPath := filepath.Join(tmpDir, "pom.xml")
	require.NoError(t, os.WriteFile(pomPath, []byte(pomContent), 0644))

	fs := afero.NewOsFs()
	missing, err := checkSecurityDependencies(tmpDir, fs, []SecurityType{SecuritySession})
	require.NoError(t, err)

	for _, dep := range missing {
		assert.NotEqual(t, "spring-boot-starter-security", dep.ArtifactId)
		assert.NotEqual(t, "spring-boot-starter-data-jpa", dep.ArtifactId)
	}
}

func TestAddSecurityDependencies(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test-add-deps-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>1.0.0</version>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
    </dependencies>
</project>`

	pomPath := filepath.Join(tmpDir, "pom.xml")
	require.NoError(t, os.WriteFile(pomPath, []byte(pomContent), 0644))

	fs := afero.NewOsFs()
	deps := []securityDependency{
		{Name: "Spring Security", GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-security"},
	}

	err = addSecurityDependencies(tmpDir, fs, deps)
	require.NoError(t, err)

	updatedContent, err := os.ReadFile(pomPath)
	require.NoError(t, err)
	assert.Contains(t, string(updatedContent), "spring-boot-starter-security")
}

func TestGenerateSecurityMultipleTypes(t *testing.T) {
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	tmpDir, err := os.MkdirTemp("", "test-security-multi-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo")
	require.NoError(t, os.MkdirAll(srcPath, 0755))

	require.NoError(t, os.Chdir(tmpDir))

	profile := &detector.ProjectProfile{
		BasePackage:  "com.example.demo",
		Architecture: detector.ArchLayered,
	}

	cfg := securityConfig{
		SecurityTypes:    []SecurityType{SecurityJWT, SecurityOAuth2},
		UserEntityName:   "User",
		GenerateEntities: false,
	}

	err = generateSecurity(profile, cfg)
	require.NoError(t, err)

	securityPath := filepath.Join(tmpDir, "src", "main", "java", "com", "example", "demo", "security")

	assert.FileExists(t, filepath.Join(securityPath, "JwtUtil.java"))
	assert.FileExists(t, filepath.Join(securityPath, "OAuth2UserService.java"))
}
