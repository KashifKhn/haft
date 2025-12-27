package detector

import "time"

type ProjectProfile struct {
	DetectedAt  time.Time `yaml:"detected_at"`
	ProjectRoot string    `yaml:"project_root"`

	Architecture   ArchitectureType `yaml:"architecture"`
	ArchConfidence float64          `yaml:"arch_confidence"`
	ArchLocked     bool             `yaml:"arch_locked"`
	FeatureStyle   FeatureStyle     `yaml:"feature_style"`

	BasePackage string `yaml:"base_package"`
	SourceRoot  string `yaml:"source_root"`
	TestRoot    string `yaml:"test_root"`

	BaseEntity     *BaseClassInfo `yaml:"base_entity,omitempty"`
	BaseRepository *BaseClassInfo `yaml:"base_repository,omitempty"`
	BaseService    *BaseClassInfo `yaml:"base_service,omitempty"`

	ResponseWrapper *WrapperInfo `yaml:"response_wrapper,omitempty"`
	PageWrapper     *WrapperInfo `yaml:"page_wrapper,omitempty"`

	DTONaming        DTONamingStyle `yaml:"dto_naming"`
	DTONamingLocked  bool           `yaml:"dto_naming_locked"`
	ControllerSuffix string         `yaml:"controller_suffix"`
	ServiceSuffix    string         `yaml:"service_suffix"`

	IDType       string `yaml:"id_type"`
	IDAnnotation string `yaml:"id_annotation"`
	IDLocked     bool   `yaml:"id_locked"`

	Mapper       MapperType `yaml:"mapper"`
	MapperLocked bool       `yaml:"mapper_locked"`

	Exceptions ExceptionProfile `yaml:"exceptions"`

	Lombok LombokProfile `yaml:"lombok"`

	HasSwagger   bool         `yaml:"has_swagger"`
	SwaggerStyle SwaggerStyle `yaml:"swagger_style"`

	HasValidation   bool            `yaml:"has_validation"`
	ValidationStyle ValidationStyle `yaml:"validation_style"`

	Testing TestProfile `yaml:"testing"`

	Database       DatabaseType `yaml:"database"`
	DatabaseLocked bool         `yaml:"database_locked"`

	FeatureModules []string `yaml:"feature_modules,omitempty"`

	LockedFields []string `yaml:"locked_fields,omitempty"`
}

func NewEmptyProfile() *ProjectProfile {
	return &ProjectProfile{
		DetectedAt:       time.Now(),
		Architecture:     ArchUnknown,
		DTONaming:        DTONamingUnknown,
		Mapper:           MapperNone,
		SwaggerStyle:     SwaggerNone,
		ValidationStyle:  ValidationNone,
		Database:         DatabaseUnknown,
		ControllerSuffix: "Controller",
		ServiceSuffix:    "Service",
		IDType:           "Long",
		SourceRoot:       "src/main/java",
		TestRoot:         "src/test/java",
	}
}

func NewDefaultProfile() *ProjectProfile {
	return &ProjectProfile{
		DetectedAt:       time.Now(),
		Architecture:     ArchLayered,
		ArchConfidence:   1.0,
		DTONaming:        DTONamingRequestResponse,
		Mapper:           MapperManual,
		SwaggerStyle:     SwaggerNone,
		ValidationStyle:  ValidationJakarta,
		Database:         DatabaseJPA,
		ControllerSuffix: "Controller",
		ServiceSuffix:    "Service",
		IDType:           "Long",
		IDAnnotation:     "@GeneratedValue(strategy = GenerationType.IDENTITY)",
		SourceRoot:       "src/main/java",
		TestRoot:         "src/test/java",
		Lombok: LombokProfile{
			Detected:   true,
			UseData:    true,
			UseBuilder: false,
			UseNoArgs:  true,
			UseAllArgs: true,
		},
		Testing: TestProfile{
			Framework:       "junit5",
			HasMockito:      true,
			StructureMirror: true,
		},
	}
}

func (p *ProjectProfile) IsFieldLocked(fieldName string) bool {
	for _, locked := range p.LockedFields {
		if locked == fieldName {
			return true
		}
	}
	return false
}

func (p *ProjectProfile) LockField(fieldName string) {
	if !p.IsFieldLocked(fieldName) {
		p.LockedFields = append(p.LockedFields, fieldName)
	}
}

func (p *ProjectProfile) UnlockField(fieldName string) {
	for i, locked := range p.LockedFields {
		if locked == fieldName {
			p.LockedFields = append(p.LockedFields[:i], p.LockedFields[i+1:]...)
			return
		}
	}
}

func (p *ProjectProfile) IsValid() bool {
	return p.Architecture != ArchUnknown && p.BasePackage != ""
}

func (p *ProjectProfile) IsStale(maxAge time.Duration) bool {
	return time.Since(p.DetectedAt) > maxAge
}

func (p *ProjectProfile) GetEntityPackage(resourceName string) string {
	return p.computePackagePath(resourceName, "entity")
}

func (p *ProjectProfile) GetControllerPackage(resourceName string) string {
	return p.computePackagePath(resourceName, "controller")
}

func (p *ProjectProfile) GetServicePackage(resourceName string) string {
	return p.computePackagePath(resourceName, "service")
}

func (p *ProjectProfile) GetRepositoryPackage(resourceName string) string {
	return p.computePackagePath(resourceName, "repository")
}

func (p *ProjectProfile) GetDTOPackage(resourceName string) string {
	return p.computePackagePath(resourceName, "dto")
}

func (p *ProjectProfile) GetMapperPackage(resourceName string) string {
	return p.computePackagePath(resourceName, "mapper")
}

func (p *ProjectProfile) computePackagePath(resourceName, layerName string) string {
	resourceLower := toLowerFirst(resourceName)

	switch p.Architecture {
	case ArchFeature:
		if p.FeatureStyle == FeatureStyleFlat {
			if layerName == "dto" {
				return p.BasePackage + "." + resourceLower + ".dto"
			}
			return p.BasePackage + "." + resourceLower
		}
		return p.BasePackage + "." + resourceLower + "." + layerName
	case ArchHexagonal:
		return p.getHexagonalPackage(resourceLower, layerName)
	case ArchClean:
		return p.getCleanPackage(resourceLower, layerName)
	case ArchModular:
		return p.BasePackage + "." + resourceLower + ".internal." + layerName
	default:
		return p.BasePackage + "." + layerName
	}
}

func (p *ProjectProfile) getHexagonalPackage(resource, layerName string) string {
	switch layerName {
	case "controller":
		return p.BasePackage + ".adapter.in.web"
	case "service":
		return p.BasePackage + ".application.service"
	case "entity":
		return p.BasePackage + ".domain.model"
	case "repository":
		return p.BasePackage + ".application.port.out"
	case "dto":
		return p.BasePackage + ".adapter.in.web.dto"
	case "mapper":
		return p.BasePackage + ".adapter.in.web.mapper"
	default:
		return p.BasePackage + "." + layerName
	}
}

func (p *ProjectProfile) getCleanPackage(resource, layerName string) string {
	switch layerName {
	case "controller":
		return p.BasePackage + ".infrastructure.web"
	case "service":
		return p.BasePackage + ".application.usecase"
	case "entity":
		return p.BasePackage + ".domain.entity"
	case "repository":
		return p.BasePackage + ".application.gateway"
	case "dto":
		return p.BasePackage + ".infrastructure.web.dto"
	case "mapper":
		return p.BasePackage + ".infrastructure.web.mapper"
	default:
		return p.BasePackage + "." + layerName
	}
}

func (p *ProjectProfile) GetDTORequestSuffix() string {
	switch p.DTONaming {
	case DTONamingRequestResponse:
		return "Request"
	case DTONamingDTOUpper:
		return "DTO"
	case DTONamingDTOLower:
		return "Dto"
	default:
		return "Request"
	}
}

func (p *ProjectProfile) GetDTOResponseSuffix() string {
	switch p.DTONaming {
	case DTONamingRequestResponse:
		return "Response"
	case DTONamingDTOUpper:
		return "DTO"
	case DTONamingDTOLower:
		return "Dto"
	default:
		return "Response"
	}
}

func (p *ProjectProfile) GetIDImport() string {
	switch p.IDType {
	case "UUID":
		return "java.util.UUID"
	default:
		return ""
	}
}

func (p *ProjectProfile) NeedsBaseEntityImport() bool {
	return p.BaseEntity != nil && p.BaseEntity.Package != ""
}

func (p *ProjectProfile) GetBaseEntityImport() string {
	if p.BaseEntity == nil {
		return ""
	}
	return p.BaseEntity.Package + "." + p.BaseEntity.Name
}

func (p *ProjectProfile) NeedsResponseWrapperImport() bool {
	return p.ResponseWrapper != nil && p.ResponseWrapper.Package != ""
}

func (p *ProjectProfile) GetResponseWrapperImport() string {
	if p.ResponseWrapper == nil {
		return ""
	}
	return p.ResponseWrapper.Package + "." + p.ResponseWrapper.Name
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]|32) + s[1:]
}
