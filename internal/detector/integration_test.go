package detector

import (
	"os"
	"testing"

	"github.com/spf13/afero"
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

	if profile.Architecture != ArchFeature {
		t.Errorf("Expected feature architecture, got %s", profile.Architecture)
	}

	if profile.DTONaming != DTONamingRequestResponse {
		t.Errorf("Expected request_response DTO naming, got %s", profile.DTONaming)
	}

	if profile.IDType != "UUID" {
		t.Errorf("Expected UUID ID type, got %s", profile.IDType)
	}

	if profile.Mapper != MapperMapStruct {
		t.Errorf("Expected MapStruct mapper, got %s", profile.Mapper)
	}
}
