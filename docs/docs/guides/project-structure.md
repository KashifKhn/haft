---
sidebar_position: 3
title: Project Structure
description: Understanding the generated project structure
---

# Project Structure

Haft generates Spring Boot projects following standard Maven/Gradle conventions and Spring best practices.

## Generated Structure

### Maven Project

When you run `haft init my-app` with Maven:

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       └── MyAppApplication.java
│   │   └── resources/
│   │       └── application.yml
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── .gitignore
├── .haft.json
├── mvnw
├── mvnw.cmd
└── pom.xml
```

### Gradle Project

When you run `haft init my-app --build gradle` or `--build gradle-kotlin`:

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/myapp/
│   │   │       └── MyAppApplication.java
│   │   └── resources/
│   │       └── application.yml
│   └── test/
│       └── java/
│           └── com/example/myapp/
│               └── MyAppApplicationTests.java
├── gradle/
│   └── wrapper/
│       └── gradle-wrapper.properties
├── .gitignore
├── .haft.json
├── build.gradle          # or build.gradle.kts for Kotlin DSL
├── settings.gradle       # or settings.gradle.kts for Kotlin DSL
├── gradlew
└── gradlew.bat
```

## File Details

### MyAppApplication.java

The main application class with Spring Boot entry point.

```java
package com.example.myapp;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class MyAppApplication {

    public static void main(String[] args) {
        SpringApplication.run(MyAppApplication.class, args);
    }
}
```

### application.yml

Configuration file (or `application.properties` if you chose Properties format).

```yaml
spring:
  application:
    name: my-app
```

### MyAppApplicationTests.java

Basic application context test.

```java
package com.example.myapp;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest
class MyAppApplicationTests {

    @Test
    void contextLoads() {
    }
}
```

### pom.xml (Maven)

Maven project configuration with:

- Spring Boot parent
- Selected dependencies
- Java version property
- Spring Boot Maven plugin

### build.gradle / build.gradle.kts (Gradle)

Gradle build configuration with:

- Spring Boot plugin
- Selected dependencies
- Java toolchain configuration
- Spring Boot Gradle plugin

### .gitignore

Pre-configured ignore patterns for:

- Build outputs (`target/`, `build/`)
- IDE files (`.idea/`, `.vscode/`, `*.iml`)
- Environment files (`.env`)
- Log files
- OS files (`.DS_Store`, `Thumbs.db`)

### Maven Wrapper

`mvnw` and `mvnw.cmd` allow building without installing Maven:

```bash
# Unix/macOS
./mvnw spring-boot:run

# Windows
mvnw.cmd spring-boot:run
```

### Gradle Wrapper

`gradlew` and `gradlew.bat` allow building without installing Gradle:

```bash
# Unix/macOS
./gradlew bootRun

# Windows
gradlew.bat bootRun
```

## With Resources Generated

After running `haft generate resource User`, the structure expands:

```
my-app/
├── src/
│   └── main/
│       └── java/
│           └── com/example/myapp/
│               ├── MyAppApplication.java
│               ├── controller/
│               │   └── UserController.java
│               ├── dto/
│               │   ├── UserRequest.java
│               │   └── UserResponse.java
│               ├── entity/
│               │   └── User.java
│               ├── exception/
│               │   └── ResourceNotFoundException.java
│               ├── mapper/
│               │   └── UserMapper.java
│               ├── repository/
│               │   └── UserRepository.java
│               └── service/
│                   ├── UserService.java
│                   └── impl/
│                       └── UserServiceImpl.java
```

## Package Organization

Haft follows a **package-by-layer** structure:

| Package | Contents |
|---------|----------|
| `controller/` | REST controllers |
| `service/` | Service interfaces |
| `service/impl/` | Service implementations |
| `repository/` | Spring Data repositories |
| `entity/` | JPA entities |
| `dto/` | Data Transfer Objects |
| `mapper/` | MapStruct mappers |
| `exception/` | Custom exceptions |

## Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Entity | PascalCase, singular | `User`, `OrderItem` |
| Repository | Entity + Repository | `UserRepository` |
| Service | Entity + Service | `UserService` |
| ServiceImpl | Entity + ServiceImpl | `UserServiceImpl` |
| Controller | Entity + Controller | `UserController` |
| Request DTO | Entity + Request | `UserRequest` |
| Response DTO | Entity + Response | `UserResponse` |
| Mapper | Entity + Mapper | `UserMapper` |
| Exception | ResourceNotFoundException | `ResourceNotFoundException` (shared) |

## REST Endpoints

Generated controllers follow RESTful conventions:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | List all |
| `GET` | `/api/users/{id}` | Get by ID |
| `POST` | `/api/users` | Create |
| `PUT` | `/api/users/{id}` | Update |
| `DELETE` | `/api/users/{id}` | Delete |

Endpoint path is the **plural lowercase** form of the entity name.

## Configuration Format

### Properties Format

```properties
spring.application.name=my-app
server.port=8080
```

### YAML Format (Default)

```yaml
spring:
  application:
    name: my-app

server:
  port: 8080
```

YAML is recommended for:
- Better readability
- Hierarchical configuration
- Profile management
