package detector

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDetector(t *testing.T) {
	d := NewDetector("/project")

	assert.NotNil(t, d)
	assert.Equal(t, "/project", d.projectDir)
	assert.NotNil(t, d.scanner)
	assert.NotNil(t, d.confidenceCalculator)
	assert.Equal(t, DefaultCacheMaxAge, d.cacheMaxAge)
}

func TestNewDetectorWithOptions(t *testing.T) {
	fs := afero.NewMemMapFs()
	d := NewDetector("/project",
		WithFileSystem(fs),
		WithCacheMaxAge(DefaultCacheMaxAge*2),
	)

	assert.NotNil(t, d)
	assert.Equal(t, DefaultCacheMaxAge*2, d.cacheMaxAge)
}

func TestDetectorDetectLayeredArchitecture(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/app/controller/UserController.java": `package com.example.app.controller;

@RestController
public class UserController {
}`,
		"/project/src/main/java/com/example/app/service/UserService.java": `package com.example.app.service;

@Service
public class UserService {
}`,
		"/project/src/main/java/com/example/app/entity/User.java": `package com.example.app.entity;

@Entity
public class User {
}`,
		"/project/src/main/java/com/example/app/repository/UserRepository.java": `package com.example.app.repository;

@Repository
public interface UserRepository {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.Equal(t, ArchLayered, profile.Architecture)
	assert.GreaterOrEqual(t, profile.ArchConfidence, 0.5)
	assert.Equal(t, "com.example.app", profile.BasePackage)
}

func TestDetectorDetectFeatureArchitecture(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/app/user/controller/UserController.java": `package com.example.app.user.controller;

@RestController
public class UserController {
}`,
		"/project/src/main/java/com/example/app/user/service/UserService.java": `package com.example.app.user.service;

@Service
public class UserService {
}`,
		"/project/src/main/java/com/example/app/user/entity/User.java": `package com.example.app.user.entity;

@Entity
public class User {
}`,
		"/project/src/main/java/com/example/app/user/repository/UserRepository.java": `package com.example.app.user.repository;

@Repository
public interface UserRepository {
}`,
		"/project/src/main/java/com/example/app/auth/controller/AuthController.java": `package com.example.app.auth.controller;

@RestController
public class AuthController {
}`,
		"/project/src/main/java/com/example/app/auth/service/AuthService.java": `package com.example.app.auth.service;

@Service
public class AuthService {
}`,
		"/project/src/main/java/com/example/app/common/entity/BaseEntity.java": `package com.example.app.common.entity;

@MappedSuperclass
public abstract class BaseEntity {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.Equal(t, ArchFeature, profile.Architecture)
	assert.Contains(t, profile.FeatureModules, "user")
	assert.Contains(t, profile.FeatureModules, "auth")
}

func TestDetectorDetectBaseEntity(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/common/entity/BaseEntity.java": `package com.example.common.entity;

import java.util.UUID;

@MappedSuperclass
public abstract class BaseEntity {
    @Id
    private UUID id;
}`,
		"/project/src/main/java/com/example/user/entity/User.java": `package com.example.user.entity;

@Entity
public class User extends BaseEntity {
}`,
		"/project/src/main/java/com/example/product/entity/Product.java": `package com.example.product.entity;

@Entity
public class Product extends BaseEntity {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	require.NotNil(t, profile.BaseEntity)
	assert.Equal(t, "BaseEntity", profile.BaseEntity.Name)
	assert.Equal(t, "com.example.common.entity", profile.BaseEntity.Package)
}

func TestDetectorDetectDTONaming(t *testing.T) {
	tests := []struct {
		name     string
		dtoFiles map[string]string
		expected DTONamingStyle
	}{
		{
			name: "request response style",
			dtoFiles: map[string]string{
				"/project/src/main/java/com/example/dto/UserRequest.java":       `package com.example.dto; public class UserRequest {}`,
				"/project/src/main/java/com/example/dto/UserResponse.java":      `package com.example.dto; public class UserResponse {}`,
				"/project/src/main/java/com/example/dto/CreateUserRequest.java": `package com.example.dto; public class CreateUserRequest {}`,
			},
			expected: DTONamingRequestResponse,
		},
		{
			name: "DTO uppercase style",
			dtoFiles: map[string]string{
				"/project/src/main/java/com/example/dto/UserDTO.java":    `package com.example.dto; public class UserDTO {}`,
				"/project/src/main/java/com/example/dto/ProductDTO.java": `package com.example.dto; public class ProductDTO {}`,
			},
			expected: DTONamingDTOUpper,
		},
		{
			name: "Dto lowercase style",
			dtoFiles: map[string]string{
				"/project/src/main/java/com/example/dto/UserDto.java":    `package com.example.dto; public class UserDto {}`,
				"/project/src/main/java/com/example/dto/ProductDto.java": `package com.example.dto; public class ProductDto {}`,
			},
			expected: DTONamingDTOLower,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tt.dtoFiles {
				require.NoError(t, createJavaFile(fs, path, content))
			}
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.expected, profile.DTONaming)
		})
	}
}

func TestDetectorDetectIDType(t *testing.T) {
	tests := []struct {
		name        string
		entityFiles map[string]string
		expectedID  string
	}{
		{
			name: "UUID id type",
			entityFiles: map[string]string{
				"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

import java.util.UUID;

@Entity
public class User {
    @Id
    private UUID id;
}`,
			},
			expectedID: "UUID",
		},
		{
			name: "Long id type (default)",
			entityFiles: map[string]string{
				"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

@Entity
public class User {
    @Id
    private Long id;
}`,
			},
			expectedID: "Long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tt.entityFiles {
				require.NoError(t, createJavaFile(fs, path, content))
			}
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedID, profile.IDType)
		})
	}
}

func TestDetectorDetectMapStruct(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/mapper/UserMapper.java": `package com.example.mapper;

import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface UserMapper {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.Equal(t, MapperMapStruct, profile.Mapper)
}

func TestDetectorDetectLombok(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
public class User {
}`,
		"/project/src/main/java/com/example/service/UserService.java": `package com.example.service;

@Slf4j
@Service
@RequiredArgsConstructor
public class UserService {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.True(t, profile.Lombok.Detected)
	assert.True(t, profile.Lombok.UseData)
	assert.True(t, profile.Lombok.UseBuilder)
	assert.True(t, profile.Lombok.UseNoArgs)
	assert.True(t, profile.Lombok.UseAllArgs)
	assert.True(t, profile.Lombok.UseSlf4j)
	assert.True(t, profile.Lombok.UseRequiredArgs)
}

func TestDetectorDetectExceptions(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/exception/GlobalExceptionHandler.java": `package com.example.exception;

@RestControllerAdvice
public class GlobalExceptionHandler {
}`,
		"/project/src/main/java/com/example/exception/ResourceNotFoundException.java": `package com.example.exception;

public class ResourceNotFoundException extends RuntimeException {
}`,
		"/project/src/main/java/com/example/exception/BadRequestException.java": `package com.example.exception;

public class BadRequestException extends RuntimeException {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.True(t, profile.Exceptions.HasGlobalHandler)
	assert.Equal(t, "com.example.exception", profile.Exceptions.HandlerPackage)
	assert.GreaterOrEqual(t, len(profile.Exceptions.CustomExceptions), 2)
}

func TestDetectorDetectSwagger(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		hasSwagger bool
		style      SwaggerStyle
	}{
		{
			name: "OpenAPI 3",
			content: `package com.example;

import io.swagger.v3.oas.annotations.Operation;

@RestController
public class UserController {
}`,
			hasSwagger: true,
			style:      SwaggerOpenAPI3,
		},
		{
			name: "Swagger 2",
			content: `package com.example;

import io.swagger.annotations.Api;

@Api
@RestController
public class UserController {
}`,
			hasSwagger: true,
			style:      SwaggerV2,
		},
		{
			name: "No Swagger",
			content: `package com.example;

@RestController
public class UserController {
}`,
			hasSwagger: false,
			style:      SwaggerNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			require.NoError(t, createJavaFile(fs, "/project/src/main/java/UserController.java", tt.content))
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.hasSwagger, profile.HasSwagger)
			assert.Equal(t, tt.style, profile.SwaggerStyle)
		})
	}
}

func TestDetectorDetectValidation(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		hasValidation   bool
		validationStyle ValidationStyle
	}{
		{
			name: "Jakarta validation",
			content: `package com.example;

import jakarta.validation.constraints.NotNull;

public class UserRequest {
    @NotNull
    private String name;
}`,
			hasValidation:   true,
			validationStyle: ValidationJakarta,
		},
		{
			name: "Javax validation",
			content: `package com.example;

import javax.validation.constraints.NotNull;

public class UserRequest {
    @NotNull
    private String name;
}`,
			hasValidation:   true,
			validationStyle: ValidationJavax,
		},
		{
			name: "No validation",
			content: `package com.example;

public class UserRequest {
    private String name;
}`,
			hasValidation:   false,
			validationStyle: ValidationNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			require.NoError(t, createJavaFile(fs, "/project/src/main/java/UserRequest.java", tt.content))
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.hasValidation, profile.HasValidation)
			assert.Equal(t, tt.validationStyle, profile.ValidationStyle)
		})
	}
}

func TestDetectorDetectDatabase(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected DatabaseType
	}{
		{
			name: "JPA only",
			files: map[string]string{
				"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

@Entity
public class User {
}`,
			},
			expected: DatabaseJPA,
		},
		{
			name: "MongoDB",
			files: map[string]string{
				"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

@Document
public class User {
}`,
			},
			expected: DatabaseMongo,
		},
		{
			name: "Multi database",
			files: map[string]string{
				"/project/src/main/java/com/example/entity/User.java": `package com.example.entity;

@Entity
public class User {
}`,
				"/project/src/main/java/com/example/entity/Log.java": `package com.example.entity;

import org.springframework.data.cassandra.core.mapping.Table;

@Table
public class Log {
}`,
			},
			expected: DatabaseMulti,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tt.files {
				require.NoError(t, createJavaFile(fs, path, content))
			}
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.expected, profile.Database)
		})
	}
}

func TestDetectorDetectTestingProfile(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"/project/src/main/java/com/example/service/UserService.java": `package com.example.service;

@Service
public class UserService {
}`,
		"/project/src/test/java/com/example/service/UserServiceTest.java": `package com.example.service;

import org.mockito.Mock;
import org.testcontainers.containers.PostgreSQLContainer;

@SpringBootTest
@Testcontainers
public class UserServiceTest {
}`,
		"/project/pom.xml": `<project></project>`,
	}

	for path, content := range files {
		require.NoError(t, createJavaFile(fs, path, content))
	}

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.Equal(t, "junit5", profile.Testing.Framework)
	assert.True(t, profile.Testing.HasMockito)
	assert.True(t, profile.Testing.HasTestcontainers)
}

func TestDetectorEmptyProject(t *testing.T) {
	fs := afero.NewMemMapFs()
	require.NoError(t, fs.MkdirAll("/project/src/main/java", 0755))
	require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

	d := NewDetector("/project", WithFileSystem(fs))
	profile, err := d.Detect()

	require.NoError(t, err)
	assert.Equal(t, ArchLayered, profile.Architecture)
	assert.Equal(t, 1.0, profile.ArchConfidence)
}

func TestHelperFunctions(t *testing.T) {
	t.Run("containsPackagePart", func(t *testing.T) {
		assert.True(t, containsPackagePart("com.example.controller", "controller"))
		assert.True(t, containsPackagePart("com.example.controller", "example"))
		assert.False(t, containsPackagePart("com.example.controller", "service"))
	})

	t.Run("isDirectLayerPackage", func(t *testing.T) {
		assert.True(t, isDirectLayerPackage("com.example.controller", "com.example", "controller"))
		assert.False(t, isDirectLayerPackage("com.example.controller.sub", "com.example", "controller"))
		assert.False(t, isDirectLayerPackage("com.example.user.controller", "com.example", "controller"))
	})

	t.Run("hasPackagePrefix", func(t *testing.T) {
		assert.True(t, hasPackagePrefix("com.example.app", "com.example"))
		assert.True(t, hasPackagePrefix("com.example", "com.example"))
		assert.False(t, hasPackagePrefix("com.example", "com.example.app"))
		assert.False(t, hasPackagePrefix("com.exampleapp", "com.example"))
	})

	t.Run("splitPackage", func(t *testing.T) {
		assert.Equal(t, []string{"com", "example", "app"}, splitPackage("com.example.app"))
		assert.Equal(t, []string{"app"}, splitPackage("app"))
		assert.Nil(t, splitPackage(""))
	})

	t.Run("isLayerName", func(t *testing.T) {
		assert.True(t, isLayerName("controller"))
		assert.True(t, isLayerName("service"))
		assert.True(t, isLayerName("repository"))
		assert.False(t, isLayerName("user"))
		assert.False(t, isLayerName("auth"))
	})

	t.Run("endsWith", func(t *testing.T) {
		assert.True(t, endsWith("UserController", "Controller"))
		assert.False(t, endsWith("UserService", "Controller"))
		assert.True(t, endsWith("Controller", "Controller"))
	})

	t.Run("endsWithAny", func(t *testing.T) {
		assert.True(t, endsWithAny("UserRequest", "Request", "Response"))
		assert.True(t, endsWithAny("UserResponse", "Request", "Response"))
		assert.False(t, endsWithAny("UserDTO", "Request", "Response"))
	})

	t.Run("containsString", func(t *testing.T) {
		assert.True(t, containsString("org.mapstruct.Mapper", "mapstruct"))
		assert.False(t, containsString("org.example.Mapper", "mapstruct"))
	})
}

func TestDetectorControllerSuffixDetection(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		expectedSuffix string
	}{
		{
			name: "Controller suffix",
			files: map[string]string{
				"/project/src/main/java/com/example/UserController.java": `package com.example;

@RestController
public class UserController {}`,
				"/project/src/main/java/com/example/AuthController.java": `package com.example;

@RestController
public class AuthController {}`,
			},
			expectedSuffix: "Controller",
		},
		{
			name: "Resource suffix",
			files: map[string]string{
				"/project/src/main/java/com/example/UserResource.java": `package com.example;

@RestController
public class UserResource {}`,
				"/project/src/main/java/com/example/AuthResource.java": `package com.example;

@RestController
public class AuthResource {}`,
			},
			expectedSuffix: "Resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tt.files {
				require.NoError(t, createJavaFile(fs, path, content))
			}
			require.NoError(t, afero.WriteFile(fs, "/project/pom.xml", []byte("<project></project>"), 0644))

			d := NewDetector("/project", WithFileSystem(fs))
			profile, err := d.Detect()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSuffix, profile.ControllerSuffix)
		})
	}
}
