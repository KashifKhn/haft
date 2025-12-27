package detector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyProfile(t *testing.T) {
	profile := NewEmptyProfile()

	assert.NotNil(t, profile)
	assert.Equal(t, ArchUnknown, profile.Architecture)
	assert.Equal(t, DTONamingUnknown, profile.DTONaming)
	assert.Equal(t, MapperNone, profile.Mapper)
	assert.Equal(t, "Controller", profile.ControllerSuffix)
	assert.Equal(t, "Service", profile.ServiceSuffix)
	assert.Equal(t, "Long", profile.IDType)
	assert.Equal(t, "src/main/java", profile.SourceRoot)
	assert.Equal(t, "src/test/java", profile.TestRoot)
}

func TestNewDefaultProfile(t *testing.T) {
	profile := NewDefaultProfile()

	assert.NotNil(t, profile)
	assert.Equal(t, ArchLayered, profile.Architecture)
	assert.Equal(t, 1.0, profile.ArchConfidence)
	assert.Equal(t, DTONamingRequestResponse, profile.DTONaming)
	assert.Equal(t, MapperManual, profile.Mapper)
	assert.True(t, profile.Lombok.Detected)
	assert.True(t, profile.Lombok.UseData)
}

func TestProfileFieldLocking(t *testing.T) {
	profile := NewEmptyProfile()

	assert.False(t, profile.IsFieldLocked("Architecture"))

	profile.LockField("Architecture")
	assert.True(t, profile.IsFieldLocked("Architecture"))

	profile.LockField("Architecture")
	assert.Equal(t, 1, len(profile.LockedFields))

	profile.UnlockField("Architecture")
	assert.False(t, profile.IsFieldLocked("Architecture"))

	profile.UnlockField("NonExistent")
}

func TestProfileIsValid(t *testing.T) {
	tests := []struct {
		name     string
		profile  *ProjectProfile
		expected bool
	}{
		{
			name:     "empty profile is invalid",
			profile:  NewEmptyProfile(),
			expected: false,
		},
		{
			name: "profile with architecture but no package is invalid",
			profile: &ProjectProfile{
				Architecture: ArchLayered,
			},
			expected: false,
		},
		{
			name: "profile with package but unknown architecture is invalid",
			profile: &ProjectProfile{
				Architecture: ArchUnknown,
				BasePackage:  "com.example",
			},
			expected: false,
		},
		{
			name: "profile with architecture and package is valid",
			profile: &ProjectProfile{
				Architecture: ArchLayered,
				BasePackage:  "com.example",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.profile.IsValid())
		})
	}
}

func TestProfileIsStale(t *testing.T) {
	profile := NewEmptyProfile()
	profile.DetectedAt = time.Now()

	assert.False(t, profile.IsStale(time.Hour))

	profile.DetectedAt = time.Now().Add(-2 * time.Hour)
	assert.True(t, profile.IsStale(time.Hour))
}

func TestProfilePackagePathsLayered(t *testing.T) {
	profile := &ProjectProfile{
		Architecture: ArchLayered,
		BasePackage:  "com.example.app",
	}

	assert.Equal(t, "com.example.app.entity", profile.GetEntityPackage("User"))
	assert.Equal(t, "com.example.app.controller", profile.GetControllerPackage("User"))
	assert.Equal(t, "com.example.app.service", profile.GetServicePackage("User"))
	assert.Equal(t, "com.example.app.repository", profile.GetRepositoryPackage("User"))
	assert.Equal(t, "com.example.app.dto", profile.GetDTOPackage("User"))
	assert.Equal(t, "com.example.app.mapper", profile.GetMapperPackage("User"))
}

func TestProfilePackagePathsFeature(t *testing.T) {
	profile := &ProjectProfile{
		Architecture: ArchFeature,
		BasePackage:  "com.example.app",
	}

	assert.Equal(t, "com.example.app.user.entity", profile.GetEntityPackage("User"))
	assert.Equal(t, "com.example.app.user.controller", profile.GetControllerPackage("User"))
	assert.Equal(t, "com.example.app.user.service", profile.GetServicePackage("User"))
	assert.Equal(t, "com.example.app.user.repository", profile.GetRepositoryPackage("User"))
	assert.Equal(t, "com.example.app.user.dto", profile.GetDTOPackage("User"))
	assert.Equal(t, "com.example.app.user.mapper", profile.GetMapperPackage("User"))
}

func TestProfilePackagePathsHexagonal(t *testing.T) {
	profile := &ProjectProfile{
		Architecture: ArchHexagonal,
		BasePackage:  "com.example.app",
	}

	assert.Equal(t, "com.example.app.domain.model", profile.GetEntityPackage("User"))
	assert.Equal(t, "com.example.app.adapter.in.web", profile.GetControllerPackage("User"))
	assert.Equal(t, "com.example.app.application.service", profile.GetServicePackage("User"))
	assert.Equal(t, "com.example.app.application.port.out", profile.GetRepositoryPackage("User"))
	assert.Equal(t, "com.example.app.adapter.in.web.dto", profile.GetDTOPackage("User"))
	assert.Equal(t, "com.example.app.adapter.in.web.mapper", profile.GetMapperPackage("User"))
}

func TestProfilePackagePathsClean(t *testing.T) {
	profile := &ProjectProfile{
		Architecture: ArchClean,
		BasePackage:  "com.example.app",
	}

	assert.Equal(t, "com.example.app.domain.entity", profile.GetEntityPackage("User"))
	assert.Equal(t, "com.example.app.infrastructure.web", profile.GetControllerPackage("User"))
	assert.Equal(t, "com.example.app.application.usecase", profile.GetServicePackage("User"))
	assert.Equal(t, "com.example.app.application.gateway", profile.GetRepositoryPackage("User"))
}

func TestProfilePackagePathsModular(t *testing.T) {
	profile := &ProjectProfile{
		Architecture: ArchModular,
		BasePackage:  "com.example.app",
	}

	assert.Equal(t, "com.example.app.user.internal.entity", profile.GetEntityPackage("User"))
	assert.Equal(t, "com.example.app.user.internal.controller", profile.GetControllerPackage("User"))
}

func TestProfileDTOSuffixes(t *testing.T) {
	tests := []struct {
		naming           DTONamingStyle
		expectedRequest  string
		expectedResponse string
	}{
		{DTONamingRequestResponse, "Request", "Response"},
		{DTONamingDTOUpper, "DTO", "DTO"},
		{DTONamingDTOLower, "Dto", "Dto"},
		{DTONamingUnknown, "Request", "Response"},
	}

	for _, tt := range tests {
		t.Run(string(tt.naming), func(t *testing.T) {
			profile := &ProjectProfile{DTONaming: tt.naming}
			assert.Equal(t, tt.expectedRequest, profile.GetDTORequestSuffix())
			assert.Equal(t, tt.expectedResponse, profile.GetDTOResponseSuffix())
		})
	}
}

func TestProfileIDImport(t *testing.T) {
	tests := []struct {
		idType   string
		expected string
	}{
		{"UUID", "java.util.UUID"},
		{"Long", ""},
		{"Integer", ""},
		{"String", ""},
	}

	for _, tt := range tests {
		t.Run(tt.idType, func(t *testing.T) {
			profile := &ProjectProfile{IDType: tt.idType}
			assert.Equal(t, tt.expected, profile.GetIDImport())
		})
	}
}

func TestProfileBaseEntityImport(t *testing.T) {
	t.Run("no base entity", func(t *testing.T) {
		profile := &ProjectProfile{}
		assert.False(t, profile.NeedsBaseEntityImport())
		assert.Equal(t, "", profile.GetBaseEntityImport())
	})

	t.Run("with base entity", func(t *testing.T) {
		profile := &ProjectProfile{
			BaseEntity: &BaseClassInfo{
				Name:    "BaseEntity",
				Package: "com.example.common.entity",
			},
		}
		assert.True(t, profile.NeedsBaseEntityImport())
		assert.Equal(t, "com.example.common.entity.BaseEntity", profile.GetBaseEntityImport())
	})
}

func TestProfileResponseWrapperImport(t *testing.T) {
	t.Run("no wrapper", func(t *testing.T) {
		profile := &ProjectProfile{}
		assert.False(t, profile.NeedsResponseWrapperImport())
		assert.Equal(t, "", profile.GetResponseWrapperImport())
	})

	t.Run("with wrapper", func(t *testing.T) {
		profile := &ProjectProfile{
			ResponseWrapper: &WrapperInfo{
				Name:    "ApiResponse",
				Package: "com.example.common.dto",
			},
		}
		assert.True(t, profile.NeedsResponseWrapperImport())
		assert.Equal(t, "com.example.common.dto.ApiResponse", profile.GetResponseWrapperImport())
	})
}

func TestToLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"UserProfile", "userProfile"},
		{"A", "a"},
		{"", ""},
		{"user", "user"},
		{"ABC", "aBC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, toLowerFirst(tt.input))
		})
	}
}
