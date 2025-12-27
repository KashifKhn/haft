---
sidebar_position: 8
title: haft routes
description: List REST API endpoints in your project
---

# haft routes

Scan and display all REST API endpoints in your Spring Boot project.

## Usage

```bash
haft routes [flags]
```

## Description

The `routes` command scans your Spring Boot project for REST controller classes and extracts all HTTP endpoint mappings. It provides a quick overview of your API structure without running the application.

## Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON format |

## Examples

```bash
# List all routes
haft routes

# Output as JSON
haft routes --json
```

## How It Works

The routes command:

1. Scans all `.java` files in the `src/main/java` directory
2. Identifies classes annotated with `@RestController` or `@Controller`
3. Extracts the base path from `@RequestMapping` on the class
4. Finds all method-level mappings (`@GetMapping`, `@PostMapping`, etc.)
5. Combines class and method paths to build complete endpoints

## Supported Annotations

### Class-Level
- `@RestController`
- `@Controller`
- `@RequestMapping`

### Method-Level
- `@GetMapping`
- `@PostMapping`
- `@PutMapping`
- `@DeleteMapping`
- `@PatchMapping`
- `@RequestMapping`

## Sample Output

```
  REST API Endpoints
──────────────────────────────────────────────────────────────────────────────────
  Method       Path                              Handler
──────────────────────────────────────────────────────────────────────────────────
  GET          /api/users                        UserController.getAllUsers
  GET          /api/users/{id}                   UserController.getUserById
  POST         /api/users                        UserController.createUser
  PUT          /api/users/{id}                   UserController.updateUser
  DELETE       /api/users/{id}                   UserController.deleteUser
  GET          /api/products                     ProductController.getAllProducts
  GET          /api/products/{id}                ProductController.getProductById
  POST         /api/products                     ProductController.createProduct
  PUT          /api/products/{id}                ProductController.updateProduct
  DELETE       /api/products/{id}                ProductController.deleteProduct
──────────────────────────────────────────────────────────────────────────────────
  Total: 10 endpoints
```

## JSON Output

With `--json` flag:

```json
{
  "routes": [
    {
      "method": "GET",
      "path": "/api/users",
      "handler": "UserController.getAllUsers"
    },
    {
      "method": "GET",
      "path": "/api/users/{id}",
      "handler": "UserController.getUserById"
    },
    {
      "method": "POST",
      "path": "/api/users",
      "handler": "UserController.createUser"
    },
    {
      "method": "PUT",
      "path": "/api/users/{id}",
      "handler": "UserController.updateUser"
    },
    {
      "method": "DELETE",
      "path": "/api/users/{id}",
      "handler": "UserController.deleteUser"
    }
  ],
  "total": 5
}
```

## Method Colors

In terminal output, HTTP methods are color-coded:

| Method | Color |
|--------|-------|
| GET | Green |
| POST | Blue |
| PUT | Yellow |
| DELETE | Red |
| PATCH | Cyan |

## Limitations

- Only scans Java source files (not Kotlin)
- Requires standard Spring MVC annotation patterns
- Does not evaluate SpEL expressions in paths
- Does not detect routes defined programmatically

## Use Cases

- Quick API documentation during development
- Verify endpoint structure before testing
- Review API surface area for security audits
- Generate API documentation snippets

## See Also

- [haft info](/docs/commands/info) - Project information
- [haft stats](/docs/commands/stats) - Code statistics
- [haft generate](/docs/commands/generate) - Generate controllers
