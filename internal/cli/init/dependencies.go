package init

type Dependency struct {
	GroupId    string
	ArtifactId string
	Version    string
	Scope      string
}

func (d Dependency) GradleScope() string {
	switch d.Scope {
	case "provided":
		return "compileOnly"
	case "runtime":
		return "runtimeOnly"
	case "test":
		return "testImplementation"
	case "annotationProcessor":
		return "annotationProcessor"
	default:
		return "implementation"
	}
}

var dependencyMap = map[string]Dependency{
	"web": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-web",
	},
	"webflux": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-webflux",
	},
	"data-jpa": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-data-jpa",
	},
	"jpa": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-data-jpa",
	},
	"security": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-security",
	},
	"validation": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-validation",
	},
	"actuator": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-actuator",
	},
	"devtools": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-devtools",
		Scope:      "runtime",
	},
	"lombok": {
		GroupId:    "org.projectlombok",
		ArtifactId: "lombok",
		Scope:      "provided",
	},
	"h2": {
		GroupId:    "com.h2database",
		ArtifactId: "h2",
		Scope:      "runtime",
	},
	"postgresql": {
		GroupId:    "org.postgresql",
		ArtifactId: "postgresql",
		Scope:      "runtime",
	},
	"mysql": {
		GroupId:    "com.mysql",
		ArtifactId: "mysql-connector-j",
		Scope:      "runtime",
	},
	"data-mongodb": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-data-mongodb",
	},
	"data-redis": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-data-redis",
	},
	"thymeleaf": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-thymeleaf",
	},
	"freemarker": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-freemarker",
	},
	"mail": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-mail",
	},
	"amqp": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-amqp",
	},
	"kafka": {
		GroupId:    "org.springframework.kafka",
		ArtifactId: "spring-kafka",
	},
	"websocket": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-websocket",
	},
	"graphql": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-graphql",
	},
	"data-rest": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-data-rest",
	},
	"hateoas": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-hateoas",
	},
	"oauth2-client": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-oauth2-client",
	},
	"oauth2-resource-server": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-oauth2-resource-server",
	},
	"cache": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-cache",
	},
	"batch": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-batch",
	},
	"quartz": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-starter-quartz",
	},
	"flyway": {
		GroupId:    "org.flywaydb",
		ArtifactId: "flyway-core",
	},
	"liquibase": {
		GroupId:    "org.liquibase",
		ArtifactId: "liquibase-core",
	},
	"testcontainers": {
		GroupId:    "org.testcontainers",
		ArtifactId: "testcontainers",
		Scope:      "test",
	},
	"docker-compose": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-docker-compose",
		Scope:      "runtime",
	},
	"configuration-processor": {
		GroupId:    "org.springframework.boot",
		ArtifactId: "spring-boot-configuration-processor",
		Scope:      "annotationProcessor",
	},
}

func buildDependencies(deps []string) []Dependency {
	normalizedDeps := normalizeDependencies(deps)

	var result []Dependency
	seen := make(map[string]bool)

	for _, dep := range normalizedDeps {
		if d, ok := dependencyMap[dep]; ok {
			key := d.GroupId + ":" + d.ArtifactId
			if !seen[key] {
				result = append(result, d)
				seen[key] = true
			}
		}
	}
	return result
}

func normalizeDependencies(deps []string) []string {
	hasJpa := contains(deps, "data-jpa") || contains(deps, "jpa")
	hasH2 := contains(deps, "h2")
	hasPostgres := contains(deps, "postgresql")
	hasMysql := contains(deps, "mysql")
	hasMongo := contains(deps, "data-mongodb")

	if hasJpa && !hasH2 && !hasPostgres && !hasMysql && !hasMongo {
		deps = append(deps, "h2")
	}

	return deps
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
