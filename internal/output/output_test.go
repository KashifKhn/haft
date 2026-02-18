package output

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestJSON_ProjectInfo(t *testing.T) {
	info := ProjectInfo{
		Name:              "demo",
		GroupID:           "com.example",
		ArtifactID:        "demo",
		Version:           "0.0.1-SNAPSHOT",
		BuildTool:         "maven",
		BuildFile:         "pom.xml",
		JavaVersion:       "21",
		SpringBootVersion: "3.4.1",
		Dependencies: &DependencyInfo{
			Total:          10,
			SpringStarters: 5,
			SpringLibs:     2,
			TestDeps:       3,
		},
		Features: &FeatureInfo{
			HasWeb:        true,
			HasJPA:        true,
			HasLombok:     true,
			HasValidation: true,
			HasMapStruct:  false,
			HasSecurity:   false,
		},
	}

	output := captureOutput(func() {
		err := JSON(info)
		assert.NoError(t, err)
	})

	var result ProjectInfo
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.Equal(t, "demo", result.Name)
	assert.Equal(t, "com.example", result.GroupID)
	assert.Equal(t, "maven", result.BuildTool)
	assert.Equal(t, 10, result.Dependencies.Total)
	assert.True(t, result.Features.HasWeb)
	assert.True(t, result.Features.HasJPA)
	assert.False(t, result.Features.HasSecurity)
}

func TestJSON_Response_Success(t *testing.T) {
	data := map[string]string{"key": "value"}

	output := captureOutput(func() {
		err := Success(data)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Data)
}

func TestJSON_Response_Error(t *testing.T) {
	output := captureOutput(func() {
		err := Error("NOT_FOUND", "Resource not found", "User with ID 123 not found")
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.False(t, result.Success)
	assert.Nil(t, result.Data)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "NOT_FOUND", result.Error.Code)
	assert.Equal(t, "Resource not found", result.Error.Message)
	assert.Equal(t, "User with ID 123 not found", result.Error.Details)
}

func TestJSON_RoutesOutput(t *testing.T) {
	routes := RoutesOutput{
		Routes: []RouteInfo{
			{Method: "GET", Path: "/api/users", Controller: "UserController", Handler: "getAll", File: "UserController.java", Line: 25},
			{Method: "POST", Path: "/api/users", Controller: "UserController", Handler: "create", File: "UserController.java", Line: 35},
			{Method: "GET", Path: "/api/users/{id}", Controller: "UserController", Handler: "getById", File: "UserController.java", Line: 30},
		},
		Total: 3,
	}

	output := captureOutput(func() {
		err := Success(routes)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)

	dataBytes, _ := json.Marshal(result.Data)
	var routesResult RoutesOutput
	err = json.Unmarshal(dataBytes, &routesResult)
	require.NoError(t, err)

	assert.Equal(t, 3, routesResult.Total)
	assert.Len(t, routesResult.Routes, 3)
	assert.Equal(t, "GET", routesResult.Routes[0].Method)
	assert.Equal(t, "/api/users", routesResult.Routes[0].Path)
}

func TestJSON_StatsOutput(t *testing.T) {
	stats := StatsOutput{
		Languages: []LanguageStats{
			{Name: "Java", Files: 50, Lines: 5000, Code: 4000, Comments: 500, Blanks: 500},
			{Name: "XML", Files: 10, Lines: 500, Code: 400, Comments: 50, Blanks: 50},
		},
		Summary: CodeStatsInfo{
			TotalFiles:  60,
			LinesOfCode: 4400,
			Comments:    550,
			BlankLines:  550,
		},
	}

	output := captureOutput(func() {
		err := Success(stats)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_GenerateOutput(t *testing.T) {
	gen := GenerateOutput{
		Results: []GenerateResult{
			{
				Type:      "resource",
				Name:      "User",
				Generated: []string{"UserController.java", "UserService.java", "User.java"},
				Skipped:   []string{"UserRepository.java"},
			},
		},
		TotalGenerated: 3,
		TotalSkipped:   1,
	}

	output := captureOutput(func() {
		err := Success(gen)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_CatalogOutput(t *testing.T) {
	catalog := CatalogOutput{
		Categories: []CatalogCategory{
			{
				Name: "Web",
				Dependencies: []CatalogItem{
					{Shortcut: "web", Name: "Spring Web", Description: "Build web applications"},
					{Shortcut: "webflux", Name: "Spring WebFlux", Description: "Reactive web applications"},
				},
			},
			{
				Name: "SQL",
				Dependencies: []CatalogItem{
					{Shortcut: "jpa", Name: "Spring Data JPA", Description: "JPA support"},
					{Shortcut: "mysql", Name: "MySQL Driver", Description: "MySQL database driver"},
				},
			},
		},
		Total: 4,
	}

	output := captureOutput(func() {
		err := Success(catalog)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_DependencyList(t *testing.T) {
	info := ProjectInfo{
		Name: "demo",
		Dependencies: &DependencyInfo{
			Total: 3,
			List: []DependencyListItem{
				{GroupID: "org.springframework.boot", ArtifactID: "spring-boot-starter-web"},
				{GroupID: "org.projectlombok", ArtifactID: "lombok", Scope: "provided"},
				{GroupID: "com.h2database", ArtifactID: "h2", Version: "2.2.224", Scope: "test"},
			},
		},
	}

	output := captureOutput(func() {
		err := Success(info)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_AddRemoveResult(t *testing.T) {
	result := AddRemoveResult{
		Action:  "add",
		Added:   []string{"lombok", "validation"},
		Skipped: []string{"jpa"},
	}

	output := captureOutput(func() {
		err := Success(result)
		assert.NoError(t, err)
	})

	var resp Response
	err := json.Unmarshal([]byte(output), &resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
}

func TestJSON_GeneratorTypes(t *testing.T) {
	types := GeneratorTypesOutput{
		Types: []GeneratorType{
			{Name: "resource", Alias: "r", Description: "Generate complete CRUD resource"},
			{Name: "controller", Alias: "co", Description: "Generate REST controller"},
			{Name: "service", Alias: "s", Description: "Generate service interface and implementation"},
			{Name: "security", Alias: "sec", Description: "Generate security configuration"},
		},
	}

	output := captureOutput(func() {
		err := Success(types)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_SecurityTypes(t *testing.T) {
	types := SecurityTypesOutput{
		Types: []SecurityType{
			{Name: "jwt", Description: "Stateless token-based authentication"},
			{Name: "session", Description: "Form-based session authentication"},
			{Name: "oauth2", Description: "Social login (Google, GitHub, Facebook)"},
		},
	}

	output := captureOutput(func() {
		err := Success(types)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_CodeStatsWithCocomo(t *testing.T) {
	stats := CodeStatsInfo{
		TotalFiles:      100,
		LinesOfCode:     10000,
		Comments:        2000,
		BlankLines:      1500,
		TotalBytes:      500000,
		TotalLines:      13500,
		EstimatedCost:   150000.50,
		EstimatedMonths: 12.5,
		EstimatedPeople: 2.5,
	}

	output := captureOutput(func() {
		err := Success(stats)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_EmptyRoutes(t *testing.T) {
	routes := RoutesOutput{
		Routes: []RouteInfo{},
		Total:  0,
	}

	output := captureOutput(func() {
		err := Success(routes)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestJSON_ArchitectureInfo(t *testing.T) {
	info := ProjectInfo{
		Name: "demo",
		Architecture: &ArchitectureInfo{
			Pattern:      "layered",
			FeatureStyle: "flat",
			DTONaming:    "request-response",
			IDType:       "Long",
			MapperType:   "mapstruct",
		},
	}

	output := captureOutput(func() {
		err := Success(info)
		assert.NoError(t, err)
	})

	var result Response
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err)

	assert.True(t, result.Success)
}

func TestPrint_JSONFormat(t *testing.T) {
	data := map[string]string{"test": "value"}
	textCalled := false

	output := captureOutput(func() {
		err := Print(FormatJSON, func() error {
			textCalled = true
			return nil
		}, data)
		assert.NoError(t, err)
	})

	assert.False(t, textCalled)
	assert.Contains(t, output, "success")
	assert.Contains(t, output, "true")
}

func TestPrint_TextFormat(t *testing.T) {
	data := map[string]string{"test": "value"}
	textCalled := false

	output := captureOutput(func() {
		err := Print(FormatText, func() error {
			textCalled = true
			return nil
		}, data)
		assert.NoError(t, err)
	})

	assert.True(t, textCalled)
	assert.Empty(t, output)
}

func TestFormatConstants(t *testing.T) {
	assert.Equal(t, Format("text"), FormatText)
	assert.Equal(t, Format("json"), FormatJSON)
}
