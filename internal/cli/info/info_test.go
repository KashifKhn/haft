package info

import (
	"testing"

	"github.com/KashifKhn/haft/internal/buildtool"
	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "info", cmd.Use)
	assert.Equal(t, "Show project information", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)

	locFlag := cmd.Flags().Lookup("loc")
	assert.NotNil(t, locFlag)
}

func TestCountByPrefix(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-jpa"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-test", Scope: "test"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-devtools"},
		{GroupId: "org.projectlombok", ArtifactId: "lombok"},
		{GroupId: "org.mapstruct", ArtifactId: "mapstruct"},
		{GroupId: "com.mysql", ArtifactId: "mysql-connector-j"},
	}

	tests := []struct {
		name     string
		prefix   string
		expected int
	}{
		{"spring-boot-starter prefix", "spring-boot-starter", 3},
		{"spring-boot prefix", "spring-boot", 4},
		{"lombok prefix", "lombok", 1},
		{"mapstruct prefix", "mapstruct", 1},
		{"nonexistent prefix", "hibernate", 0},
		{"empty prefix", "", 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countByPrefix(deps, tt.prefix)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountByPrefix_EmptyDeps(t *testing.T) {
	deps := []buildtool.Dependency{}
	result := countByPrefix(deps, "spring")
	assert.Equal(t, 0, result)
}

func TestCountByScope(t *testing.T) {
	deps := []buildtool.Dependency{
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-test", Scope: "test"},
		{GroupId: "org.junit.jupiter", ArtifactId: "junit-jupiter", Scope: "test"},
		{GroupId: "org.mockito", ArtifactId: "mockito-core", Scope: "test"},
		{GroupId: "com.h2database", ArtifactId: "h2", Scope: "runtime"},
		{GroupId: "org.projectlombok", ArtifactId: "lombok", Scope: "provided"},
	}

	tests := []struct {
		name     string
		scope    string
		expected int
	}{
		{"test scope", "test", 3},
		{"runtime scope", "runtime", 1},
		{"provided scope", "provided", 1},
		{"compile scope (empty)", "compile", 0},
		{"empty scope (no scope set)", "", 1},
		{"case insensitive TEST", "TEST", 3},
		{"case insensitive Test", "Test", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countByScope(deps, tt.scope)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCountByScope_EmptyDeps(t *testing.T) {
	deps := []buildtool.Dependency{}
	result := countByScope(deps, "test")
	assert.Equal(t, 0, result)
}

func TestCountByPrefix_MultiplePrefixMatches(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "spring-boot-starter-web"},
		{ArtifactId: "spring-boot-starter-data-jpa"},
		{ArtifactId: "spring-boot-starter-security"},
		{ArtifactId: "spring-boot-starter-validation"},
		{ArtifactId: "spring-boot-starter-actuator"},
		{ArtifactId: "spring-boot-devtools"},
		{ArtifactId: "spring-security-test"},
		{ArtifactId: "spring-data-commons"},
	}

	starterCount := countByPrefix(deps, "spring-boot-starter")
	assert.Equal(t, 5, starterCount)

	springCount := countByPrefix(deps, "spring")
	assert.Equal(t, 8, springCount)

	springBootCount := countByPrefix(deps, "spring-boot")
	assert.Equal(t, 6, springBootCount)
}

func TestCountByScope_AllTestDeps(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "spring-boot-starter-test", Scope: "test"},
		{ArtifactId: "junit-jupiter", Scope: "test"},
		{ArtifactId: "mockito-core", Scope: "test"},
		{ArtifactId: "assertj-core", Scope: "test"},
		{ArtifactId: "h2", Scope: "test"},
	}

	result := countByScope(deps, "test")
	assert.Equal(t, 5, result)
}

func TestCountByScope_MixedCase(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "dep1", Scope: "test"},
		{ArtifactId: "dep2", Scope: "Test"},
		{ArtifactId: "dep3", Scope: "TEST"},
		{ArtifactId: "dep4", Scope: "tEsT"},
	}

	result := countByScope(deps, "test")
	assert.Equal(t, 4, result)
}

func TestCountByPrefix_PartialMatch(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "spring-boot-starter-web"},
		{ArtifactId: "springdoc-openapi-ui"},
		{ArtifactId: "spring-security-core"},
	}

	result := countByPrefix(deps, "spring-boot")
	assert.Equal(t, 1, result)

	result = countByPrefix(deps, "spring")
	assert.Equal(t, 3, result)
}

func TestCountByPrefix_ExactMatch(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "lombok"},
		{ArtifactId: "lombok-mapstruct-binding"},
	}

	result := countByPrefix(deps, "lombok")
	assert.Equal(t, 2, result)
}

func TestCountByScope_RuntimeDeps(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "mysql-connector-j", Scope: "runtime"},
		{ArtifactId: "postgresql", Scope: "runtime"},
		{ArtifactId: "h2", Scope: "test"},
		{ArtifactId: "spring-boot-devtools", Scope: "runtime"},
	}

	result := countByScope(deps, "runtime")
	assert.Equal(t, 3, result)
}

func TestCountByScope_ProvidedDeps(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "lombok", Scope: "provided"},
		{ArtifactId: "servlet-api", Scope: "provided"},
		{ArtifactId: "spring-boot-starter-web"},
	}

	result := countByScope(deps, "provided")
	assert.Equal(t, 2, result)
}

func TestCountByPrefix_NoMatch(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "spring-boot-starter-web"},
		{ArtifactId: "lombok"},
		{ArtifactId: "mysql-connector-j"},
	}

	result := countByPrefix(deps, "hibernate")
	assert.Equal(t, 0, result)

	result = countByPrefix(deps, "jakarta")
	assert.Equal(t, 0, result)
}

func TestCountByScope_NoMatch(t *testing.T) {
	deps := []buildtool.Dependency{
		{ArtifactId: "spring-boot-starter-web"},
		{ArtifactId: "lombok", Scope: "provided"},
	}

	result := countByScope(deps, "system")
	assert.Equal(t, 0, result)
}

func TestCountByPrefix_LargeDependencyList(t *testing.T) {
	var deps []buildtool.Dependency
	for i := 0; i < 100; i++ {
		deps = append(deps, buildtool.Dependency{
			ArtifactId: "spring-boot-starter-custom-" + string(rune('a'+i%26)),
		})
	}
	for i := 0; i < 50; i++ {
		deps = append(deps, buildtool.Dependency{
			ArtifactId: "other-lib-" + string(rune('a'+i%26)),
		})
	}

	result := countByPrefix(deps, "spring-boot-starter")
	assert.Equal(t, 100, result)

	result = countByPrefix(deps, "other-lib")
	assert.Equal(t, 50, result)
}

func TestCountByScope_LargeDependencyList(t *testing.T) {
	var deps []buildtool.Dependency
	for i := 0; i < 50; i++ {
		deps = append(deps, buildtool.Dependency{
			ArtifactId: "test-lib-" + string(rune('a'+i%26)),
			Scope:      "test",
		})
	}
	for i := 0; i < 30; i++ {
		deps = append(deps, buildtool.Dependency{
			ArtifactId: "runtime-lib-" + string(rune('a'+i%26)),
			Scope:      "runtime",
		})
	}
	for i := 0; i < 20; i++ {
		deps = append(deps, buildtool.Dependency{
			ArtifactId: "compile-lib-" + string(rune('a'+i%26)),
		})
	}

	result := countByScope(deps, "test")
	assert.Equal(t, 50, result)

	result = countByScope(deps, "runtime")
	assert.Equal(t, 30, result)

	result = countByScope(deps, "")
	assert.Equal(t, 20, result)
}

func TestNewCommand_HasDepsFlag(t *testing.T) {
	cmd := NewCommand()

	depsFlag := cmd.Flags().Lookup("deps")
	assert.NotNil(t, depsFlag)
	assert.Equal(t, "false", depsFlag.DefValue)
}

func TestNewCommand_FlagDefaults(t *testing.T) {
	cmd := NewCommand()

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue)

	locFlag := cmd.Flags().Lookup("loc")
	assert.NotNil(t, locFlag)
	assert.Equal(t, "false", locFlag.DefValue)

	depsFlag := cmd.Flags().Lookup("deps")
	assert.NotNil(t, depsFlag)
	assert.Equal(t, "false", depsFlag.DefValue)
}

func TestNewCommand_Examples(t *testing.T) {
	cmd := NewCommand()

	assert.Contains(t, cmd.Example, "haft info")
	assert.Contains(t, cmd.Example, "--json")
	assert.Contains(t, cmd.Example, "--loc")
	assert.Contains(t, cmd.Example, "--deps")
}
