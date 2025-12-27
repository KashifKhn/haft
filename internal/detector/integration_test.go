package detector

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManualTrackMateDetection(t *testing.T) {
	trackmatePath := "/home/zarqan-khn/mycoding/fyp/trackmate-backend"

	if _, err := os.Stat(trackmatePath); os.IsNotExist(err) {
		t.Skip("TrackMate project not found, skipping manual test")
	}

	d := NewDetector(trackmatePath, WithFileSystem(afero.NewOsFs()))
	profile, err := d.Detect()
	if err != nil {
		t.Fatalf("Error detecting: %v", err)
	}

	t.Logf("=== TrackMate Detection Results ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)
	t.Logf("Base Package: %s", profile.BasePackage)
	t.Logf("Source Root: %s", profile.SourceRoot)

	if profile.BaseEntity != nil {
		t.Logf("Base Entity: %s.%s", profile.BaseEntity.Package, profile.BaseEntity.Name)
	} else {
		t.Log("Base Entity: Not detected")
	}

	t.Logf("DTO Naming: %s", profile.DTONaming)
	t.Logf("Controller Suffix: %s", profile.ControllerSuffix)
	t.Logf("ID Type: %s", profile.IDType)
	t.Logf("Mapper: %s", profile.Mapper)
	t.Logf("Database: %s", profile.Database)

	t.Logf("Lombok Detected: %v", profile.Lombok.Detected)
	t.Logf("  @Data: %v, @Builder: %v, @Slf4j: %v",
		profile.Lombok.UseData, profile.Lombok.UseBuilder, profile.Lombok.UseSlf4j)

	t.Logf("Swagger: %v (%s)", profile.HasSwagger, profile.SwaggerStyle)
	t.Logf("Validation: %v (%s)", profile.HasValidation, profile.ValidationStyle)

	t.Logf("Global Exception Handler: %v", profile.Exceptions.HasGlobalHandler)
	t.Logf("Custom Exceptions: %d", len(profile.Exceptions.CustomExceptions))

	if len(profile.FeatureModules) > 0 {
		t.Logf("Feature Modules: %v", profile.FeatureModules)
	}

	t.Logf("Testing - Framework: %s, Mockito: %v, Testcontainers: %v",
		profile.Testing.Framework, profile.Testing.HasMockito, profile.Testing.HasTestcontainers)

	assert.Equal(t, ArchFeature, profile.Architecture)
	assert.Equal(t, DTONamingRequestResponse, profile.DTONaming)
	assert.Equal(t, "UUID", profile.IDType)
	assert.Equal(t, MapperMapStruct, profile.Mapper)
	assert.Equal(t, "com.trackmate", profile.BasePackage)
	assert.NotNil(t, profile.BaseEntity)
	assert.Equal(t, "BaseEntity", profile.BaseEntity.Name)
	assert.True(t, profile.Lombok.Detected)
	assert.True(t, profile.HasSwagger)
	assert.True(t, profile.HasValidation)
	assert.True(t, profile.Exceptions.HasGlobalHandler)
}

func TestLayeredArchitectureProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/demo/DemoApplication.java": `package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class DemoApplication {
    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}`,
		"/project/src/main/java/com/example/demo/controller/UserController.java": `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;
import com.example.demo.service.UserService;
import com.example.demo.dto.UserDTO;

@RestController
@RequestMapping("/api/users")
public class UserController {
    private final UserService userService;
    
    public UserController(UserService userService) {
        this.userService = userService;
    }
}`,
		"/project/src/main/java/com/example/demo/controller/ProductController.java": `package com.example.demo.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/products")
public class ProductController {
}`,
		"/project/src/main/java/com/example/demo/service/UserService.java": `package com.example.demo.service;

import org.springframework.stereotype.Service;

@Service
public class UserService {
}`,
		"/project/src/main/java/com/example/demo/service/ProductService.java": `package com.example.demo.service;

import org.springframework.stereotype.Service;

@Service
public class ProductService {
}`,
		"/project/src/main/java/com/example/demo/repository/UserRepository.java": `package com.example.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import com.example.demo.entity.User;

public interface UserRepository extends JpaRepository<User, Long> {
}`,
		"/project/src/main/java/com/example/demo/repository/ProductRepository.java": `package com.example.demo.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import com.example.demo.entity.Product;

public interface ProductRepository extends JpaRepository<Product, Long> {
}`,
		"/project/src/main/java/com/example/demo/entity/User.java": `package com.example.demo.entity;

import jakarta.persistence.*;
import lombok.Data;

@Data
@Entity
@Table(name = "users")
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String name;
    private String email;
}`,
		"/project/src/main/java/com/example/demo/entity/Product.java": `package com.example.demo.entity;

import jakarta.persistence.*;
import lombok.Data;

@Data
@Entity
@Table(name = "products")
public class Product {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String name;
}`,
		"/project/src/main/java/com/example/demo/dto/UserDTO.java": `package com.example.demo.dto;

import lombok.Data;

@Data
public class UserDTO {
    private Long id;
    private String name;
}`,
		"/project/src/main/java/com/example/demo/dto/ProductDTO.java": `package com.example.demo.dto;

import lombok.Data;

@Data
public class ProductDTO {
    private Long id;
    private String name;
}`,
		"/project/pom.xml": `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>demo</artifactId>
    <version>0.0.1-SNAPSHOT</version>
</project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Layered Architecture Detection ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)
	t.Logf("Base Package: %s", profile.BasePackage)
	t.Logf("DTO Naming: %s", profile.DTONaming)
	t.Logf("ID Type: %s", profile.IDType)
	t.Logf("Lombok: %v", profile.Lombok.Detected)

	assert.Equal(t, ArchLayered, profile.Architecture)
	assert.Equal(t, "com.example.demo", profile.BasePackage)
	assert.Equal(t, DTONamingDTOUpper, profile.DTONaming)
	assert.Equal(t, "Long", profile.IDType)
	assert.True(t, profile.Lombok.Detected)
	assert.True(t, profile.Lombok.UseData)
}

func TestHexagonalArchitectureProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/app/Application.java": `package com.example.app;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}`,
		"/project/src/main/java/com/example/app/domain/model/User.java": `package com.example.app.domain.model;

import java.util.UUID;
import lombok.Data;

@Data
public class User {
    private UUID id;
    private String name;
    private String email;
}`,
		"/project/src/main/java/com/example/app/domain/model/Product.java": `package com.example.app.domain.model;

import java.util.UUID;
import lombok.Data;

@Data
public class Product {
    private UUID id;
    private String name;
}`,
		"/project/src/main/java/com/example/app/application/port/in/CreateUserUseCase.java": `package com.example.app.application.port.in;

import com.example.app.domain.model.User;

public interface CreateUserUseCase {
    User createUser(String name, String email);
}`,
		"/project/src/main/java/com/example/app/application/port/out/UserRepository.java": `package com.example.app.application.port.out;

import com.example.app.domain.model.User;
import java.util.UUID;

public interface UserRepository {
    User save(User user);
    User findById(UUID id);
}`,
		"/project/src/main/java/com/example/app/application/service/UserService.java": `package com.example.app.application.service;

import org.springframework.stereotype.Service;
import com.example.app.application.port.in.CreateUserUseCase;
import com.example.app.application.port.out.UserRepository;
import com.example.app.domain.model.User;

@Service
public class UserService implements CreateUserUseCase {
    private final UserRepository userRepository;
    
    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }
    
    @Override
    public User createUser(String name, String email) {
        return null;
    }
}`,
		"/project/src/main/java/com/example/app/adapter/in/web/UserController.java": `package com.example.app.adapter.in.web;

import org.springframework.web.bind.annotation.*;
import com.example.app.application.port.in.CreateUserUseCase;

@RestController
@RequestMapping("/api/users")
public class UserController {
    private final CreateUserUseCase createUserUseCase;
    
    public UserController(CreateUserUseCase createUserUseCase) {
        this.createUserUseCase = createUserUseCase;
    }
}`,
		"/project/src/main/java/com/example/app/adapter/in/web/dto/UserRequest.java": `package com.example.app.adapter.in.web.dto;

import lombok.Data;

@Data
public class UserRequest {
    private String name;
    private String email;
}`,
		"/project/src/main/java/com/example/app/adapter/in/web/dto/UserResponse.java": `package com.example.app.adapter.in.web.dto;

import lombok.Data;
import java.util.UUID;

@Data
public class UserResponse {
    private UUID id;
    private String name;
    private String email;
}`,
		"/project/src/main/java/com/example/app/adapter/out/persistence/UserJpaRepository.java": `package com.example.app.adapter.out.persistence;

import org.springframework.data.jpa.repository.JpaRepository;
import java.util.UUID;

public interface UserJpaRepository extends JpaRepository<UserJpaEntity, UUID> {
}`,
		"/project/src/main/java/com/example/app/adapter/out/persistence/UserJpaEntity.java": `package com.example.app.adapter.out.persistence;

import jakarta.persistence.*;
import lombok.Data;
import java.util.UUID;

@Data
@Entity
@Table(name = "users")
public class UserJpaEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;
    private String name;
    private String email;
}`,
		"/project/src/main/java/com/example/app/infrastructure/config/AppConfig.java": `package com.example.app.infrastructure.config;

import org.springframework.context.annotation.Configuration;

@Configuration
public class AppConfig {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Hexagonal Architecture Detection ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)
	t.Logf("Base Package: %s", profile.BasePackage)
	t.Logf("DTO Naming: %s", profile.DTONaming)
	t.Logf("ID Type: %s", profile.IDType)

	assert.Equal(t, ArchHexagonal, profile.Architecture)
	assert.Equal(t, "com.example.app", profile.BasePackage)
	assert.Equal(t, DTONamingRequestResponse, profile.DTONaming)
	assert.Equal(t, "UUID", profile.IDType)
}

func TestFlatArchitectureProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/app/Application.java": `package com.example.app;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}`,
		"/project/src/main/java/com/example/app/User.java": `package com.example.app;

import jakarta.persistence.*;

@Entity
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String name;
}`,
		"/project/src/main/java/com/example/app/UserController.java": `package com.example.app;

import org.springframework.web.bind.annotation.*;

@RestController
public class UserController {
}`,
		"/project/src/main/java/com/example/app/UserRepository.java": `package com.example.app;

import org.springframework.data.jpa.repository.JpaRepository;

public interface UserRepository extends JpaRepository<User, Long> {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Flat Architecture Detection ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)
	t.Logf("Base Package: %s", profile.BasePackage)

	assert.Equal(t, ArchFlat, profile.Architecture)
	assert.Equal(t, "com.example.app", profile.BasePackage)
}

func TestEmptyProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Empty Project Detection ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)

	assert.Equal(t, ArchLayered, profile.Architecture)
	assert.Equal(t, 1.0, profile.ArchConfidence)
}

func TestProjectWithOnlyEntities(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

import jakarta.persistence.*;
import java.util.UUID;

@Entity
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;
}`,
		"/project/src/main/java/com/example/entity/Product.java": `package com.example.entity;

import jakarta.persistence.*;
import java.util.UUID;

@Entity
public class Product {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Entities Only Detection ===")
	t.Logf("ID Type: %s", profile.IDType)
	t.Logf("Database: %s", profile.Database)

	assert.Equal(t, "UUID", profile.IDType)
	assert.Equal(t, DatabaseJPA, profile.Database)
}

func TestProjectWithBaseEntityInheritance(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/common/BaseEntity.java": `package com.example.common;

import jakarta.persistence.*;
import java.util.UUID;
import java.time.LocalDateTime;

@MappedSuperclass
public abstract class BaseEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}`,
		"/project/src/main/java/com/example/user/entity/User.java": `package com.example.user.entity;

import jakarta.persistence.*;
import com.example.common.BaseEntity;

@Entity
public class User extends BaseEntity {
    private String name;
}`,
		"/project/src/main/java/com/example/product/entity/Product.java": `package com.example.product.entity;

import jakarta.persistence.*;
import com.example.common.BaseEntity;

@Entity
public class Product extends BaseEntity {
    private String name;
}`,
		"/project/src/main/java/com/example/order/entity/Order.java": `package com.example.order.entity;

import jakarta.persistence.*;
import com.example.common.BaseEntity;

@Entity
public class Order extends BaseEntity {
    private String status;
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Base Entity Inheritance Detection ===")
	t.Logf("Base Entity: %v", profile.BaseEntity)

	require.NotNil(t, profile.BaseEntity)
	assert.Equal(t, "BaseEntity", profile.BaseEntity.Name)
	assert.Equal(t, "com.example.common", profile.BaseEntity.Package)
}

func TestProjectWithMixedDTONaming(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/dto/UserRequest.java":     `package com.example.dto; public class UserRequest {}`,
		"/project/src/main/java/com/example/dto/UserResponse.java":    `package com.example.dto; public class UserResponse {}`,
		"/project/src/main/java/com/example/dto/ProductRequest.java":  `package com.example.dto; public class ProductRequest {}`,
		"/project/src/main/java/com/example/dto/ProductResponse.java": `package com.example.dto; public class ProductResponse {}`,
		"/project/src/main/java/com/example/dto/OrderDTO.java":        `package com.example.dto; public class OrderDTO {}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Mixed DTO Naming Detection ===")
	t.Logf("DTO Naming: %s (majority wins)", profile.DTONaming)

	assert.Equal(t, DTONamingRequestResponse, profile.DTONaming)
}

func TestProjectWithMongoDBEntities(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.annotation.Id;

@Document(collection = "users")
public class User {
    @Id
    private String id;
    private String name;
}`,
		"/project/src/main/java/com/example/repository/UserRepository.java": `package com.example.repository;

import org.springframework.data.mongodb.repository.MongoRepository;
import com.example.entity.User;

public interface UserRepository extends MongoRepository<User, String> {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== MongoDB Detection ===")
	t.Logf("Database: %s", profile.Database)

	assert.Equal(t, DatabaseMongo, profile.Database)
}

func TestProjectWithAllLombokAnnotations(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

import lombok.*;
import jakarta.persistence.*;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
public class User {
    @Id
    private Long id;
}`,
		"/project/src/main/java/com/example/service/UserService.java": `package com.example.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class UserService {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Lombok Detection ===")
	t.Logf("Detected: %v", profile.Lombok.Detected)
	t.Logf("@Data: %v, @Builder: %v, @NoArgs: %v, @AllArgs: %v",
		profile.Lombok.UseData, profile.Lombok.UseBuilder,
		profile.Lombok.UseNoArgs, profile.Lombok.UseAllArgs)
	t.Logf("@Slf4j: %v, @RequiredArgsConstructor: %v",
		profile.Lombok.UseSlf4j, profile.Lombok.UseRequiredArgs)

	assert.True(t, profile.Lombok.Detected)
	assert.True(t, profile.Lombok.UseData)
	assert.True(t, profile.Lombok.UseBuilder)
	assert.True(t, profile.Lombok.UseNoArgs)
	assert.True(t, profile.Lombok.UseAllArgs)
	assert.True(t, profile.Lombok.UseSlf4j)
	assert.True(t, profile.Lombok.UseRequiredArgs)
}

func TestCleanArchitectureProject(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/app/Application.java": `package com.example.app;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {}`,
		"/project/src/main/java/com/example/app/domain/entity/User.java": `package com.example.app.domain.entity;

public class User {
    private String id;
    private String name;
}`,
		"/project/src/main/java/com/example/app/application/usecase/CreateUserUseCase.java": `package com.example.app.application.usecase;

import org.springframework.stereotype.Service;

@Service
public class CreateUserUseCase {
}`,
		"/project/src/main/java/com/example/app/application/gateway/UserGateway.java": `package com.example.app.application.gateway;

public interface UserGateway {
}`,
		"/project/src/main/java/com/example/app/infrastructure/web/UserController.java": `package com.example.app.infrastructure.web;

import org.springframework.web.bind.annotation.*;

@RestController
public class UserController {
}`,
		"/project/src/main/java/com/example/app/infrastructure/persistence/UserJpaGateway.java": `package com.example.app.infrastructure.persistence;

import org.springframework.stereotype.Repository;
import com.example.app.application.gateway.UserGateway;

@Repository
public class UserJpaGateway implements UserGateway {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		dir := path[:len(path)-len(path[len(path)-len(getFileName(path)):])-1]
		require.NoError(t, fs.MkdirAll(dir, 0755))
		require.NoError(t, afero.WriteFile(fs, path, []byte(content), 0644))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()
	require.NoError(t, err)

	t.Logf("=== Clean Architecture Detection ===")
	t.Logf("Architecture: %s (%.0f%% confidence)", profile.Architecture, profile.ArchConfidence*100)

	assert.Equal(t, ArchClean, profile.Architecture)
}

func getFileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
