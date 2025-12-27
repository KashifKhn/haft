package routes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"path with leading slash", "/users", "/users"},
		{"path without leading slash", "users", "/users"},
		{"path with whitespace", "  /users  ", "/users"},
		{"path with double quotes", `"/users"`, "/users"},
		{"path with single quotes", `'/users'`, "/users"},
		{"path with id param", "/users/{id}", "/users/{id}"},
		{"nested path", "/api/v1/users", "/api/v1/users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinPaths(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		subPath  string
		expected string
	}{
		{"both empty", "", "", "/"},
		{"base only", "/api", "", "/api"},
		{"sub only", "", "/users", "/users"},
		{"both paths", "/api", "/users", "/api/users"},
		{"base with trailing slash", "/api/", "/users", "/api/users"},
		{"sub without leading slash", "/api", "users", "/api/users"},
		{"nested paths", "/api/v1", "/users/{id}", "/api/v1/users/{id}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		{Method: "PUT", Path: "/users"},
		{Method: "PATCH", Path: "/users"},
	}

	sortRoutes(routes)

	assert.Equal(t, "/api", routes[0].Path)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/api", routes[1].Path)
	assert.Equal(t, "POST", routes[1].Method)
	assert.Equal(t, "/users", routes[2].Path)
	assert.Equal(t, "GET", routes[2].Method)
	assert.Equal(t, "/users", routes[3].Path)
	assert.Equal(t, "PUT", routes[3].Method)
	assert.Equal(t, "/users", routes[4].Path)
	assert.Equal(t, "PATCH", routes[4].Method)
	assert.Equal(t, "/users", routes[5].Path)
	assert.Equal(t, "DELETE", routes[5].Method)
}

func TestParseFileForRoutes_RestController(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class UserController {

    @GetMapping
    public List<User> getAllUsers() {
        return userService.findAll();
    }

    @GetMapping("/{id}")
    public User getUserById(@PathVariable Long id) {
        return userService.findById(id);
    }

    @PostMapping
    public User createUser(@RequestBody User user) {
        return userService.save(user);
    }

    @PutMapping("/{id}")
    public User updateUser(@PathVariable Long id, @RequestBody User user) {
        return userService.update(id, user);
    }

    @DeleteMapping("/{id}")
    public void deleteUser(@PathVariable Long id) {
        userService.delete(id);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "UserController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 5)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "GET /api/users")
	assert.Contains(t, routeMap, "GET /api/users/{id}")
	assert.Contains(t, routeMap, "POST /api/users")
	assert.Contains(t, routeMap, "PUT /api/users/{id}")
	assert.Contains(t, routeMap, "DELETE /api/users/{id}")

	for _, r := range routes {
		assert.Equal(t, "UserController", r.Controller)
		assert.Equal(t, filePath, r.File)
		assert.Greater(t, r.Line, 0)
	}
}

func TestParseFileForRoutes_ControllerAnnotation(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

@Controller
@RequestMapping("/web")
public class WebController {

    @GetMapping("/home")
    public String home() {
        return "home";
    }

    @PostMapping("/submit")
    public String submit(@ModelAttribute Form form) {
        return "result";
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "WebController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "GET /web/home")
	assert.Contains(t, routeMap, "POST /web/submit")
}

func TestParseFileForRoutes_PatchMapping(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/products")
public class ProductController {

    @PatchMapping("/{id}")
    public Product partialUpdate(@PathVariable Long id, @RequestBody Map<String, Object> updates) {
        return productService.partialUpdate(id, updates);
    }

    @PatchMapping("/{id}/status")
    public Product updateStatus(@PathVariable Long id, @RequestParam String status) {
        return productService.updateStatus(id, status);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "ProductController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "PATCH /api/products/{id}")
	assert.Contains(t, routeMap, "PATCH /api/products/{id}/status")
}

func TestParseFileForRoutes_NoClassLevelMapping(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
public class HealthController {

    @GetMapping("/health")
    public String health() {
        return "OK";
    }

    @GetMapping("/ready")
    public String ready() {
        return "READY";
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "HealthController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "GET /health")
	assert.Contains(t, routeMap, "GET /ready")
}

func TestParseFileForRoutes_NonController(t *testing.T) {
	content := `package com.example.demo.service;

import org.springframework.stereotype.Service;

@Service
public class UserService {

    public List<User> findAll() {
        return userRepository.findAll();
    }

    public User findById(Long id) {
        return userRepository.findById(id);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "UserService.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Nil(t, routes)
}

func TestParseFileForRoutes_MappingWithValueAttribute(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/api/orders")
public class OrderController {

    @GetMapping(value = "/{id}")
    public Order getOrder(@PathVariable Long id) {
        return orderService.findById(id);
    }

    @PostMapping(value = "/create")
    public Order createOrder(@RequestBody Order order) {
        return orderService.save(order);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "OrderController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "GET /api/orders/{id}")
	assert.Contains(t, routeMap, "POST /api/orders/create")
}

func TestParseFileForRoutes_HandlerNames(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/items")
public class ItemController {

    @GetMapping
    public List<Item> listItems() {
        return itemService.findAll();
    }

    @GetMapping("/{id}")
    public Item getItemById(@PathVariable Long id) {
        return itemService.findById(id);
    }

    @PostMapping
    public Item addNewItem(@RequestBody Item item) {
        return itemService.save(item);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "ItemController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 3)

	handlerMap := make(map[string]string)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		handlerMap[key] = r.Handler
	}

	assert.Equal(t, "listItems", handlerMap["GET /api/items"])
	assert.Equal(t, "getItemById", handlerMap["GET /api/items/{id}"])
	assert.Equal(t, "addNewItem", handlerMap["POST /api/items"])
}

func TestScanForRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	javaDir := filepath.Join(tmpDir, "src", "main", "java", "com", "example")
	err := os.MkdirAll(javaDir, 0755)
	require.NoError(t, err)

	controller1 := `package com.example;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class UserController {
    @GetMapping
    public List<User> list() { return null; }
}
`
	err = os.WriteFile(filepath.Join(javaDir, "UserController.java"), []byte(controller1), 0644)
	require.NoError(t, err)

	controller2 := `package com.example;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/products")
public class ProductController {
    @GetMapping
    public List<Product> list() { return null; }
    
    @PostMapping
    public Product create(@RequestBody Product p) { return null; }
}
`
	err = os.WriteFile(filepath.Join(javaDir, "ProductController.java"), []byte(controller2), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	routes, err := scanForRoutes(srcDir)
	require.NoError(t, err)

	assert.Len(t, routes, 3)

	routeMap := make(map[string]bool)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = true
	}

	assert.True(t, routeMap["GET /api/users"])
	assert.True(t, routeMap["GET /api/products"])
	assert.True(t, routeMap["POST /api/products"])
}

func TestScanForRoutes_IgnoresNonControllers(t *testing.T) {
	tmpDir := t.TempDir()
	javaDir := filepath.Join(tmpDir, "src", "main", "java", "com", "example")
	err := os.MkdirAll(javaDir, 0755)
	require.NoError(t, err)

	controller := `package com.example;

import org.springframework.web.bind.annotation.*;

@RestController
public class ApiController {
    @GetMapping("/api")
    public String api() { return "ok"; }
}
`
	err = os.WriteFile(filepath.Join(javaDir, "ApiController.java"), []byte(controller), 0644)
	require.NoError(t, err)

	service := `package com.example;

import org.springframework.stereotype.Service;

@Service
public class MyService {
    public void doSomething() {}
}
`
	err = os.WriteFile(filepath.Join(javaDir, "MyService.java"), []byte(service), 0644)
	require.NoError(t, err)

	repository := `package com.example;

import org.springframework.stereotype.Repository;

@Repository
public interface MyRepository {
}
`
	err = os.WriteFile(filepath.Join(javaDir, "MyRepository.java"), []byte(repository), 0644)
	require.NoError(t, err)

	srcDir := filepath.Join(tmpDir, "src", "main", "java")
	routes, err := scanForRoutes(srcDir)
	require.NoError(t, err)

	assert.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/api", routes[0].Path)
}

func TestScanForRoutes_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	javaDir := filepath.Join(tmpDir, "src", "main", "java")
	err := os.MkdirAll(javaDir, 0755)
	require.NoError(t, err)

	routes, err := scanForRoutes(javaDir)
	require.NoError(t, err)
	assert.Len(t, routes, 0)
}

func TestParseFileForRoutes_MultipleAnnotationsOnClass(t *testing.T) {
	content := `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.validation.annotation.Validated;

@RestController
@RequestMapping("/api/v2/accounts")
@Validated
public class AccountController {

    @GetMapping
    public List<Account> getAll() {
        return accountService.findAll();
    }

    @DeleteMapping("/{id}")
    public void delete(@PathVariable Long id) {
        accountService.delete(id);
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "AccountController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	routeMap := make(map[string]Route)
	for _, r := range routes {
		key := r.Method + " " + r.Path
		routeMap[key] = r
	}

	assert.Contains(t, routeMap, "GET /api/v2/accounts")
	assert.Contains(t, routeMap, "DELETE /api/v2/accounts/{id}")
}

func TestParseFileForRoutes_LineNumbers(t *testing.T) {
	content := `package com.example;

import org.springframework.web.bind.annotation.*;

@RestController
public class TestController {

    @GetMapping("/first")
    public String first() {
        return "first";
    }

    @PostMapping("/second")
    public String second() {
        return "second";
    }
}
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "TestController.java")
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(t, err)

	routes, err := parseFileForRoutes(filePath)
	require.NoError(t, err)

	assert.Len(t, routes, 2)

	var getRoute, postRoute Route
	for _, r := range routes {
		if r.Method == "GET" {
			getRoute = r
		} else if r.Method == "POST" {
			postRoute = r
		}
	}

	assert.Equal(t, 8, getRoute.Line)
	assert.Equal(t, 13, postRoute.Line)
}
