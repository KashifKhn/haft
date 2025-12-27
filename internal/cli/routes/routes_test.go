package routes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()

	assert.Equal(t, "routes", cmd.Use)
	assert.Equal(t, "List all REST endpoints", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)

	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)

	filesFlag := cmd.Flags().Lookup("files")
	assert.NotNil(t, filesFlag)
	assert.Equal(t, "f", filesFlag.Shorthand)
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"/users", "/users"},
		{"users", "/users"},
		{"  /users  ", "/users"},
		{`"/users"`, "/users"},
		{`'/users'`, "/users"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := cleanPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinPaths(t *testing.T) {
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

func TestSortRoutes(t *testing.T) {
	routes := []Route{
		{Method: "DELETE", Path: "/users"},
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/api"},
		{Method: "GET", Path: "/api"},
	}

	sortRoutes(routes)

	assert.Equal(t, "/api", routes[0].Path)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/api", routes[1].Path)
	assert.Equal(t, "POST", routes[1].Method)
	assert.Equal(t, "/users", routes[2].Path)
	assert.Equal(t, "GET", routes[2].Method)
	assert.Equal(t, "/users", routes[3].Path)
	assert.Equal(t, "DELETE", routes[3].Method)
}

func TestRouteStruct(t *testing.T) {
	route := Route{
		Method:     "GET",
		Path:       "/api/users",
		Controller: "UserController",
		Handler:    "getUsers",
		File:       "src/main/java/com/example/UserController.java",
		Line:       25,
	}

	assert.Equal(t, "GET", route.Method)
	assert.Equal(t, "/api/users", route.Path)
	assert.Equal(t, "UserController", route.Controller)
	assert.Equal(t, "getUsers", route.Handler)
	assert.Equal(t, 25, route.Line)
}
