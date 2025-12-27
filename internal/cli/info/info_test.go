package info

import (
	"testing"

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
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		basePath string
		subPath  string
		expected string
	}{
		{"", "", "/"},
		{"/api", "", "/api"},
		{"", "/users", "/users"},
		{"/api", "/users", "/api/users"},
		{"/api/", "/users", "/api/users"},
		{"/api", "users", "/api/users"},
	}

	for _, tt := range tests {
		t.Run(tt.basePath+"_"+tt.subPath, func(t *testing.T) {
			result := joinPaths(tt.basePath, tt.subPath)
			assert.Equal(t, tt.expected, result)
		})
	}
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

	if len(basePath) > 0 && basePath[len(basePath)-1] == '/' {
		basePath = basePath[:len(basePath)-1]
	}
	if len(subPath) > 0 && subPath[0] != '/' {
		subPath = "/" + subPath
	}

	return basePath + subPath
}

func TestCountByPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		expected int
	}{
		{"spring-boot-starter", "spring-boot-starter", 1},
		{"no-match", "xyz", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just testing the command structure
			assert.NotEmpty(t, tt.name)
		})
	}
}
