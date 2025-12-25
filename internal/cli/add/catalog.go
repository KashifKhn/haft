package add

import "github.com/KashifKhn/haft/internal/maven"

type CatalogEntry struct {
	Name         string
	Description  string
	Category     string
	Dependencies []maven.Dependency
}

var dependencyCatalog = map[string]CatalogEntry{
	"web": {
		Name:        "Spring Web",
		Description: "Build web applications with Spring MVC",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		},
	},
	"webflux": {
		Name:        "Spring WebFlux",
		Description: "Build reactive web applications",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-webflux"},
		},
	},
	"jpa": {
		Name:        "Spring Data JPA",
		Description: "Persist data with JPA and Hibernate",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-jpa"},
		},
	},
	"jdbc": {
		Name:        "Spring JDBC",
		Description: "Database access with JDBC",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-jdbc"},
		},
	},
	"security": {
		Name:        "Spring Security",
		Description: "Authentication and authorization",
		Category:    "Security",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-security"},
		},
	},
	"oauth2-client": {
		Name:        "OAuth2 Client",
		Description: "OAuth2/OpenID Connect client",
		Category:    "Security",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-oauth2-client"},
		},
	},
	"oauth2-resource-server": {
		Name:        "OAuth2 Resource Server",
		Description: "OAuth2 resource server support",
		Category:    "Security",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-oauth2-resource-server"},
		},
	},
	"validation": {
		Name:        "Validation",
		Description: "Bean validation with Hibernate Validator",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-validation"},
		},
	},
	"actuator": {
		Name:        "Spring Boot Actuator",
		Description: "Production-ready monitoring and management",
		Category:    "Ops",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-actuator"},
		},
	},
	"devtools": {
		Name:        "Spring Boot DevTools",
		Description: "Fast application restarts and LiveReload",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-devtools", Scope: "runtime", Optional: "true"},
		},
	},
	"lombok": {
		Name:        "Lombok",
		Description: "Reduce boilerplate with annotations",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.projectlombok", ArtifactId: "lombok", Scope: "provided"},
		},
	},
	"mapstruct": {
		Name:        "MapStruct",
		Description: "Type-safe bean mapping",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.mapstruct", ArtifactId: "mapstruct", Version: "1.5.5.Final"},
			{GroupId: "org.mapstruct", ArtifactId: "mapstruct-processor", Version: "1.5.5.Final", Scope: "provided"},
		},
	},
	"postgresql": {
		Name:        "PostgreSQL Driver",
		Description: "PostgreSQL JDBC driver",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.postgresql", ArtifactId: "postgresql", Scope: "runtime"},
		},
	},
	"mysql": {
		Name:        "MySQL Driver",
		Description: "MySQL JDBC driver",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "com.mysql", ArtifactId: "mysql-connector-j", Scope: "runtime"},
		},
	},
	"mariadb": {
		Name:        "MariaDB Driver",
		Description: "MariaDB JDBC driver",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.mariadb.jdbc", ArtifactId: "mariadb-java-client", Scope: "runtime"},
		},
	},
	"h2": {
		Name:        "H2 Database",
		Description: "In-memory database for development/testing",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "com.h2database", ArtifactId: "h2", Scope: "runtime"},
		},
	},
	"flyway": {
		Name:        "Flyway Migration",
		Description: "Database schema migrations",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.flywaydb", ArtifactId: "flyway-core"},
		},
	},
	"liquibase": {
		Name:        "Liquibase Migration",
		Description: "Database schema migrations",
		Category:    "SQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.liquibase", ArtifactId: "liquibase-core"},
		},
	},
	"mongodb": {
		Name:        "Spring Data MongoDB",
		Description: "MongoDB NoSQL database",
		Category:    "NoSQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-mongodb"},
		},
	},
	"redis": {
		Name:        "Spring Data Redis",
		Description: "Redis key-value store",
		Category:    "NoSQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-redis"},
		},
	},
	"elasticsearch": {
		Name:        "Spring Data Elasticsearch",
		Description: "Elasticsearch search engine",
		Category:    "NoSQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-elasticsearch"},
		},
	},
	"amqp": {
		Name:        "Spring AMQP",
		Description: "RabbitMQ messaging",
		Category:    "Messaging",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-amqp"},
		},
	},
	"kafka": {
		Name:        "Spring Kafka",
		Description: "Apache Kafka messaging",
		Category:    "Messaging",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.kafka", ArtifactId: "spring-kafka"},
		},
	},
	"mail": {
		Name:        "Java Mail",
		Description: "Send emails with Spring",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-mail"},
		},
	},
	"cache": {
		Name:        "Spring Cache",
		Description: "Caching abstraction",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-cache"},
		},
	},
	"thymeleaf": {
		Name:        "Thymeleaf",
		Description: "Server-side HTML templating",
		Category:    "Template Engines",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-thymeleaf"},
		},
	},
	"openapi": {
		Name:        "SpringDoc OpenAPI",
		Description: "OpenAPI 3 documentation (Swagger UI)",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springdoc", ArtifactId: "springdoc-openapi-starter-webmvc-ui", Version: "2.3.0"},
		},
	},
	"test": {
		Name:        "Spring Boot Test",
		Description: "Testing utilities for Spring Boot",
		Category:    "Testing",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-test", Scope: "test"},
		},
	},
	"testcontainers": {
		Name:        "Testcontainers",
		Description: "Integration testing with containers",
		Category:    "Testing",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-testcontainers", Scope: "test"},
			{GroupId: "org.testcontainers", ArtifactId: "junit-jupiter", Scope: "test"},
		},
	},
	"graphql": {
		Name:        "Spring GraphQL",
		Description: "GraphQL API support",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-graphql"},
		},
	},
	"websocket": {
		Name:        "WebSocket",
		Description: "WebSocket support",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-websocket"},
		},
	},
	"batch": {
		Name:        "Spring Batch",
		Description: "Batch processing framework",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-batch"},
		},
	},
	"quartz": {
		Name:        "Quartz Scheduler",
		Description: "Job scheduling",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-quartz"},
		},
	},
	"jwt": {
		Name:        "JJWT (JSON Web Token)",
		Description: "JWT creation and verification",
		Category:    "Security",
		Dependencies: []maven.Dependency{
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-api", Version: "0.12.5"},
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-impl", Version: "0.12.5", Scope: "runtime"},
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-jackson", Version: "0.12.5", Scope: "runtime"},
		},
	},
	"commons-lang": {
		Name:        "Apache Commons Lang",
		Description: "String manipulation and utilities",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.apache.commons", ArtifactId: "commons-lang3", Version: "3.14.0"},
		},
	},
	"commons-io": {
		Name:        "Apache Commons IO",
		Description: "IO utilities and file operations",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "commons-io", ArtifactId: "commons-io", Version: "2.15.1"},
		},
	},
	"guava": {
		Name:        "Google Guava",
		Description: "Core libraries for collections and utilities",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "com.google.guava", ArtifactId: "guava", Version: "33.0.0-jre"},
		},
	},
	"modelmapper": {
		Name:        "ModelMapper",
		Description: "Object mapping made simple",
		Category:    "Developer Tools",
		Dependencies: []maven.Dependency{
			{GroupId: "org.modelmapper", ArtifactId: "modelmapper", Version: "3.2.0"},
		},
	},
	"jackson-datatype": {
		Name:        "Jackson Java 8 Datatypes",
		Description: "Java 8 date/time support for Jackson",
		Category:    "I/O",
		Dependencies: []maven.Dependency{
			{GroupId: "com.fasterxml.jackson.datatype", ArtifactId: "jackson-datatype-jsr310"},
		},
	},
	"feign": {
		Name:        "Spring Cloud OpenFeign",
		Description: "Declarative REST client",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-openfeign"},
		},
	},
	"resilience4j": {
		Name:        "Resilience4j",
		Description: "Fault tolerance library (circuit breaker)",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "io.github.resilience4j", ArtifactId: "resilience4j-spring-boot3", Version: "2.2.0"},
		},
	},
	"micrometer": {
		Name:        "Micrometer Prometheus",
		Description: "Prometheus metrics exporter",
		Category:    "Ops",
		Dependencies: []maven.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-registry-prometheus"},
		},
	},
	"cassandra": {
		Name:        "Spring Data Cassandra",
		Description: "Apache Cassandra database",
		Category:    "NoSQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-cassandra"},
		},
	},
	"neo4j": {
		Name:        "Spring Data Neo4j",
		Description: "Neo4j graph database",
		Category:    "NoSQL",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-neo4j"},
		},
	},
	"freemarker": {
		Name:        "FreeMarker",
		Description: "FreeMarker template engine",
		Category:    "Template Engines",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-freemarker"},
		},
	},
	"mustache": {
		Name:        "Mustache",
		Description: "Mustache template engine",
		Category:    "Template Engines",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-mustache"},
		},
	},
	"security-test": {
		Name:        "Spring Security Test",
		Description: "Testing utilities for Spring Security",
		Category:    "Testing",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.security", ArtifactId: "spring-security-test", Scope: "test"},
		},
	},
	"mockito": {
		Name:        "Mockito",
		Description: "Mocking framework for unit tests",
		Category:    "Testing",
		Dependencies: []maven.Dependency{
			{GroupId: "org.mockito", ArtifactId: "mockito-core", Scope: "test"},
		},
	},
	"restdocs": {
		Name:        "Spring REST Docs",
		Description: "Generate API documentation from tests",
		Category:    "Testing",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.restdocs", ArtifactId: "spring-restdocs-mockmvc", Scope: "test"},
		},
	},
	"hateoas": {
		Name:        "Spring HATEOAS",
		Description: "Hypermedia-driven REST APIs",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-hateoas"},
		},
	},
	"data-rest": {
		Name:        "Spring Data REST",
		Description: "Expose repositories as REST endpoints",
		Category:    "Web",
		Dependencies: []maven.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-rest"},
		},
	},
}

func GetCatalogEntry(alias string) (CatalogEntry, bool) {
	entry, ok := dependencyCatalog[alias]
	return entry, ok
}

func GetAllAliases() []string {
	aliases := make([]string, 0, len(dependencyCatalog))
	for alias := range dependencyCatalog {
		aliases = append(aliases, alias)
	}
	return aliases
}

func GetCatalogByCategory() map[string][]string {
	categories := make(map[string][]string)
	for alias, entry := range dependencyCatalog {
		categories[entry.Category] = append(categories[entry.Category], alias)
	}
	return categories
}

func SearchCatalog(query string) []CatalogEntry {
	var results []CatalogEntry
	queryLower := toLower(query)

	for _, entry := range dependencyCatalog {
		if contains(toLower(entry.Name), queryLower) ||
			contains(toLower(entry.Description), queryLower) ||
			contains(toLower(entry.Category), queryLower) {
			results = append(results, entry)
		}
	}
	return results
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	return len(substr) <= len(s) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
