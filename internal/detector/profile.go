package detector

import "time"

type ProjectProfile struct {
	DetectedAt  time.Time `json:"detected_at"`
	ProjectRoot string    `json:"project_root"`

	Architecture   ArchitectureType `json:"architecture"`
	ArchConfidence float64          `json:"arch_confidence"`
	ArchLocked     bool             `json:"arch_locked"`
	FeatureStyle   FeatureStyle     `json:"feature_style"`

	BasePackage string `json:"base_package"`
	SourceRoot  string `json:"source_root"`
	TestRoot    string `json:"test_root"`

	BaseEntity     *BaseClassInfo `json:"base_entity,omitempty"`
	BaseRepository *BaseClassInfo `json:"base_repository,omitempty"`
	BaseService    *BaseClassInfo `json:"base_service,omitempty"`

	ResponseWrapper *WrapperInfo `json:"response_wrapper,omitempty"`
	PageWrapper     *WrapperInfo `json:"page_wrapper,omitempty"`

	DTONaming        DTONamingStyle `json:"dto_naming"`
	DTONamingLocked  bool           `json:"dto_naming_locked"`
	ControllerSuffix string         `json:"controller_suffix"`
	ServiceSuffix    string         `json:"service_suffix"`

	IDType       string `json:"id_type"`
	IDAnnotation string `json:"id_annotation"`
	IDLocked     bool   `json:"id_locked"`

	Mapper       MapperType `json:"mapper"`
	MapperLocked bool       `json:"mapper_locked"`

	Exceptions ExceptionProfile `json:"exceptions"`

	Lombok LombokProfile `json:"lombok"`

	HasSwagger   bool         `json:"has_swagger"`
	SwaggerStyle SwaggerStyle `json:"swagger_style"`

	HasValidation   bool            `json:"has_validation"`
	ValidationStyle ValidationStyle `json:"validation_style"`

	Testing TestProfile `json:"testing"`

	Database       DatabaseType `json:"database"`
	DatabaseLocked bool         `json:"database_locked"`

	FeatureModules []string `json:"feature_modules,omitempty"`

	LockedFields []string `json:"locked_fields,omitempty"`
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
