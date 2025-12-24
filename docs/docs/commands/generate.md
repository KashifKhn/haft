---
sidebar_position: 2
title: haft generate
description: Generate CRUD resources and components
---

# haft generate

Generate boilerplate code for Spring Boot applications.

## Usage

```bash
haft generate <subcommand> [name] [flags]
haft g <subcommand> [name] [flags]  # alias
```

## Subcommands

### haft generate resource

Generate a complete CRUD resource with all layers.

```bash
# Interactive mode (recommended)
haft generate resource

# With resource name
haft generate resource User
haft g resource Product
```

This generates:

| File | Description |
|------|-------------|
| `controller/UserController.java` | REST controller with CRUD endpoints |
| `service/UserService.java` | Service interface |
| `service/impl/UserServiceImpl.java` | Service implementation |
| `repository/UserRepository.java` | Spring Data JPA repository |
| `entity/User.java` | JPA entity |
| `dto/UserRequest.java` | Request DTO |
| `dto/UserResponse.java` | Response DTO |
| `mapper/UserMapper.java` | Entity to DTO mapper |
| `exception/ResourceNotFoundException.java` | Shared exception (created once) |

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--package` | `-p` | Override base package (auto-detected from pom.xml) |
| `--no-interactive` | | Skip interactive wizard |
| `--skip-entity` | | Skip entity generation |
| `--skip-repository` | | Skip repository generation |

### Examples

```bash
# Interactive mode - wizard guides you through options
haft generate resource

# Non-interactive with name
haft generate resource User --no-interactive

# Override base package
haft generate resource Product --package com.mycompany.store

# Skip database layer (service-only pattern)
haft generate resource Payment --skip-entity --skip-repository
```

## Smart Detection

Haft reads your `pom.xml` to automatically detect and customize generated code:

| Dependency | Detection | Effect |
|------------|-----------|--------|
| **Lombok** | `org.projectlombok:lombok` | Generates `@Getter`, `@Setter`, `@Builder`, etc. |
| **Spring Data JPA** | `spring-boot-starter-data-jpa` | Generates Entity and Repository with `@Transactional` |
| **Validation** | `spring-boot-starter-validation` | Adds `@Valid` to controller parameters |

## Generated Code Examples

### Controller

```java
package com.example.demo.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import jakarta.validation.Valid;
import com.example.demo.service.UserService;
import com.example.demo.dto.UserRequest;
import com.example.demo.dto.UserResponse;

import java.util.List;

@RestController
@RequestMapping("/api/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping
    public ResponseEntity<List<UserResponse>> getAll() {
        return ResponseEntity.ok(userService.findAll());
    }

    @GetMapping("/{id}")
    public ResponseEntity<UserResponse> getById(@PathVariable Long id) {
        return ResponseEntity.ok(userService.findById(id));
    }

    @PostMapping
    public ResponseEntity<UserResponse> create(@Valid @RequestBody UserRequest request) {
        return ResponseEntity.ok(userService.create(request));
    }

    @PutMapping("/{id}")
    public ResponseEntity<UserResponse> update(@PathVariable Long id, @Valid @RequestBody UserRequest request) {
        return ResponseEntity.ok(userService.update(id, request));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable Long id) {
        userService.delete(id);
        return ResponseEntity.noContent().build();
    }
}
```

### Entity (with Lombok)

```java
package com.example.demo.entity;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "users")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
}
```

### Entity (without Lombok)

```java
package com.example.demo.entity;

import jakarta.persistence.*;

@Entity
@Table(name = "users")
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }
}
```

### Service Implementation

```java
package com.example.demo.service.impl;

import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import com.example.demo.service.UserService;
import com.example.demo.dto.UserRequest;
import com.example.demo.dto.UserResponse;
import com.example.demo.mapper.UserMapper;
import com.example.demo.repository.UserRepository;
import com.example.demo.entity.User;
import com.example.demo.exception.ResourceNotFoundException;

import java.util.List;

@Service
@Transactional
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final UserMapper userMapper;

    public UserServiceImpl(UserRepository userRepository, UserMapper userMapper) {
        this.userRepository = userRepository;
        this.userMapper = userMapper;
    }

    @Override
    @Transactional(readOnly = true)
    public List<UserResponse> findAll() {
        return userRepository.findAll().stream()
                .map(userMapper::toResponse)
                .toList();
    }

    @Override
    @Transactional(readOnly = true)
    public UserResponse findById(Long id) {
        return userRepository.findById(id)
                .map(userMapper::toResponse)
                .orElseThrow(() -> new ResourceNotFoundException("User not found with id: " + id));
    }

    @Override
    public UserResponse create(UserRequest request) {
        User user = userMapper.toEntity(request);
        User saved = userRepository.save(user);
        return userMapper.toResponse(saved);
    }

    @Override
    public UserResponse update(Long id, UserRequest request) {
        User user = userRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("User not found with id: " + id));
        userMapper.updateEntity(user, request);
        User updated = userRepository.save(user);
        return userMapper.toResponse(updated);
    }

    @Override
    public void delete(Long id) {
        if (!userRepository.existsById(id)) {
            throw new ResourceNotFoundException("User not found with id: " + id);
        }
        userRepository.deleteById(id);
    }
}
```

## File Safety

Haft never overwrites existing files. If a file already exists, it will be skipped with a warning:

```
WARN File already exists, skipping path=src/main/java/.../UserController.java
```

This allows you to safely re-run the command without losing custom code.

## Coming Soon

Individual component generators are planned for future releases:

```bash
# Generate only specific components (planned)
haft generate controller Product
haft generate service Order
haft generate entity Customer
haft generate repository Invoice
```

## See Also

- [haft init](/docs/commands/init) - Initialize a new project
- [Project Structure](/docs/guides/project-structure) - Where files are generated
