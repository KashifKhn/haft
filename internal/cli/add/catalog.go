package add

import "github.com/KashifKhn/haft/internal/buildtool"

type CatalogEntry struct {
	Name         string
	Description  string
	Category     string
	Dependencies []buildtool.Dependency
}

var dependencyCatalog = map[string]CatalogEntry{
	"web": {
		Name:        "Spring Web",
		Description: "Build web applications with Spring MVC",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		},
	},
	"webflux": {
		Name:        "Spring WebFlux",
		Description: "Build reactive web applications",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-webflux"},
		},
	},
	"jpa": {
		Name:        "Spring Data JPA",
		Description: "Persist data with JPA and Hibernate",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-jpa"},
		},
	},
	"jdbc": {
		Name:        "Spring JDBC",
		Description: "Database access with JDBC",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-jdbc"},
		},
	},
	"security": {
		Name:        "Spring Security",
		Description: "Authentication and authorization",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-security"},
		},
	},
	"oauth2-client": {
		Name:        "OAuth2 Client",
		Description: "OAuth2/OpenID Connect client",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-oauth2-client"},
		},
	},
	"oauth2-resource-server": {
		Name:        "OAuth2 Resource Server",
		Description: "OAuth2 resource server support",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-oauth2-resource-server"},
		},
	},
	"validation": {
		Name:        "Validation",
		Description: "Bean validation with Hibernate Validator",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-validation"},
		},
	},
	"actuator": {
		Name:        "Spring Boot Actuator",
		Description: "Production-ready monitoring and management",
		Category:    "Ops",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-actuator"},
		},
	},
	"devtools": {
		Name:        "Spring Boot DevTools",
		Description: "Fast application restarts and LiveReload",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-devtools", Scope: "runtime", Optional: true},
		},
	},
	"lombok": {
		Name:        "Lombok",
		Description: "Reduce boilerplate with annotations",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.projectlombok", ArtifactId: "lombok", Scope: "provided"},
		},
	},
	"mapstruct": {
		Name:        "MapStruct",
		Description: "Type-safe bean mapping",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.mapstruct", ArtifactId: "mapstruct", Version: "1.5.5.Final"},
			{GroupId: "org.mapstruct", ArtifactId: "mapstruct-processor", Version: "1.5.5.Final", Scope: "provided"},
		},
	},
	"postgresql": {
		Name:        "PostgreSQL Driver",
		Description: "PostgreSQL JDBC driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.postgresql", ArtifactId: "postgresql", Scope: "runtime"},
		},
	},
	"mysql": {
		Name:        "MySQL Driver",
		Description: "MySQL JDBC driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.mysql", ArtifactId: "mysql-connector-j", Scope: "runtime"},
		},
	},
	"mariadb": {
		Name:        "MariaDB Driver",
		Description: "MariaDB JDBC driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.mariadb.jdbc", ArtifactId: "mariadb-java-client", Scope: "runtime"},
		},
	},
	"h2": {
		Name:        "H2 Database",
		Description: "In-memory database for development/testing",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.h2database", ArtifactId: "h2", Scope: "runtime"},
		},
	},
	"flyway": {
		Name:        "Flyway Migration",
		Description: "Database schema migrations",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.flywaydb", ArtifactId: "flyway-core"},
		},
	},
	"liquibase": {
		Name:        "Liquibase Migration",
		Description: "Database schema migrations",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.liquibase", ArtifactId: "liquibase-core"},
		},
	},
	"mongodb": {
		Name:        "Spring Data MongoDB",
		Description: "MongoDB NoSQL database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-mongodb"},
		},
	},
	"redis": {
		Name:        "Spring Data Redis",
		Description: "Redis key-value store",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-redis"},
		},
	},
	"elasticsearch": {
		Name:        "Spring Data Elasticsearch",
		Description: "Elasticsearch search engine",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-elasticsearch"},
		},
	},
	"amqp": {
		Name:        "Spring AMQP",
		Description: "RabbitMQ messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-amqp"},
		},
	},
	"kafka": {
		Name:        "Spring Kafka",
		Description: "Apache Kafka messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.kafka", ArtifactId: "spring-kafka"},
		},
	},
	"mail": {
		Name:        "Java Mail",
		Description: "Send emails with Spring",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-mail"},
		},
	},
	"cache": {
		Name:        "Spring Cache",
		Description: "Caching abstraction",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-cache"},
		},
	},
	"thymeleaf": {
		Name:        "Thymeleaf",
		Description: "Server-side HTML templating",
		Category:    "Template Engines",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-thymeleaf"},
		},
	},
	"openapi": {
		Name:        "SpringDoc OpenAPI",
		Description: "OpenAPI 3 documentation (Swagger UI)",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springdoc", ArtifactId: "springdoc-openapi-starter-webmvc-ui", Version: "2.3.0"},
		},
	},
	"test": {
		Name:        "Spring Boot Test",
		Description: "Testing utilities for Spring Boot",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-test", Scope: "test"},
		},
	},
	"testcontainers": {
		Name:        "Testcontainers",
		Description: "Integration testing with containers",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-testcontainers", Scope: "test"},
			{GroupId: "org.testcontainers", ArtifactId: "junit-jupiter", Scope: "test"},
		},
	},
	"graphql": {
		Name:        "Spring GraphQL",
		Description: "GraphQL API support",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-graphql"},
		},
	},
	"websocket": {
		Name:        "WebSocket",
		Description: "WebSocket support",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-websocket"},
		},
	},
	"batch": {
		Name:        "Spring Batch",
		Description: "Batch processing framework",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-batch"},
		},
	},
	"quartz": {
		Name:        "Quartz Scheduler",
		Description: "Job scheduling",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-quartz"},
		},
	},
	"jwt": {
		Name:        "JJWT (JSON Web Token)",
		Description: "JWT creation and verification",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-api", Version: "0.12.5"},
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-impl", Version: "0.12.5", Scope: "runtime"},
			{GroupId: "io.jsonwebtoken", ArtifactId: "jjwt-jackson", Version: "0.12.5", Scope: "runtime"},
		},
	},
	"commons-lang": {
		Name:        "Apache Commons Lang",
		Description: "String manipulation and utilities",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.commons", ArtifactId: "commons-lang3", Version: "3.14.0"},
		},
	},
	"commons-io": {
		Name:        "Apache Commons IO",
		Description: "IO utilities and file operations",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "commons-io", ArtifactId: "commons-io", Version: "2.15.1"},
		},
	},
	"guava": {
		Name:        "Google Guava",
		Description: "Core libraries for collections and utilities",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.guava", ArtifactId: "guava", Version: "33.0.0-jre"},
		},
	},
	"modelmapper": {
		Name:        "ModelMapper",
		Description: "Object mapping made simple",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.modelmapper", ArtifactId: "modelmapper", Version: "3.2.0"},
		},
	},
	"jackson-datatype": {
		Name:        "Jackson Java 8 Datatypes",
		Description: "Java 8 date/time support for Jackson",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.fasterxml.jackson.datatype", ArtifactId: "jackson-datatype-jsr310"},
		},
	},
	"feign": {
		Name:        "Spring Cloud OpenFeign",
		Description: "Declarative REST client",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-openfeign"},
		},
	},
	"resilience4j": {
		Name:        "Resilience4j",
		Description: "Fault tolerance library (circuit breaker)",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.github.resilience4j", ArtifactId: "resilience4j-spring-boot3", Version: "2.2.0"},
		},
	},
	"micrometer": {
		Name:        "Micrometer Prometheus",
		Description: "Prometheus metrics exporter",
		Category:    "Ops",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-registry-prometheus"},
		},
	},
	"cassandra": {
		Name:        "Spring Data Cassandra",
		Description: "Apache Cassandra database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-cassandra"},
		},
	},
	"neo4j": {
		Name:        "Spring Data Neo4j",
		Description: "Neo4j graph database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-neo4j"},
		},
	},
	"freemarker": {
		Name:        "FreeMarker",
		Description: "FreeMarker template engine",
		Category:    "Template Engines",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-freemarker"},
		},
	},
	"mustache": {
		Name:        "Mustache",
		Description: "Mustache template engine",
		Category:    "Template Engines",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-mustache"},
		},
	},
	"security-test": {
		Name:        "Spring Security Test",
		Description: "Testing utilities for Spring Security",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.security", ArtifactId: "spring-security-test", Scope: "test"},
		},
	},
	"mockito": {
		Name:        "Mockito",
		Description: "Mocking framework for unit tests",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.mockito", ArtifactId: "mockito-core", Scope: "test"},
		},
	},
	"restdocs": {
		Name:        "Spring REST Docs",
		Description: "Generate API documentation from tests",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.restdocs", ArtifactId: "spring-restdocs-mockmvc", Scope: "test"},
		},
	},
	"hateoas": {
		Name:        "Spring HATEOAS",
		Description: "Hypermedia-driven REST APIs",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-hateoas"},
		},
	},
	"data-rest": {
		Name:        "Spring Data REST",
		Description: "Expose repositories as REST endpoints",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-rest"},
		},
	},

	// AI Dependencies
	"openai": {
		Name:        "Spring AI OpenAI",
		Description: "OpenAI (ChatGPT, GPT-4) integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-openai-spring-boot-starter"},
		},
	},
	"anthropic": {
		Name:        "Spring AI Anthropic",
		Description: "Anthropic Claude integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-anthropic-spring-boot-starter"},
		},
	},
	"ollama": {
		Name:        "Spring AI Ollama",
		Description: "Ollama local LLM integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-ollama-spring-boot-starter"},
		},
	},
	"azure-openai": {
		Name:        "Spring AI Azure OpenAI",
		Description: "Azure OpenAI integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-azure-openai-spring-boot-starter"},
		},
	},
	"bedrock": {
		Name:        "Spring AI Amazon Bedrock",
		Description: "Amazon Bedrock AI integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-bedrock-ai-spring-boot-starter"},
		},
	},
	"vertex-ai": {
		Name:        "Spring AI Vertex AI",
		Description: "Google Vertex AI Gemini integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-vertex-ai-gemini-spring-boot-starter"},
		},
	},
	"mistral": {
		Name:        "Spring AI Mistral",
		Description: "Mistral AI integration",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-mistral-ai-spring-boot-starter"},
		},
	},
	"pgvector": {
		Name:        "PGVector",
		Description: "PostgreSQL vector database for AI",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.ai", ArtifactId: "spring-ai-pgvector-store-spring-boot-starter"},
		},
	},

	// Security - Additional
	"oauth2-authorization-server": {
		Name:        "OAuth2 Authorization Server",
		Description: "Build your own OAuth2 authorization server",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-oauth2-authorization-server"},
		},
	},
	"ldap": {
		Name:        "Spring LDAP",
		Description: "LDAP authentication and operations",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-ldap"},
		},
	},

	// Payments and External Services
	"stripe": {
		Name:        "Stripe",
		Description: "Stripe payment processing",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.stripe", ArtifactId: "stripe-java", Version: "26.0.0"},
		},
	},
	"paypal": {
		Name:        "PayPal",
		Description: "PayPal payment SDK",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.paypal.sdk", ArtifactId: "checkout-sdk", Version: "2.0.0"},
		},
	},

	// Cloud - AWS
	"aws-s3": {
		Name:        "AWS S3",
		Description: "Amazon S3 file storage",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "s3", Version: "2.25.0"},
		},
	},
	"aws-sqs": {
		Name:        "AWS SQS",
		Description: "Amazon Simple Queue Service",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "sqs", Version: "2.25.0"},
		},
	},
	"aws-ses": {
		Name:        "AWS SES",
		Description: "Amazon Simple Email Service",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "ses", Version: "2.25.0"},
		},
	},
	"aws-dynamodb": {
		Name:        "AWS DynamoDB",
		Description: "Amazon DynamoDB NoSQL database",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "dynamodb-enhanced", Version: "2.25.0"},
		},
	},

	// Observability
	"prometheus": {
		Name:        "Prometheus",
		Description: "Prometheus metrics registry",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-registry-prometheus"},
		},
	},
	"zipkin": {
		Name:        "Zipkin",
		Description: "Distributed tracing with Zipkin",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.zipkin.reporter2", ArtifactId: "zipkin-reporter-brave"},
		},
	},
	"opentelemetry": {
		Name:        "OpenTelemetry",
		Description: "OpenTelemetry distributed tracing",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-tracing-bridge-otel"},
			{GroupId: "io.opentelemetry", ArtifactId: "opentelemetry-exporter-otlp"},
		},
	},

	// Additional useful libraries
	"r2dbc": {
		Name:        "Spring Data R2DBC",
		Description: "Reactive database access",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-r2dbc"},
		},
	},
	"r2dbc-postgresql": {
		Name:        "R2DBC PostgreSQL",
		Description: "Reactive PostgreSQL driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.postgresql", ArtifactId: "r2dbc-postgresql", Scope: "runtime"},
		},
	},
	"mybatis": {
		Name:        "MyBatis",
		Description: "MyBatis SQL mapping framework",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.mybatis.spring.boot", ArtifactId: "mybatis-spring-boot-starter", Version: "3.0.3"},
		},
	},
	"jooq": {
		Name:        "jOOQ",
		Description: "Type-safe SQL query builder",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-jooq"},
		},
	},
	"oracle": {
		Name:        "Oracle Driver",
		Description: "Oracle database JDBC driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.oracle.database.jdbc", ArtifactId: "ojdbc11", Scope: "runtime"},
		},
	},
	"sqlserver": {
		Name:        "SQL Server Driver",
		Description: "Microsoft SQL Server JDBC driver",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.microsoft.sqlserver", ArtifactId: "mssql-jdbc", Scope: "runtime"},
		},
	},

	// Reactive MongoDB
	"mongodb-reactive": {
		Name:        "Spring Data MongoDB Reactive",
		Description: "Reactive MongoDB support",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-mongodb-reactive"},
		},
	},
	"couchbase": {
		Name:        "Spring Data Couchbase",
		Description: "Couchbase NoSQL database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-data-couchbase"},
		},
	},

	// Messaging - Additional
	"pulsar": {
		Name:        "Spring Pulsar",
		Description: "Apache Pulsar messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-pulsar"},
		},
	},
	"activemq": {
		Name:        "Apache ActiveMQ",
		Description: "ActiveMQ JMS messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-activemq"},
		},
	},
	"rsocket": {
		Name:        "RSocket",
		Description: "RSocket reactive messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-rsocket"},
		},
	},

	// Session management
	"session-redis": {
		Name:        "Spring Session Redis",
		Description: "Distributed sessions with Redis",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.session", ArtifactId: "spring-session-data-redis"},
		},
	},
	"session-jdbc": {
		Name:        "Spring Session JDBC",
		Description: "Distributed sessions with JDBC",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.session", ArtifactId: "spring-session-jdbc"},
		},
	},

	// HTTP Clients
	"restclient": {
		Name:        "Spring RestClient",
		Description: "Synchronous HTTP client",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-web"},
		},
	},
	"webclient": {
		Name:        "Spring WebClient",
		Description: "Reactive HTTP client",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-webflux"},
		},
	},

	// File handling
	"apache-poi": {
		Name:        "Apache POI",
		Description: "Microsoft Office file handling (Excel, Word)",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.poi", ArtifactId: "poi-ooxml", Version: "5.2.5"},
		},
	},
	"itext": {
		Name:        "iText PDF",
		Description: "PDF generation and manipulation",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.itextpdf", ArtifactId: "itext7-core", Version: "8.0.2"},
		},
	},
	"minio": {
		Name:        "MinIO",
		Description: "S3-compatible object storage client",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.minio", ArtifactId: "minio", Version: "8.5.7"},
		},
	},

	// Scheduler
	"scheduler": {
		Name:        "Spring Scheduler",
		Description: "Task scheduling with @Scheduled",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter"},
		},
	},

	// Validation extras
	"passay": {
		Name:        "Passay",
		Description: "Password validation library",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.passay", ArtifactId: "passay", Version: "1.6.4"},
		},
	},

	// JSON processing
	"json-path": {
		Name:        "JsonPath",
		Description: "JSON path query library",
		Category:    "I/O",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.jayway.jsonpath", ArtifactId: "json-path", Version: "2.9.0"},
		},
	},

	// gRPC
	"grpc": {
		Name:        "gRPC",
		Description: "gRPC framework for Spring Boot",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.devh", ArtifactId: "grpc-spring-boot-starter", Version: "3.0.0.RELEASE"},
		},
	},

	// Rate limiting
	"bucket4j": {
		Name:        "Bucket4j",
		Description: "Rate limiting library",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.bucket4j", ArtifactId: "bucket4j-core", Version: "8.7.0"},
		},
	},

	// Configuration
	"config-processor": {
		Name:        "Configuration Processor",
		Description: "IDE support for custom @ConfigurationProperties",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-configuration-processor", Scope: "provided"},
		},
	},
	"docker-compose": {
		Name:        "Docker Compose Support",
		Description: "Docker Compose development support",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-docker-compose", Scope: "runtime"},
		},
	},

	// Native
	"native": {
		Name:        "GraalVM Native",
		Description: "GraalVM native image support",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.graalvm.buildtools", ArtifactId: "native-maven-plugin"},
		},
	},

	// ============================================
	// COMMUNICATION / NOTIFICATIONS
	// ============================================
	"twilio": {
		Name:        "Twilio",
		Description: "SMS, WhatsApp, Voice API",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.twilio.sdk", ArtifactId: "twilio", Version: "10.1.0"},
		},
	},
	"sendgrid": {
		Name:        "SendGrid",
		Description: "Email delivery service",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.sendgrid", ArtifactId: "sendgrid-java", Version: "4.10.2"},
		},
	},
	"mailersend": {
		Name:        "MailerSend",
		Description: "Transactional email service",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.mailersend", ArtifactId: "java-sdk", Version: "1.0.0"},
		},
	},
	"firebase-admin": {
		Name:        "Firebase Admin",
		Description: "Firebase push notifications and services",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.firebase", ArtifactId: "firebase-admin", Version: "9.2.0"},
		},
	},
	"pusher": {
		Name:        "Pusher",
		Description: "Real-time messaging and websockets",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.pusher", ArtifactId: "pusher-http-java", Version: "1.3.3"},
		},
	},
	"slack": {
		Name:        "Slack SDK",
		Description: "Slack API integration",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.slack.api", ArtifactId: "slack-api-client", Version: "1.38.0"},
		},
	},
	"discord": {
		Name:        "Discord JDA",
		Description: "Discord bot and API integration",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.dv8tion", ArtifactId: "JDA", Version: "5.0.0-beta.21"},
		},
	},
	"telegram": {
		Name:        "Telegram Bot",
		Description: "Telegram bot API",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.telegram", ArtifactId: "telegrambots", Version: "6.9.7.1"},
		},
	},
	"mailgun": {
		Name:        "Mailgun",
		Description: "Email delivery via Mailgun",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.sargue", ArtifactId: "mailgun", Version: "1.10.0"},
		},
	},
	"onesignal": {
		Name:        "OneSignal",
		Description: "Push notification service",
		Category:    "Notifications",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.onesignal", ArtifactId: "onesignal-java-client", Version: "2.0.2"},
		},
	},

	// ============================================
	// CLOUD - AWS (Additional)
	// ============================================
	"aws-sns": {
		Name:        "AWS SNS",
		Description: "Amazon Simple Notification Service",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "sns", Version: "2.25.0"},
		},
	},
	"aws-lambda": {
		Name:        "AWS Lambda",
		Description: "Amazon Lambda serverless functions",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.amazonaws", ArtifactId: "aws-lambda-java-core", Version: "1.2.3"},
			{GroupId: "com.amazonaws", ArtifactId: "aws-lambda-java-events", Version: "3.11.4"},
		},
	},
	"aws-cognito": {
		Name:        "AWS Cognito",
		Description: "Amazon Cognito user authentication",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "cognitoidentityprovider", Version: "2.25.0"},
		},
	},
	"aws-secretsmanager": {
		Name:        "AWS Secrets Manager",
		Description: "Amazon Secrets Manager",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "secretsmanager", Version: "2.25.0"},
		},
	},
	"aws-cloudwatch": {
		Name:        "AWS CloudWatch",
		Description: "Amazon CloudWatch monitoring",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "cloudwatch", Version: "2.25.0"},
		},
	},
	"aws-kinesis": {
		Name:        "AWS Kinesis",
		Description: "Amazon Kinesis data streaming",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "kinesis", Version: "2.25.0"},
		},
	},

	// CLOUD - GCP
	"gcp-storage": {
		Name:        "Google Cloud Storage",
		Description: "GCP object storage",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud", ArtifactId: "google-cloud-storage", Version: "2.35.0"},
		},
	},
	"gcp-pubsub": {
		Name:        "Google Cloud Pub/Sub",
		Description: "GCP messaging service",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud", ArtifactId: "google-cloud-pubsub", Version: "1.127.0"},
		},
	},
	"gcp-bigquery": {
		Name:        "Google BigQuery",
		Description: "GCP data warehouse",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud", ArtifactId: "google-cloud-bigquery", Version: "2.38.1"},
		},
	},
	"gcp-firestore": {
		Name:        "Google Firestore",
		Description: "GCP Firestore NoSQL database",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud", ArtifactId: "google-cloud-firestore", Version: "3.18.0"},
		},
	},
	"gcp-secretmanager": {
		Name:        "Google Secret Manager",
		Description: "GCP secrets management",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud", ArtifactId: "google-cloud-secretmanager", Version: "2.37.0"},
		},
	},
	"gcp-functions": {
		Name:        "Google Cloud Functions",
		Description: "GCP serverless functions",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud.functions", ArtifactId: "functions-framework-api", Version: "1.1.0"},
		},
	},

	// CLOUD - Azure
	"azure-storage": {
		Name:        "Azure Blob Storage",
		Description: "Azure object storage",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.azure", ArtifactId: "azure-storage-blob", Version: "12.25.2"},
		},
	},
	"azure-servicebus": {
		Name:        "Azure Service Bus",
		Description: "Azure messaging service",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.azure", ArtifactId: "azure-messaging-servicebus", Version: "7.15.1"},
		},
	},
	"azure-keyvault": {
		Name:        "Azure Key Vault",
		Description: "Azure secrets management",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.azure", ArtifactId: "azure-security-keyvault-secrets", Version: "4.8.0"},
		},
	},
	"azure-cosmosdb": {
		Name:        "Azure Cosmos DB",
		Description: "Azure NoSQL database",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.azure", ArtifactId: "azure-cosmos", Version: "4.55.0"},
		},
	},
	"azure-functions": {
		Name:        "Azure Functions",
		Description: "Azure serverless functions",
		Category:    "Cloud",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.microsoft.azure.functions", ArtifactId: "azure-functions-java-library", Version: "3.0.0"},
		},
	},

	// ============================================
	// SEARCH
	// ============================================
	"algolia": {
		Name:        "Algolia",
		Description: "Search-as-a-service",
		Category:    "Search",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.algolia", ArtifactId: "algoliasearch", Version: "3.16.9"},
		},
	},
	"meilisearch": {
		Name:        "Meilisearch",
		Description: "Fast open-source search engine",
		Category:    "Search",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.meilisearch.sdk", ArtifactId: "meilisearch-java", Version: "0.11.8"},
		},
	},
	"typesense": {
		Name:        "Typesense",
		Description: "Open-source typo-tolerant search",
		Category:    "Search",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.typesense", ArtifactId: "typesense-java", Version: "0.5.0"},
		},
	},
	"solr": {
		Name:        "Apache Solr",
		Description: "Enterprise search platform",
		Category:    "Search",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.solr", ArtifactId: "solr-solrj", Version: "9.5.0"},
		},
	},
	"opensearch": {
		Name:        "OpenSearch",
		Description: "OpenSearch client (Elasticsearch fork)",
		Category:    "Search",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.opensearch.client", ArtifactId: "opensearch-rest-high-level-client", Version: "2.12.0"},
		},
	},

	// ============================================
	// DATABASES - Additional NoSQL / TimeSeries
	// ============================================
	"scylladb": {
		Name:        "ScyllaDB",
		Description: "High-performance Cassandra-compatible DB",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.scylladb", ArtifactId: "java-driver-core", Version: "4.17.0.0"},
		},
	},
	"cockroachdb": {
		Name:        "CockroachDB",
		Description: "Distributed SQL database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.postgresql", ArtifactId: "postgresql", Scope: "runtime"},
		},
	},
	"influxdb": {
		Name:        "InfluxDB",
		Description: "Time-series database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.influxdb", ArtifactId: "influxdb-client-java", Version: "7.0.0"},
		},
	},
	"timescaledb": {
		Name:        "TimescaleDB",
		Description: "PostgreSQL-based time-series database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.postgresql", ArtifactId: "postgresql", Scope: "runtime"},
		},
	},
	"clickhouse": {
		Name:        "ClickHouse",
		Description: "Column-oriented OLAP database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.clickhouse", ArtifactId: "clickhouse-jdbc", Version: "0.6.0"},
		},
	},
	"dynamodb-local": {
		Name:        "DynamoDB Enhanced",
		Description: "DynamoDB enhanced client with mapping",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "software.amazon.awssdk", ArtifactId: "dynamodb-enhanced", Version: "2.25.0"},
		},
	},
	"arangodb": {
		Name:        "ArangoDB",
		Description: "Multi-model graph database",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.arangodb", ArtifactId: "arangodb-java-driver", Version: "7.5.0"},
		},
	},
	"voltdb": {
		Name:        "VoltDB",
		Description: "In-memory ACID SQL database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.voltdb", ArtifactId: "voltdbclient", Version: "13.0"},
		},
	},
	"hazelcast": {
		Name:        "Hazelcast",
		Description: "In-memory data grid",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.hazelcast", ArtifactId: "hazelcast-spring", Version: "5.3.6"},
		},
	},
	"ignite": {
		Name:        "Apache Ignite",
		Description: "Distributed in-memory computing",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.ignite", ArtifactId: "ignite-spring-boot-autoconfigure-ext", Version: "1.0.0"},
		},
	},
	"memcached": {
		Name:        "Memcached",
		Description: "Distributed memory caching",
		Category:    "NoSQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.code.simple-spring-memcached", ArtifactId: "spymemcached-provider", Version: "4.1.3"},
		},
	},

	// ============================================
	// UTILITIES
	// ============================================
	"okhttp": {
		Name:        "OkHttp",
		Description: "Efficient HTTP client",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.squareup.okhttp3", ArtifactId: "okhttp", Version: "4.12.0"},
		},
	},
	"retrofit": {
		Name:        "Retrofit",
		Description: "Type-safe REST client",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.squareup.retrofit2", ArtifactId: "retrofit", Version: "2.9.0"},
			{GroupId: "com.squareup.retrofit2", ArtifactId: "converter-jackson", Version: "2.9.0"},
		},
	},
	"jsoup": {
		Name:        "Jsoup",
		Description: "HTML parsing and web scraping",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jsoup", ArtifactId: "jsoup", Version: "1.17.2"},
		},
	},
	"zxing": {
		Name:        "ZXing",
		Description: "QR code and barcode generation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.zxing", ArtifactId: "core", Version: "3.5.3"},
			{GroupId: "com.google.zxing", ArtifactId: "javase", Version: "3.5.3"},
		},
	},
	"thumbnailator": {
		Name:        "Thumbnailator",
		Description: "Image resizing and processing",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.coobird", ArtifactId: "thumbnailator", Version: "0.4.20"},
		},
	},
	"jasperreports": {
		Name:        "JasperReports",
		Description: "Report generation library",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.sf.jasperreports", ArtifactId: "jasperreports", Version: "6.21.2"},
		},
	},
	"jfreechart": {
		Name:        "JFreeChart",
		Description: "Chart and graph generation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jfree", ArtifactId: "jfreechart", Version: "1.5.4"},
		},
	},
	"opencsv": {
		Name:        "OpenCSV",
		Description: "CSV parsing and writing",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.opencsv", ArtifactId: "opencsv", Version: "5.9"},
		},
	},
	"snakeyaml": {
		Name:        "SnakeYAML",
		Description: "YAML parsing and writing",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.yaml", ArtifactId: "snakeyaml", Version: "2.2"},
		},
	},
	"jna": {
		Name:        "JNA",
		Description: "Java Native Access",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.java.dev.jna", ArtifactId: "jna", Version: "5.14.0"},
		},
	},
	"joda-time": {
		Name:        "Joda-Time",
		Description: "Date and time library (legacy)",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "joda-time", ArtifactId: "joda-time", Version: "2.12.7"},
		},
	},
	"slug": {
		Name:        "Slugify",
		Description: "URL slug generation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.slugify", ArtifactId: "slugify", Version: "3.0.6"},
		},
	},
	"libphonenumber": {
		Name:        "LibPhoneNumber",
		Description: "Phone number parsing and validation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.googlecode.libphonenumber", ArtifactId: "libphonenumber", Version: "8.13.30"},
		},
	},
	"emoji": {
		Name:        "Emoji Java",
		Description: "Emoji parsing and manipulation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.vdurmont", ArtifactId: "emoji-java", Version: "5.1.1"},
		},
	},
	"commonmark": {
		Name:        "CommonMark",
		Description: "Markdown parsing and rendering",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.commonmark", ArtifactId: "commonmark", Version: "0.21.0"},
		},
	},
	"flexmark": {
		Name:        "Flexmark",
		Description: "Advanced Markdown processor",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.vladsch.flexmark", ArtifactId: "flexmark-all", Version: "0.64.8"},
		},
	},
	"httpclient5": {
		Name:        "Apache HttpClient 5",
		Description: "Apache HTTP client",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.httpcomponents.client5", ArtifactId: "httpclient5", Version: "5.3.1"},
		},
	},
	"jsch": {
		Name:        "JSch",
		Description: "SSH and SFTP client",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.mwiede", ArtifactId: "jsch", Version: "0.2.16"},
		},
	},
	"sshj": {
		Name:        "SSHJ",
		Description: "Modern SSH client library",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.hierynomus", ArtifactId: "sshj", Version: "0.38.0"},
		},
	},
	"jaxb": {
		Name:        "JAXB",
		Description: "XML binding (Java 11+)",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "jakarta.xml.bind", ArtifactId: "jakarta.xml.bind-api", Version: "4.0.1"},
			{GroupId: "org.glassfish.jaxb", ArtifactId: "jaxb-runtime", Version: "4.0.4"},
		},
	},
	"dom4j": {
		Name:        "Dom4j",
		Description: "XML parsing and manipulation",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.dom4j", ArtifactId: "dom4j", Version: "2.1.4"},
		},
	},
	"xstream": {
		Name:        "XStream",
		Description: "XML serialization",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.thoughtworks.xstream", ArtifactId: "xstream", Version: "1.4.20"},
		},
	},
	"protobuf": {
		Name:        "Protocol Buffers",
		Description: "Google Protocol Buffers",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.protobuf", ArtifactId: "protobuf-java", Version: "3.25.3"},
		},
	},
	"avro": {
		Name:        "Apache Avro",
		Description: "Data serialization system",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.avro", ArtifactId: "avro", Version: "1.11.3"},
		},
	},
	"msgpack": {
		Name:        "MessagePack",
		Description: "Binary serialization format",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.msgpack", ArtifactId: "msgpack-core", Version: "0.9.8"},
		},
	},
	"kryo": {
		Name:        "Kryo",
		Description: "Fast binary serialization",
		Category:    "Utilities",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.esotericsoftware", ArtifactId: "kryo", Version: "5.6.0"},
		},
	},

	// ============================================
	// ADDITIONAL MESSAGING
	// ============================================
	"rabbitmq": {
		Name:        "RabbitMQ Client",
		Description: "RabbitMQ Java client (direct)",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.rabbitmq", ArtifactId: "amqp-client", Version: "5.20.0"},
		},
	},
	"nats": {
		Name:        "NATS",
		Description: "NATS messaging system",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.nats", ArtifactId: "jnats", Version: "2.17.3"},
		},
	},
	"zeromq": {
		Name:        "ZeroMQ",
		Description: "High-performance messaging",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.zeromq", ArtifactId: "jeromq", Version: "0.6.0"},
		},
	},
	"artemis": {
		Name:        "Apache Artemis",
		Description: "Apache ActiveMQ Artemis",
		Category:    "Messaging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-artemis"},
		},
	},

	// ============================================
	// ADDITIONAL SECURITY
	// ============================================
	"bouncy-castle": {
		Name:        "Bouncy Castle",
		Description: "Cryptography provider",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.bouncycastle", ArtifactId: "bcprov-jdk18on", Version: "1.77"},
		},
	},
	"jasypt": {
		Name:        "Jasypt",
		Description: "Encryption for application properties",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.ulisesbocchio", ArtifactId: "jasypt-spring-boot-starter", Version: "3.0.5"},
		},
	},
	"keycloak": {
		Name:        "Keycloak",
		Description: "Keycloak admin client",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.keycloak", ArtifactId: "keycloak-admin-client", Version: "24.0.1"},
		},
	},
	"auth0": {
		Name:        "Auth0",
		Description: "Auth0 authentication",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.auth0", ArtifactId: "auth0-spring-security-api", Version: "1.5.3"},
		},
	},
	"nimbus-jose": {
		Name:        "Nimbus JOSE+JWT",
		Description: "JWT and JOSE implementation",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.nimbusds", ArtifactId: "nimbus-jose-jwt", Version: "9.37.3"},
		},
	},
	"otp": {
		Name:        "Google Authenticator",
		Description: "TOTP/HOTP for 2FA",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "dev.samstevens.totp", ArtifactId: "totp", Version: "1.7.1"},
		},
	},

	// ============================================
	// ADDITIONAL TESTING
	// ============================================
	"wiremock": {
		Name:        "WireMock",
		Description: "HTTP mock server for testing",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.wiremock", ArtifactId: "wiremock-standalone", Version: "3.4.2", Scope: "test"},
		},
	},
	"awaitility": {
		Name:        "Awaitility",
		Description: "Async testing utility",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.awaitility", ArtifactId: "awaitility", Version: "4.2.0", Scope: "test"},
		},
	},
	"archunit": {
		Name:        "ArchUnit",
		Description: "Architecture testing",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.tngtech.archunit", ArtifactId: "archunit-junit5", Version: "1.2.1", Scope: "test"},
		},
	},
	"assertj": {
		Name:        "AssertJ",
		Description: "Fluent assertions library",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.assertj", ArtifactId: "assertj-core", Version: "3.25.3", Scope: "test"},
		},
	},
	"jsonassert": {
		Name:        "JSONAssert",
		Description: "JSON comparison for tests",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.skyscreamer", ArtifactId: "jsonassert", Version: "1.5.1", Scope: "test"},
		},
	},
	"rest-assured": {
		Name:        "REST Assured",
		Description: "REST API testing",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.rest-assured", ArtifactId: "rest-assured", Version: "5.4.0", Scope: "test"},
		},
	},
	"gatling": {
		Name:        "Gatling",
		Description: "Load testing framework",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.gatling.highcharts", ArtifactId: "gatling-charts-highcharts", Version: "3.10.4", Scope: "test"},
		},
	},
	"jmh": {
		Name:        "JMH",
		Description: "Java Microbenchmark Harness",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.openjdk.jmh", ArtifactId: "jmh-core", Version: "1.37", Scope: "test"},
			{GroupId: "org.openjdk.jmh", ArtifactId: "jmh-generator-annprocess", Version: "1.37", Scope: "test"},
		},
	},
	"faker": {
		Name:        "DataFaker",
		Description: "Generate fake test data",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.datafaker", ArtifactId: "datafaker", Version: "2.1.0", Scope: "test"},
		},
	},
	"greenmail": {
		Name:        "GreenMail",
		Description: "Email testing server",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.icegreen", ArtifactId: "greenmail-junit5", Version: "2.0.1", Scope: "test"},
		},
	},

	// ============================================
	// ADDITIONAL DEVELOPER TOOLS
	// ============================================
	"spotbugs": {
		Name:        "SpotBugs",
		Description: "Static analysis annotations",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.spotbugs", ArtifactId: "spotbugs-annotations", Version: "4.8.3"},
		},
	},
	"error-prone": {
		Name:        "Error Prone",
		Description: "Static analysis annotations",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.errorprone", ArtifactId: "error_prone_annotations", Version: "2.25.0"},
		},
	},
	"checker-qual": {
		Name:        "Checker Framework",
		Description: "Type annotations for static analysis",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.checkerframework", ArtifactId: "checker-qual", Version: "3.42.0"},
		},
	},
	"jmolecules": {
		Name:        "jMolecules",
		Description: "DDD architectural concepts",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jmolecules", ArtifactId: "jmolecules-ddd", Version: "1.9.0"},
		},
	},
	"vavr": {
		Name:        "Vavr",
		Description: "Functional programming library",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.vavr", ArtifactId: "vavr", Version: "0.10.4"},
		},
	},
	"immutables": {
		Name:        "Immutables",
		Description: "Immutable object generation",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.immutables", ArtifactId: "value", Version: "2.10.1", Scope: "provided"},
		},
	},
	"record-builder": {
		Name:        "RecordBuilder",
		Description: "Builder pattern for Java records",
		Category:    "Developer Tools",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.soabase.record-builder", ArtifactId: "record-builder-processor", Version: "40", Scope: "provided"},
		},
	},

	// ============================================
	// ADDITIONAL OBSERVABILITY
	// ============================================
	"jaeger": {
		Name:        "Jaeger",
		Description: "Distributed tracing with Jaeger",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.opentracing.contrib", ArtifactId: "opentracing-spring-jaeger-cloud-starter", Version: "3.3.1"},
		},
	},
	"datadog": {
		Name:        "Datadog",
		Description: "Datadog APM metrics",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-registry-datadog"},
		},
	},
	"newrelic": {
		Name:        "New Relic",
		Description: "New Relic APM metrics",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-registry-new-relic"},
		},
	},
	"grafana": {
		Name:        "Grafana LGTM",
		Description: "Grafana Loki/Tempo metrics",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.micrometer", ArtifactId: "micrometer-tracing-bridge-otel"},
		},
	},
	"sentry": {
		Name:        "Sentry",
		Description: "Error tracking and monitoring",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.sentry", ArtifactId: "sentry-spring-boot-starter-jakarta", Version: "7.4.0"},
		},
	},
	"loki": {
		Name:        "Loki Logback",
		Description: "Grafana Loki logging",
		Category:    "Observability",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.loki4j", ArtifactId: "loki-logback-appender", Version: "1.5.1"},
		},
	},

	// ============================================
	// WORKFLOW / BPM
	// ============================================
	"camunda": {
		Name:        "Camunda",
		Description: "Workflow and process automation",
		Category:    "Workflow",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.camunda.bpm.springboot", ArtifactId: "camunda-bpm-spring-boot-starter-rest", Version: "7.21.0"},
		},
	},
	"flowable": {
		Name:        "Flowable",
		Description: "BPMN workflow engine",
		Category:    "Workflow",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.flowable", ArtifactId: "flowable-spring-boot-starter", Version: "7.0.1"},
		},
	},
	"temporal": {
		Name:        "Temporal",
		Description: "Durable workflow orchestration",
		Category:    "Workflow",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.temporal", ArtifactId: "temporal-sdk", Version: "1.22.3"},
		},
	},

	// ============================================
	// API GATEWAY / SERVICE MESH
	// ============================================
	"spring-cloud-gateway": {
		Name:        "Spring Cloud Gateway",
		Description: "API Gateway for microservices",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-gateway"},
		},
	},
	"eureka": {
		Name:        "Eureka Client",
		Description: "Service discovery with Eureka",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-netflix-eureka-client"},
		},
	},
	"consul": {
		Name:        "Consul",
		Description: "Service discovery with Consul",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-consul-discovery"},
		},
	},
	"config-server": {
		Name:        "Config Server",
		Description: "Centralized configuration",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-config-server"},
		},
	},
	"config-client": {
		Name:        "Config Client",
		Description: "Config server client",
		Category:    "Web",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-config"},
		},
	},
	"vault": {
		Name:        "Vault",
		Description: "HashiCorp Vault integration",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-vault-config"},
		},
	},

	// ============================================
	// MAPS & LOCATION SERVICES
	// ============================================
	"google-maps": {
		Name:        "Google Maps",
		Description: "Google Maps Server SDK (Geocoding, Directions, Places)",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.maps", ArtifactId: "google-maps-services", Version: "2.2.0"},
		},
	},
	"mapbox": {
		Name:        "Mapbox",
		Description: "Mapbox Java SDK (Navigation, Search, Maps)",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.mapbox.mapboxsdk", ArtifactId: "mapbox-sdk-services", Version: "6.15.0"},
		},
	},
	"graphhopper": {
		Name:        "GraphHopper",
		Description: "Open-source routing engine (Directions API)",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.graphhopper", ArtifactId: "graphhopper-core", Version: "8.0"},
		},
	},
	"h3": {
		Name:        "H3",
		Description: "Uber's Hexagonal Hierarchical Spatial Index",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.uber", ArtifactId: "h3", Version: "4.1.1"},
		},
	},
	"jts": {
		Name:        "JTS Topology Suite",
		Description: "Geometry and spatial operations",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.locationtech.jts", ArtifactId: "jts-core", Version: "1.19.0"},
		},
	},
	"geotools": {
		Name:        "GeoTools",
		Description: "Geospatial data tools (Shapefiles, CRS)",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.geotools", ArtifactId: "gt-main", Version: "30.2"},
		},
	},
	"ip2location": {
		Name:        "IP2Location",
		Description: "IP Address to Geo-location lookup",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.ip2location", ArtifactId: "ip2location-java", Version: "8.11.1"},
		},
	},
	"maxmind": {
		Name:        "MaxMind GeoIP2",
		Description: "MaxMind GeoIP2 Database integration",
		Category:    "Maps & Geo",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.maxmind.geoip2", ArtifactId: "geoip2", Version: "4.2.0"},
		},
	},

	// ============================================
	// FILE PROCESSING & MEDIA
	// ============================================
	"ffmpeg": {
		Name:        "FFmpeg",
		Description: "FFmpeg wrapper for video/audio processing",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.bramp.ffmpeg", ArtifactId: "ffmpeg", Version: "0.8.0"},
		},
	},
	"tika": {
		Name:        "Apache Tika",
		Description: "Content detection and metadata extraction",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.tika", ArtifactId: "tika-core", Version: "2.9.1"},
		},
	},
	"batik": {
		Name:        "Apache Batik",
		Description: "SVG generation and manipulation",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.xmlgraphics", ArtifactId: "batik-transcoder", Version: "1.17"},
		},
	},
	"opencv": {
		Name:        "OpenCV",
		Description: "Computer Vision library bindings",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.openpnp", ArtifactId: "opencv", Version: "4.9.0-0"},
		},
	},
	"metadata-extractor": {
		Name:        "Metadata Extractor",
		Description: "Read Exif/IPTC metadata from images",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.drewnoakes", ArtifactId: "metadata-extractor", Version: "2.19.0"},
		},
	},
	"imgscalr": {
		Name:        "ImgScalr",
		Description: "Simple image scaling library",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.imgscalr", ArtifactId: "imgscalr-lib", Version: "4.2"},
		},
	},
	"pdfbox": {
		Name:        "Apache PDFBox",
		Description: "Create and edit PDF documents",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.pdfbox", ArtifactId: "pdfbox", Version: "3.0.1"},
		},
	},
	"openhtmltopdf": {
		Name:        "OpenHTMLtoPDF",
		Description: "Convert HTML/CSS directly to PDF",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.openhtmltopdf", ArtifactId: "openhtmltopdf-pdfbox", Version: "1.0.10"},
		},
	},
	"docx4j": {
		Name:        "Docx4j",
		Description: "Manipulate Word and PowerPoint files",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.docx4j", ArtifactId: "docx4j-JAXB-ReferenceImpl", Version: "11.4.10"},
		},
	},
	"jxls": {
		Name:        "JXLS",
		Description: "Excel templates with logic",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jxls", ArtifactId: "jxls-poi", Version: "3.0.0"},
		},
	},
	"flying-saucer": {
		Name:        "Flying Saucer",
		Description: "XHTML/CSS to PDF renderer",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.xhtmlrenderer", ArtifactId: "flying-saucer-pdf", Version: "9.3.1"},
		},
	},
	"barcode4j": {
		Name:        "Barcode4J",
		Description: "Barcode generation library",
		Category:    "Media",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.sf.barcode4j", ArtifactId: "barcode4j", Version: "2.1"},
		},
	},

	// ============================================
	// FINTECH & BLOCKCHAIN
	// ============================================
	"web3j": {
		Name:        "Web3j",
		Description: "Ethereum Blockchain integration",
		Category:    "Fintech",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.web3j", ArtifactId: "core", Version: "4.11.0"},
		},
	},
	"bitcoinj": {
		Name:        "BitcoinJ",
		Description: "Bitcoin protocol implementation",
		Category:    "Fintech",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.bitcoinj", ArtifactId: "bitcoinj-core", Version: "0.16.2"},
		},
	},
	"plaid": {
		Name:        "Plaid",
		Description: "Plaid SDK (Bank account linking)",
		Category:    "Fintech",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.plaid", ArtifactId: "plaid-java", Version: "18.2.0"},
		},
	},
	"xchange": {
		Name:        "XChange",
		Description: "Crypto exchange library (Binance, Coinbase)",
		Category:    "Fintech",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.knowm.xchange", ArtifactId: "xchange-core", Version: "5.1.1"},
		},
	},
	"stellar": {
		Name:        "Stellar",
		Description: "Stellar network SDK",
		Category:    "Fintech",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.stellar", ArtifactId: "java-stellar-sdk", Version: "0.43.0"},
		},
	},
	"braintree": {
		Name:        "Braintree",
		Description: "Braintree payment processing",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.braintreepayments.gateway", ArtifactId: "braintree-java", Version: "3.28.0"},
		},
	},
	"square": {
		Name:        "Square",
		Description: "Square payment processing",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.squareup", ArtifactId: "square", Version: "35.1.0.20240320"},
		},
	},
	"razorpay": {
		Name:        "Razorpay",
		Description: "Razorpay payment gateway",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.razorpay", ArtifactId: "razorpay-java", Version: "1.4.6"},
		},
	},
	"mollie": {
		Name:        "Mollie",
		Description: "Mollie payment gateway",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.mollie", ArtifactId: "mollie-api", Version: "3.4.0"},
		},
	},
	"adyen": {
		Name:        "Adyen",
		Description: "Adyen payment platform",
		Category:    "Payments",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.adyen", ArtifactId: "adyen-java-api-library", Version: "25.1.0"},
		},
	},

	// ============================================
	// COMMUNICATION & SOCIAL
	// ============================================
	"facebook-sdk": {
		Name:        "Facebook SDK",
		Description: "Facebook Graph API integration",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.facebook.business.sdk", ArtifactId: "facebook-java-business-sdk", Version: "19.0.3"},
		},
	},
	"twitter-api": {
		Name:        "Twitter API",
		Description: "Twitter/X API v2 integration",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.github.redouane59.twitter", ArtifactId: "twittered", Version: "2.23"},
		},
	},
	"linkedin-api": {
		Name:        "LinkedIn API",
		Description: "LinkedIn API integration",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.linkedin.dex", ArtifactId: "api", Version: "3.0.0"},
		},
	},
	"zoom": {
		Name:        "Zoom",
		Description: "Zoom API (Meeting generation)",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "us.zoom", ArtifactId: "zoomsdk", Version: "5.15.7"},
		},
	},
	"agora": {
		Name:        "Agora",
		Description: "Agora SDK (Real-time Video/Audio)",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.agora", ArtifactId: "authentication", Version: "2.0.0"},
		},
	},
	"matrix": {
		Name:        "Matrix",
		Description: "Matrix.org protocol (Decentralized chat)",
		Category:    "Social",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.github.ma1uta.matrix", ArtifactId: "client", Version: "0.14.0"},
		},
	},

	// ============================================
	// AI & DATA ENGINEERING (EXPANDED)
	// ============================================
	"spark": {
		Name:        "Apache Spark",
		Description: "Apache Spark Java API",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.spark", ArtifactId: "spark-core_2.13", Version: "3.5.0"},
		},
	},
	"flink": {
		Name:        "Apache Flink",
		Description: "Stream processing framework",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.flink", ArtifactId: "flink-java", Version: "1.18.1"},
		},
	},
	"hadoop": {
		Name:        "Apache Hadoop",
		Description: "Hadoop client libraries",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.hadoop", ArtifactId: "hadoop-client", Version: "3.3.6"},
		},
	},
	"djl": {
		Name:        "Deep Java Library",
		Description: "Amazon's ML toolkit for Java",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "ai.djl", ArtifactId: "api", Version: "0.27.0"},
		},
	},
	"weka": {
		Name:        "Weka",
		Description: "Machine Learning algorithms collection",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "nz.ac.waikato.cms.weka", ArtifactId: "weka-stable", Version: "3.8.6"},
		},
	},
	"corenlp": {
		Name:        "Stanford CoreNLP",
		Description: "Natural Language Processing",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "edu.stanford.nlp", ArtifactId: "stanford-corenlp", Version: "4.5.6"},
		},
	},
	"langchain4j": {
		Name:        "LangChain4j",
		Description: "LangChain for Java (LLM orchestration)",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "dev.langchain4j", ArtifactId: "langchain4j", Version: "0.28.0"},
		},
	},
	"pinecone": {
		Name:        "Pinecone",
		Description: "Pinecone vector database client",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.pinecone", ArtifactId: "pinecone-client", Version: "1.0.0"},
		},
	},
	"weaviate": {
		Name:        "Weaviate",
		Description: "Weaviate vector database client",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.weaviate", ArtifactId: "client", Version: "4.5.1"},
		},
	},
	"qdrant": {
		Name:        "Qdrant",
		Description: "Qdrant vector database client",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.qdrant", ArtifactId: "client", Version: "1.7.2"},
		},
	},
	"chroma": {
		Name:        "Chroma",
		Description: "Chroma vector database client",
		Category:    "AI & Machine Learning",
		Dependencies: []buildtool.Dependency{
			{GroupId: "dev.langchain4j", ArtifactId: "langchain4j-chroma", Version: "0.28.0"},
		},
	},

	// ============================================
	// FEATURE FLAGS & CONFIG
	// ============================================
	"unleash": {
		Name:        "Unleash",
		Description: "Unleash feature toggle client",
		Category:    "Feature Flags",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.getunleash", ArtifactId: "unleash-client-java", Version: "9.2.0"},
		},
	},
	"launchdarkly": {
		Name:        "LaunchDarkly",
		Description: "LaunchDarkly feature flags SDK",
		Category:    "Feature Flags",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.launchdarkly", ArtifactId: "launchdarkly-java-server-sdk", Version: "7.2.5"},
		},
	},
	"flagsmith": {
		Name:        "Flagsmith",
		Description: "Flagsmith feature flags",
		Category:    "Feature Flags",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.flagsmith", ArtifactId: "flagsmith-java-client", Version: "7.1.0"},
		},
	},
	"togglz": {
		Name:        "Togglz",
		Description: "Feature Flags for Java (Spring)",
		Category:    "Feature Flags",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.togglz", ArtifactId: "togglz-spring-boot-starter", Version: "4.4.0"},
		},
	},
	"ff4j": {
		Name:        "FF4J",
		Description: "Feature Flipping for Java",
		Category:    "Feature Flags",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.ff4j", ArtifactId: "ff4j-spring-boot-starter", Version: "2.1"},
		},
	},

	// ============================================
	// IDENTITY & AUTH (EXPANDED)
	// ============================================
	"supertokens": {
		Name:        "SuperTokens",
		Description: "Open source auth provider",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.supertokens", ArtifactId: "supertokens-plugin-interface", Version: "6.0.0"},
		},
	},
	"kratos": {
		Name:        "Ory Kratos",
		Description: "Identity management client",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "sh.ory.kratos", ArtifactId: "kratos-client", Version: "1.0.0"},
		},
	},
	"cas": {
		Name:        "CAS Client",
		Description: "Central Authentication Service client",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apereo.cas.client", ArtifactId: "cas-client-support-springboot", Version: "4.0.3"},
		},
	},
	"saml": {
		Name:        "Spring Security SAML",
		Description: "SAML authentication extension",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.security", ArtifactId: "spring-security-saml2-service-provider"},
		},
	},
	"spring-session": {
		Name:        "Spring Session",
		Description: "Distributed session management",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.session", ArtifactId: "spring-session-core"},
		},
	},
	"recaptcha": {
		Name:        "reCAPTCHA",
		Description: "Google reCAPTCHA integration",
		Category:    "Security",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.recaptcha", ArtifactId: "recaptcha", Version: "1.0.0"},
		},
	},

	// ============================================
	// MICROSERVICES (ADVANCED PATTERNS)
	// ============================================
	"spring-cloud-stream": {
		Name:        "Spring Cloud Stream",
		Description: "Event-driven framework (Kafka/RabbitMQ abstraction)",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-stream-kafka"},
		},
	},
	"spring-cloud-bus": {
		Name:        "Spring Cloud Bus",
		Description: "Broadcasts state changes across nodes",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-bus-amqp"},
		},
	},
	"spring-cloud-function": {
		Name:        "Spring Cloud Function",
		Description: "Write once, run as Web/Lambda/Serverless",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-function-web"},
		},
	},
	"spring-retry": {
		Name:        "Spring Retry",
		Description: "Simple retry logic for operations",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.retry", ArtifactId: "spring-retry"},
		},
	},
	"shedlock": {
		Name:        "ShedLock",
		Description: "Distributed lock for @Scheduled tasks",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.javacrumbs.shedlock", ArtifactId: "shedlock-spring", Version: "5.10.2"},
			{GroupId: "net.javacrumbs.shedlock", ArtifactId: "shedlock-provider-jdbc-template", Version: "5.10.2"},
		},
	},
	"pact": {
		Name:        "Pact",
		Description: "Consumer-Driven Contract testing",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "au.com.dius.pact.consumer", ArtifactId: "junit5", Version: "4.6.7", Scope: "test"},
		},
	},
	"etcd": {
		Name:        "Etcd",
		Description: "Etcd client (distributed key-value store)",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.etcd", ArtifactId: "jetcd-core", Version: "0.8.0"},
		},
	},
	"zookeeper": {
		Name:        "Apache Zookeeper",
		Description: "Zookeeper client",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.curator", ArtifactId: "curator-framework", Version: "5.6.0"},
		},
	},
	"dapr": {
		Name:        "Dapr",
		Description: "Dapr SDK (Distributed Application Runtime)",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.dapr", ArtifactId: "dapr-sdk-springboot", Version: "1.11.0"},
		},
	},
	"grpc-server": {
		Name:        "gRPC Server",
		Description: "gRPC server starter",
		Category:    "Microservices",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.devh", ArtifactId: "grpc-server-spring-boot-starter", Version: "3.0.0.RELEASE"},
		},
	},

	// ============================================
	// ENTERPRISE INTEGRATION (EIP)
	// ============================================
	"spring-integration": {
		Name:        "Spring Integration",
		Description: "Enterprise Integration Patterns (EIP)",
		Category:    "Integration",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-integration"},
		},
	},
	"apache-camel": {
		Name:        "Apache Camel",
		Description: "Integration framework (Routing and mediation)",
		Category:    "Integration",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.camel.springboot", ArtifactId: "camel-spring-boot-starter", Version: "4.4.0"},
		},
	},
	"spring-cloud-data-flow": {
		Name:        "Spring Cloud Data Flow",
		Description: "Orchestration for data microservices",
		Category:    "Integration",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-dataflow-rest-client", Version: "2.11.2"},
		},
	},
	"spring-cloud-task": {
		Name:        "Spring Cloud Task",
		Description: "Short-lived microservices (ephemeral tasks)",
		Category:    "Integration",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.cloud", ArtifactId: "spring-cloud-starter-task"},
		},
	},

	// ============================================
	// IoT & HARDWARE PROTOCOLS
	// ============================================
	"mqtt-paho": {
		Name:        "Eclipse Paho MQTT",
		Description: "MQTT messaging client",
		Category:    "IoT",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.eclipse.paho", ArtifactId: "org.eclipse.paho.client.mqttv3", Version: "1.2.5"},
		},
	},
	"spring-integration-mqtt": {
		Name:        "Spring Integration MQTT",
		Description: "Spring's MQTT integration",
		Category:    "IoT",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.integration", ArtifactId: "spring-integration-mqtt"},
		},
	},
	"coap": {
		Name:        "Californium CoAP",
		Description: "CoAP for constrained devices",
		Category:    "IoT",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.eclipse.californium", ArtifactId: "californium-core", Version: "3.10.0"},
		},
	},
	"modbus": {
		Name:        "Modbus",
		Description: "Modbus protocol (Industrial automation)",
		Category:    "IoT",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.ghgande", ArtifactId: "j2mod", Version: "3.2.1"},
		},
	},
	"pi4j": {
		Name:        "Pi4J",
		Description: "Raspberry Pi I/O control for Java",
		Category:    "IoT",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.pi4j", ArtifactId: "pi4j-core", Version: "2.5.0"},
		},
	},

	// ============================================
	// DATA SCIENCE & MATH
	// ============================================
	"tablesaw": {
		Name:        "Tablesaw",
		Description: "Java Dataframes (Pandas equivalent)",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "tech.tablesaw", ArtifactId: "tablesaw-core", Version: "0.43.1"},
		},
	},
	"commons-math": {
		Name:        "Apache Commons Math",
		Description: "Math and statistics library",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.commons", ArtifactId: "commons-math3", Version: "3.6.1"},
		},
	},
	"nd4j": {
		Name:        "ND4J",
		Description: "N-Dimensional Arrays (NumPy equivalent)",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.nd4j", ArtifactId: "nd4j-native-platform", Version: "1.0.0-M2.1"},
		},
	},
	"joda-money": {
		Name:        "Joda-Money",
		Description: "Currency and money handling",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.joda", ArtifactId: "joda-money", Version: "1.0.4"},
		},
	},
	"jgrapht": {
		Name:        "JGraphT",
		Description: "Graph theory data structures and algorithms",
		Category:    "Data Processing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jgrapht", ArtifactId: "jgrapht-core", Version: "1.5.2"},
		},
	},

	// ============================================
	// CONTAINERIZATION & BUILD TOOLS
	// ============================================
	"jib": {
		Name:        "Google Jib",
		Description: "Build Docker images without Docker daemon",
		Category:    "DevOps",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.google.cloud.tools", ArtifactId: "jib-core", Version: "0.27.0"},
		},
	},
	"docker-java": {
		Name:        "Docker Java",
		Description: "Java API for Docker Engine",
		Category:    "DevOps",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.docker-java", ArtifactId: "docker-java-core", Version: "3.3.6"},
		},
	},
	"kubernetes-client": {
		Name:        "Kubernetes Client",
		Description: "Official Kubernetes Java client",
		Category:    "DevOps",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.kubernetes", ArtifactId: "client-java", Version: "20.0.0"},
		},
	},
	"fabric8": {
		Name:        "Fabric8 Kubernetes",
		Description: "Popular Kubernetes client",
		Category:    "DevOps",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.fabric8", ArtifactId: "kubernetes-client", Version: "6.10.0"},
		},
	},

	// ============================================
	// ADVANCED TESTING (E2E & UI)
	// ============================================
	"selenium": {
		Name:        "Selenium",
		Description: "Browser automation framework",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.seleniumhq.selenium", ArtifactId: "selenium-java", Version: "4.18.1", Scope: "test"},
		},
	},
	"selenide": {
		Name:        "Selenide",
		Description: "Concise UI tests (Selenium wrapper)",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.codeborne", ArtifactId: "selenide", Version: "7.2.1", Scope: "test"},
		},
	},
	"playwright": {
		Name:        "Playwright",
		Description: "Microsoft Playwright browser automation",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.microsoft.playwright", ArtifactId: "playwright", Version: "1.41.2", Scope: "test"},
		},
	},
	"cucumber": {
		Name:        "Cucumber",
		Description: "BDD (Behavior Driven Development) testing",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.cucumber", ArtifactId: "cucumber-java", Version: "7.15.0", Scope: "test"},
			{GroupId: "io.cucumber", ArtifactId: "cucumber-spring", Version: "7.15.0", Scope: "test"},
		},
	},
	"hoverfly": {
		Name:        "Hoverfly",
		Description: "API simulation/virtualization",
		Category:    "Testing",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.specto", ArtifactId: "hoverfly-java", Version: "0.18.1", Scope: "test"},
		},
	},

	// ============================================
	// CODE QUALITY & FORMATTING
	// ============================================
	"jacoco": {
		Name:        "JaCoCo",
		Description: "Java Code Coverage Library",
		Category:    "Quality",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jacoco", ArtifactId: "jacoco-maven-plugin", Version: "0.8.11"},
		},
	},
	"checkstyle": {
		Name:        "Checkstyle",
		Description: "Coding standards enforcement",
		Category:    "Quality",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.puppycrawl.tools", ArtifactId: "checkstyle", Version: "10.14.0"},
		},
	},
	"spotless": {
		Name:        "Spotless",
		Description: "Auto-format code during build",
		Category:    "Quality",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.diffplug.spotless", ArtifactId: "spotless-maven-plugin", Version: "2.43.0"},
		},
	},
	"pmd": {
		Name:        "PMD",
		Description: "Source code analyzer",
		Category:    "Quality",
		Dependencies: []buildtool.Dependency{
			{GroupId: "net.sourceforge.pmd", ArtifactId: "pmd-java", Version: "7.0.0"},
		},
	},
	"sonar": {
		Name:        "SonarQube Scanner",
		Description: "SonarQube integration",
		Category:    "Quality",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.sonarsource.scanner.maven", ArtifactId: "sonar-maven-plugin", Version: "3.10.0.2594"},
		},
	},

	// ============================================
	// ADVANCED CACHING & PERFORMANCE
	// ============================================
	"caffeine": {
		Name:        "Caffeine",
		Description: "High-performance local caching",
		Category:    "Caching",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.ben-manes.caffeine", ArtifactId: "caffeine", Version: "3.1.8"},
		},
	},
	"ehcache": {
		Name:        "Ehcache",
		Description: "Robust enterprise caching",
		Category:    "Caching",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.ehcache", ArtifactId: "ehcache", Version: "3.10.8"},
		},
	},
	"hazelcast-jet": {
		Name:        "Hazelcast Jet",
		Description: "Distributed stream processing",
		Category:    "Caching",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.hazelcast.jet", ArtifactId: "hazelcast-jet", Version: "5.3"},
		},
	},
	"infinispan": {
		Name:        "Infinispan",
		Description: "Distributed in-memory data grid",
		Category:    "Caching",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.infinispan", ArtifactId: "infinispan-spring-boot3-starter-embedded", Version: "15.0.0.Final"},
		},
	},

	// ============================================
	// CONTENT & FEEDS
	// ============================================
	"rome": {
		Name:        "Rome",
		Description: "RSS and Atom feed parser/generator",
		Category:    "Content",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.rometools", ArtifactId: "rome", Version: "2.1.0"},
		},
	},
	"htmlunit": {
		Name:        "HtmlUnit",
		Description: "Headless browser for scraping/testing",
		Category:    "Content",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.htmlunit", ArtifactId: "htmlunit", Version: "4.0.0"},
		},
	},
	"bliki": {
		Name:        "Bliki",
		Description: "Wikipedia syntax parser",
		Category:    "Content",
		Dependencies: []buildtool.Dependency{
			{GroupId: "info.bliki.wiki", ArtifactId: "bliki-core", Version: "3.1.0"},
		},
	},
	"emoji-java": {
		Name:        "Emoji Java",
		Description: "Emoji handling in strings",
		Category:    "Content",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.vdurmont", ArtifactId: "emoji-java", Version: "5.1.1"},
		},
	},

	// ============================================
	// NETWORKING & IP TOOLS
	// ============================================
	"netty": {
		Name:        "Netty",
		Description: "Async event-driven network framework",
		Category:    "Networking",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.netty", ArtifactId: "netty-all", Version: "4.1.107.Final"},
		},
	},
	"pcap4j": {
		Name:        "Pcap4J",
		Description: "Packet capture library",
		Category:    "Networking",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.pcap4j", ArtifactId: "pcap4j-core", Version: "1.8.2"},
		},
	},
	"dns-java": {
		Name:        "dnsjava",
		Description: "DNS protocol implementation",
		Category:    "Networking",
		Dependencies: []buildtool.Dependency{
			{GroupId: "dnsjava", ArtifactId: "dnsjava", Version: "3.5.3"},
		},
	},

	// ============================================
	// API & SCHEMA DESIGN
	// ============================================
	"springdoc": {
		Name:        "SpringDoc",
		Description: "OpenAPI 3 (Swagger) generation",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springdoc", ArtifactId: "springdoc-openapi-starter-webmvc-ui", Version: "2.3.0"},
		},
	},
	"asyncapi": {
		Name:        "AsyncAPI",
		Description: "Event-Driven API documentation",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.asyncapi", ArtifactId: "asyncapi-core", Version: "1.0.0-EAP"},
		},
	},
	"netflix-dgs": {
		Name:        "Netflix DGS",
		Description: "Netflix Domain Graph Service (GraphQL)",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.netflix.graphql.dgs", ArtifactId: "graphql-dgs-spring-boot-starter", Version: "8.4.0"},
		},
	},
	"avro-serializer": {
		Name:        "Avro Serializer",
		Description: "Avro serialization for Kafka schemas",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "io.confluent", ArtifactId: "kafka-avro-serializer", Version: "7.6.0"},
		},
	},
	"json-schema": {
		Name:        "JSON Schema",
		Description: "JSON Schema validation",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.networknt", ArtifactId: "json-schema-validator", Version: "1.3.3"},
		},
	},

	// ============================================
	// JOB SCHEDULING (EXPANDED)
	// ============================================
	"jobrunr": {
		Name:        "JobRunr",
		Description: "Distributed background job processing",
		Category:    "Scheduling",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.jobrunr", ArtifactId: "jobrunr-spring-boot-3-starter", Version: "7.1.0"},
		},
	},
	"db-scheduler": {
		Name:        "DB Scheduler",
		Description: "Persistent task scheduler",
		Category:    "Scheduling",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.github.kagkarlsson", ArtifactId: "db-scheduler-spring-boot-starter", Version: "14.0.1"},
		},
	},

	// ============================================
	// DATABASES (ADDITIONAL)
	// ============================================
	"sqlite": {
		Name:        "SQLite",
		Description: "SQLite embedded database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.xerial", ArtifactId: "sqlite-jdbc", Version: "3.45.1.0", Scope: "runtime"},
		},
	},
	"hsqldb": {
		Name:        "HSQLDB",
		Description: "HyperSQL embedded database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.hsqldb", ArtifactId: "hsqldb", Scope: "runtime"},
		},
	},
	"derby": {
		Name:        "Apache Derby",
		Description: "Apache Derby embedded database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.apache.derby", ArtifactId: "derby", Version: "10.17.1.0", Scope: "runtime"},
		},
	},
	"duckdb": {
		Name:        "DuckDB",
		Description: "DuckDB analytical database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.duckdb", ArtifactId: "duckdb_jdbc", Version: "0.10.0", Scope: "runtime"},
		},
	},
	"questdb": {
		Name:        "QuestDB",
		Description: "QuestDB time-series database",
		Category:    "SQL",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.questdb", ArtifactId: "questdb", Version: "7.4.0"},
		},
	},

	// ============================================
	// LOGGING (EXPLICIT)
	// ============================================
	"logback": {
		Name:        "Logback",
		Description: "Logback logging framework",
		Category:    "Logging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "ch.qos.logback", ArtifactId: "logback-classic"},
		},
	},
	"log4j2": {
		Name:        "Log4j2",
		Description: "Apache Log4j 2 logging",
		Category:    "Logging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.springframework.boot", ArtifactId: "spring-boot-starter-log4j2"},
		},
	},
	"slf4j": {
		Name:        "SLF4J",
		Description: "Simple Logging Facade for Java",
		Category:    "Logging",
		Dependencies: []buildtool.Dependency{
			{GroupId: "org.slf4j", ArtifactId: "slf4j-api"},
		},
	},

	// ============================================
	// GRAPHQL (EXPANDED)
	// ============================================
	"graphql-kickstart": {
		Name:        "GraphQL Kickstart",
		Description: "GraphQL Java Kickstart starter",
		Category:    "API",
		Dependencies: []buildtool.Dependency{
			{GroupId: "com.graphql-java-kickstart", ArtifactId: "graphql-spring-boot-starter", Version: "15.1.0"},
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
