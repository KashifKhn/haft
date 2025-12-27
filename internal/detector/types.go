package detector

type ArchitectureType string

const (
	ArchLayered   ArchitectureType = "layered"
	ArchFeature   ArchitectureType = "feature"
	ArchHexagonal ArchitectureType = "hexagonal"
	ArchClean     ArchitectureType = "clean"
	ArchModular   ArchitectureType = "modular"
	ArchFlat      ArchitectureType = "flat"
	ArchUnknown   ArchitectureType = "unknown"
)

type DTONamingStyle string

const (
	DTONamingRequestResponse DTONamingStyle = "request_response"
	DTONamingDTOUpper        DTONamingStyle = "dto_upper"
	DTONamingDTOLower        DTONamingStyle = "dto_lower"
	DTONamingUnknown         DTONamingStyle = "unknown"
)

type MapperType string

const (
	MapperMapStruct   MapperType = "mapstruct"
	MapperModelMapper MapperType = "modelmapper"
	MapperManual      MapperType = "manual"
	MapperNone        MapperType = "none"
)

type DatabaseType string

const (
	DatabaseJPA       DatabaseType = "jpa"
	DatabaseCassandra DatabaseType = "cassandra"
	DatabaseMongo     DatabaseType = "mongo"
	DatabaseR2DBC     DatabaseType = "r2dbc"
	DatabaseMulti     DatabaseType = "multi"
	DatabaseUnknown   DatabaseType = "unknown"
)

type SwaggerStyle string

const (
	SwaggerOpenAPI3 SwaggerStyle = "openapi3"
	SwaggerV2       SwaggerStyle = "swagger2"
	SwaggerNone     SwaggerStyle = "none"
)

type ValidationStyle string

const (
	ValidationJakarta ValidationStyle = "jakarta"
	ValidationJavax   ValidationStyle = "javax"
	ValidationNone    ValidationStyle = "none"
)

type JavaFileType string

const (
	FileTypeController JavaFileType = "controller"
	FileTypeService    JavaFileType = "service"
	FileTypeRepository JavaFileType = "repository"
	FileTypeEntity     JavaFileType = "entity"
	FileTypeDTO        JavaFileType = "dto"
	FileTypeMapper     JavaFileType = "mapper"
	FileTypeException  JavaFileType = "exception"
	FileTypeConfig     JavaFileType = "config"
	FileTypeTest       JavaFileType = "test"
	FileTypeUnknown    JavaFileType = "unknown"
)

type Field struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Annotations []string `yaml:"annotations,omitempty"`
}

type BaseClassInfo struct {
	Name       string  `yaml:"name"`
	Package    string  `yaml:"package"`
	FullPath   string  `yaml:"full_path"`
	Fields     []Field `yaml:"fields,omitempty"`
	IDType     string  `yaml:"id_type,omitempty"`
	IDStrategy string  `yaml:"id_strategy,omitempty"`
}

type WrapperInfo struct {
	Name           string   `yaml:"name"`
	Package        string   `yaml:"package"`
	FullPath       string   `yaml:"full_path"`
	IsGeneric      bool     `yaml:"is_generic"`
	Fields         []string `yaml:"fields,omitempty"`
	FactoryMethods []string `yaml:"factory_methods,omitempty"`
}

type ExceptionInfo struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
}

type ExceptionProfile struct {
	HasGlobalHandler bool            `yaml:"has_global_handler"`
	HandlerPackage   string          `yaml:"handler_package,omitempty"`
	CustomExceptions []ExceptionInfo `yaml:"custom_exceptions,omitempty"`
}

type LombokProfile struct {
	Detected        bool `yaml:"detected"`
	UseData         bool `yaml:"use_data"`
	UseBuilder      bool `yaml:"use_builder"`
	UseAccessors    bool `yaml:"use_accessors"`
	UseSlf4j        bool `yaml:"use_slf4j"`
	UseRequiredArgs bool `yaml:"use_required_args"`
	UseAllArgs      bool `yaml:"use_all_args"`
	UseNoArgs       bool `yaml:"use_no_args"`
}

type TestProfile struct {
	Framework         string `yaml:"framework"`
	HasMockito        bool   `yaml:"has_mockito"`
	HasTestcontainers bool   `yaml:"has_testcontainers"`
	HasRestAssured    bool   `yaml:"has_rest_assured"`
	StructureMirror   bool   `yaml:"structure_mirror"`
}

type JavaFile struct {
	Path                 string
	Package              string
	ClassName            string
	FileType             JavaFileType
	Annotations          []string
	ExtendsClass         string
	ImplementsInterfaces []string
	Imports              []string
	IsAbstract           bool
	IsInterface          bool
}

type DetectionResult struct {
	Value      interface{}
	Confidence float64
	Evidence   []string
}

func (dt ArchitectureType) String() string {
	return string(dt)
}

func (dt ArchitectureType) IsValid() bool {
	switch dt {
	case ArchLayered, ArchFeature, ArchHexagonal, ArchClean, ArchModular, ArchFlat:
		return true
	}
	return false
}

func ParseArchitectureType(s string) ArchitectureType {
	switch s {
	case "layered":
		return ArchLayered
	case "feature", "package-by-feature":
		return ArchFeature
	case "hexagonal", "ports-and-adapters":
		return ArchHexagonal
	case "clean":
		return ArchClean
	case "modular", "modular-monolith":
		return ArchModular
	case "flat":
		return ArchFlat
	default:
		return ArchUnknown
	}
}

func (dt DTONamingStyle) String() string {
	return string(dt)
}

func ParseDTONamingStyle(s string) DTONamingStyle {
	switch s {
	case "request_response":
		return DTONamingRequestResponse
	case "dto_upper", "DTO":
		return DTONamingDTOUpper
	case "dto_lower", "Dto":
		return DTONamingDTOLower
	default:
		return DTONamingUnknown
	}
}

func (mt MapperType) String() string {
	return string(mt)
}

func ParseMapperType(s string) MapperType {
	switch s {
	case "mapstruct":
		return MapperMapStruct
	case "modelmapper":
		return MapperModelMapper
	case "manual":
		return MapperManual
	default:
		return MapperNone
	}
}

func (dt DatabaseType) String() string {
	return string(dt)
}

func ParseDatabaseType(s string) DatabaseType {
	switch s {
	case "jpa":
		return DatabaseJPA
	case "cassandra":
		return DatabaseCassandra
	case "mongo", "mongodb":
		return DatabaseMongo
	case "r2dbc":
		return DatabaseR2DBC
	case "multi":
		return DatabaseMulti
	default:
		return DatabaseUnknown
	}
}
