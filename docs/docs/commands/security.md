---
sidebar_position: 10
title: haft generate security
description: Generate security configuration for Spring Boot applications
---

# haft generate security

Generate Spring Security configuration with support for JWT, Session-based, and OAuth2 authentication.

## Usage

```bash
haft generate security [flags]
haft g sec [flags]  # alias
```

## Overview

The security generator creates a complete, production-ready security setup for your Spring Boot application. It supports three authentication types:

| Type | Description | Use Case |
|------|-------------|----------|
| **JWT** | Stateless token-based authentication | REST APIs, microservices, SPAs |
| **Session** | Traditional session-based with form login | Server-rendered web apps, MVC |
| **OAuth2** | Social login (Google, GitHub, Facebook) | Apps with third-party authentication |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--jwt` | | Generate JWT authentication |
| `--session` | | Generate session-based authentication |
| `--oauth2` | | Generate OAuth2 authentication |
| `--all` | | Generate all authentication types |
| `--package` | `-p` | Override base package |
| `--no-interactive` | | Skip interactive wizard |
| `--skip-entities` | | Skip User/Role entity generation |
| `--refresh` | | Force re-scan project profile |

## Examples

```bash
# Interactive mode - select authentication types via TUI
haft generate security

# Generate JWT authentication only
haft generate security --jwt

# Generate session-based authentication
haft generate security --session

# Generate OAuth2 (Google, GitHub, Facebook)
haft generate security --oauth2

# Generate all authentication types
haft generate security --all

# Non-interactive with specific package
haft generate security --jwt --package com.example.app --no-interactive

# Skip User/Role entity generation (use existing)
haft generate security --jwt --skip-entities
```

---

## JWT Authentication

Generates stateless token-based authentication suitable for REST APIs and SPAs.

### Generated Files

| File | Description |
|------|-------------|
| `SecurityConfig.java` | Spring Security configuration with JWT filter chain |
| `JwtUtil.java` | JWT token generation, validation, and extraction |
| `JwtAuthenticationFilter.java` | Request filter for token validation |
| `AuthenticationController.java` | Login, register, and refresh token endpoints |
| `AuthRequest.java` | Login request DTO |
| `AuthResponse.java` | Token response DTO |
| `RegisterRequest.java` | Registration request DTO |
| `RefreshTokenRequest.java` | Refresh token request DTO |
| `CustomUserDetailsService.java` | Loads users from database |
| `User.java` | User entity (optional) |
| `Role.java` | Role entity (optional) |
| `UserRepository.java` | User repository (optional) |
| `RoleRepository.java` | Role repository (optional) |

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login and get tokens |
| POST | `/api/auth/refresh` | Refresh access token |

### Configuration

Add to `application.properties` or `application.yml`:

```yaml
# JWT Configuration
jwt:
  secret: your-256-bit-secret-key-here
  expiration: 86400000  # 24 hours in milliseconds
  refresh-expiration: 604800000  # 7 days in milliseconds
```

### Example Usage

```bash
# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com", "password": "secret123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "secret123"}'

# Access protected resource
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9..."

# Refresh token
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken": "eyJhbGciOiJIUzI1NiJ9..."}'
```

### Required Dependencies

The generator automatically checks and offers to add:

| Dependency | Purpose |
|------------|---------|
| `spring-boot-starter-security` | Spring Security core |
| `spring-boot-starter-data-jpa` | Database access |
| `jjwt-api` (0.12.6) | JWT API |
| `jjwt-impl` (0.12.6) | JWT implementation |
| `jjwt-jackson` (0.12.6) | JWT JSON support |

---

## Session-Based Authentication

Generates traditional form-based authentication suitable for server-rendered web applications.

### Generated Files

| File | Description |
|------|-------------|
| `SecurityConfig.java` | Spring Security with form login, remember-me, session management |
| `CustomUserDetailsService.java` | Loads users from database |
| `AuthController.java` | MVC controller for login/register pages |
| `RegisterRequest.java` | Registration form DTO |

### Features

- Form-based login (`/login`)
- Remember-me functionality
- Session fixation protection
- CSRF protection
- Logout handling

### Configuration

```yaml
# Session Configuration
server:
  servlet:
    session:
      timeout: 30m
      cookie:
        http-only: true
        secure: true  # Enable in production

spring:
  security:
    remember-me:
      key: your-remember-me-key
      token-validity-seconds: 604800  # 7 days
```

### Example Routes

| Route | Description |
|-------|-------------|
| `/login` | Login page |
| `/register` | Registration page |
| `/logout` | Logout (POST) |
| `/dashboard` | Protected page (example) |

---

## OAuth2 Authentication

Generates social login configuration supporting Google, GitHub, and Facebook.

### Generated Files

| File | Description |
|------|-------------|
| `SecurityConfig.java` | OAuth2 login configuration |
| `OAuth2UserService.java` | Custom OAuth2 user handling |
| `OAuth2SuccessHandler.java` | Post-login success handler |
| `OAuth2UserPrincipal.java` | OAuth2User implementation |

### Supported Providers

| Provider | Registration ID |
|----------|-----------------|
| Google | `google` |
| GitHub | `github` |
| Facebook | `facebook` |

### Configuration

Add to `application.yml`:

```yaml
spring:
  security:
    oauth2:
      client:
        registration:
          google:
            client-id: your-google-client-id
            client-secret: your-google-client-secret
            scope:
              - email
              - profile
          github:
            client-id: your-github-client-id
            client-secret: your-github-client-secret
            scope:
              - user:email
              - read:user
          facebook:
            client-id: your-facebook-client-id
            client-secret: your-facebook-client-secret
            scope:
              - email
              - public_profile
```

### OAuth2 Flow

1. User clicks "Login with Google/GitHub/Facebook"
2. Redirected to provider's authorization page
3. User grants permissions
4. Redirected back with authorization code
5. `OAuth2UserService` processes user info
6. `OAuth2SuccessHandler` handles post-login logic (create/update user, generate JWT, etc.)

### Required Dependencies

| Dependency | Purpose |
|------------|---------|
| `spring-boot-starter-security` | Spring Security core |
| `spring-boot-starter-oauth2-client` | OAuth2 client support |

---

## Intelligent Features

### Dependency Checking

The generator automatically:

1. Scans your `pom.xml` or `build.gradle` for existing dependencies
2. Identifies missing required dependencies
3. Prompts to add them automatically
4. Uses your project's build tool (Maven/Gradle)

```
? Missing dependencies detected:
  - spring-boot-starter-security
  - jjwt-api
  - jjwt-impl
  - jjwt-jackson
  
  Add missing dependencies? [Y/n]
```

### User Entity Detection

Before generating User/Role entities, Haft scans your project for existing user-related entities:

| Scanned Directories | Entity Names Checked |
|--------------------|--------------------|
| `entity/` | `User`, `AppUser`, `Account` |
| `model/` | `Member`, `Principal` |
| `domain/` | `UserEntity`, `ApplicationUser` |
| `user/` | |
| `auth/` | |

If found, the generator:
- Skips entity generation
- Uses your existing entity in generated code
- Prompts if you want to generate anyway

### Architecture-Aware Generation

Files are placed according to your project's architecture:

| Architecture | Package Location |
|--------------|------------------|
| **Layered** | `com.example.security` |
| **Feature** | `com.example.auth` or `com.example.security` |
| **Hexagonal** | `com.example.infrastructure.security` |
| **Clean** | `com.example.infrastructure.security` |
| **Modular** | `com.example.security` |

---

## Multiple Authentication Types

You can generate multiple authentication types in the same project:

```bash
# Generate JWT for API + OAuth2 for web
haft generate security --jwt --oauth2

# Generate all types
haft generate security --all
```

When generating multiple types:
- Shared files (like `SecurityConfig.java`) are merged appropriately
- Duplicate files are skipped with a warning
- Each type's specific files are generated

---

## Next Steps After Generation

### JWT Authentication

1. Add `jwt.secret` to `application.properties/yml`
2. Configure `jwt.expiration` (default: 24 hours)
3. Create initial admin user or use `/api/auth/register`
4. Test with `/api/auth/login` endpoint

### Session Authentication

1. Configure session timeout in `application.properties`
2. Create login/register Thymeleaf templates
3. Add CSRF token to forms
4. Configure remember-me key

### OAuth2 Authentication

1. Create OAuth apps at provider consoles:
   - [Google Cloud Console](https://console.cloud.google.com/)
   - [GitHub Developer Settings](https://github.com/settings/developers)
   - [Facebook for Developers](https://developers.facebook.com/)
2. Add client IDs and secrets to `application.yml`
3. Configure redirect URIs at providers

---

## File Safety

Haft never overwrites existing files. If a file already exists, it will be skipped:

```
WARN  File exists, skipping file=SecurityConfig.java
```

Use `--refresh` to force re-scan project profile if your project structure has changed.

## See Also

- [haft generate](/docs/commands/generate) - All generation commands
- [haft add](/docs/commands/add) - Add dependencies
- [Project Structure](/docs/guides/project-structure) - File placement
- [Templates Reference](/docs/reference/templates) - Template customization
