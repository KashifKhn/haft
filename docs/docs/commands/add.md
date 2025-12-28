---
sidebar_position: 3
title: haft add
description: Add dependencies to your project
---

# haft add

Add dependencies to an existing Spring Boot project.

## Usage

```bash
haft add                              # Interactive search picker
haft add --browse                     # Browse by category
haft add <dependency> [dependencies...]
haft add <groupId:artifactId>
haft add <groupId:artifactId:version>
```

## Description

The `add` command modifies your build file (`pom.xml` or `build.gradle`) to add new dependencies. It supports:

- **Interactive mode** — Search and select from 330+ shortcuts
- **Browse mode** — Navigate dependencies by category
- **Shortcuts** — Common dependencies like `lombok`, `jpa`, `web`, `jwt`
- **Maven coordinates** — Any dependency as `groupId:artifactId`
- **Maven Central verification** — Auto-verify and fetch latest versions

## Interactive Modes

### Search Picker (Default)

```bash
haft add
```

Opens an interactive TUI where you can:
- Type to search/filter dependencies
- Select multiple with `Space`
- Navigate with `↑`/`↓`, `PgUp`/`PgDown`
- Select all visible with `a`, none with `n`
- Confirm with `Enter`, cancel with `Esc`

### Category Browser

```bash
haft add --browse
haft add -b
```

Opens a category-based browser similar to the init wizard:
- Jump to categories with `0-9` keys
- Cycle categories with `Tab`/`Shift+Tab`
- Search within category with `/`

## Examples

### Add Using Shortcuts

```bash
# Add Lombok
haft add lombok

# Add multiple dependencies
haft add jpa validation lombok

# Add JWT (adds all 3 JJWT artifacts)
haft add jwt

# Add database driver
haft add postgresql
```

### Add Using Maven Coordinates

```bash
# Without version (fetches latest from Maven Central)
haft add org.mapstruct:mapstruct

# With specific version
haft add io.jsonwebtoken:jjwt-api:0.12.5
```

Dependencies are automatically verified against Maven Central. If a dependency doesn't exist, you'll get an error:

```
ERROR ✗ dependency 'com.fake:nonexistent' not found on Maven Central
```

### Override Scope

```bash
# Add as test dependency
haft add h2 --scope test

# Add as provided
haft add org.example:my-processor --scope provided
```

### List Available Shortcuts

```bash
haft add --list
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--browse` | `-b` | Browse dependencies by category |
| `--list` | | List available dependency shortcuts |
| `--scope` | | Set dependency scope (compile, runtime, test, provided) |
| `--version` | | Override default version |
| `--json` | | Output result as JSON |
| `--no-interactive` | | Skip interactive prompts |

## Available Shortcuts (330+)

### Web

| Shortcut | Description |
|----------|-------------|
| `web` | Spring Boot Web (Spring MVC) |
| `webflux` | Spring WebFlux (reactive) |
| `graphql` | Spring GraphQL |
| `websocket` | WebSocket support |
| `hateoas` | Hypermedia-driven REST APIs |
| `data-rest` | Expose repositories as REST endpoints |
| `feign` | Declarative REST client (OpenFeign) |
| `resilience4j` | Fault tolerance (circuit breaker) |
| `restclient` | Synchronous HTTP client |
| `webclient` | Reactive HTTP client |
| `session-redis` | Distributed sessions with Redis |
| `session-jdbc` | Distributed sessions with JDBC |
| `grpc` | gRPC framework for Spring Boot |
| `bucket4j` | Rate limiting library |
| `spring-cloud-gateway` | API Gateway for microservices |
| `eureka` | Service discovery with Eureka |
| `consul` | Service discovery with Consul |
| `config-server` | Centralized configuration server |
| `config-client` | Config server client |

### SQL

| Shortcut | Description |
|----------|-------------|
| `jpa` | Spring Data JPA |
| `jdbc` | Spring JDBC |
| `postgresql` | PostgreSQL driver |
| `mysql` | MySQL driver |
| `mariadb` | MariaDB driver |
| `h2` | H2 in-memory database |
| `flyway` | Flyway migrations |
| `liquibase` | Liquibase migrations |
| `r2dbc` | Reactive database access |
| `r2dbc-postgresql` | Reactive PostgreSQL driver |
| `mybatis` | MyBatis SQL mapping framework |
| `jooq` | Type-safe SQL query builder |
| `oracle` | Oracle database JDBC driver |
| `sqlserver` | Microsoft SQL Server JDBC driver |
| `clickhouse` | ClickHouse column-oriented OLAP |
| `cockroachdb` | CockroachDB distributed SQL |
| `timescaledb` | TimescaleDB time-series |
| `voltdb` | VoltDB in-memory SQL |
| `sqlite` | SQLite embedded database |
| `hsqldb` | HSQLDB embedded database |
| `derby` | Apache Derby embedded database |
| `duckdb` | DuckDB analytical database |
| `questdb` | QuestDB time-series database |

### NoSQL

| Shortcut | Description |
|----------|-------------|
| `mongodb` | Spring Data MongoDB |
| `mongodb-reactive` | Reactive MongoDB support |
| `redis` | Spring Data Redis |
| `elasticsearch` | Spring Data Elasticsearch |
| `cassandra` | Spring Data Cassandra |
| `neo4j` | Spring Data Neo4j |
| `couchbase` | Couchbase NoSQL database |
| `scylladb` | ScyllaDB (Cassandra-compatible) |
| `influxdb` | InfluxDB time-series database |
| `arangodb` | ArangoDB multi-model database |
| `hazelcast` | Hazelcast in-memory data grid |
| `ignite` | Apache Ignite distributed computing |
| `memcached` | Memcached distributed caching |
| `dynamodb-local` | DynamoDB enhanced client |

### Security

| Shortcut | Description |
|----------|-------------|
| `security` | Spring Security |
| `oauth2-client` | OAuth2 client |
| `oauth2-resource-server` | OAuth2 resource server |
| `oauth2-authorization-server` | Build your own OAuth2 authorization server |
| `jwt` | JJWT library (api + impl + jackson) |
| `ldap` | LDAP authentication and operations |
| `passay` | Password validation library |
| `keycloak` | Keycloak admin client |
| `auth0` | Auth0 authentication |
| `vault` | HashiCorp Vault integration |
| `bouncy-castle` | Bouncy Castle cryptography |
| `jasypt` | Encryption for properties |
| `nimbus-jose` | Nimbus JOSE+JWT |
| `otp` | TOTP/HOTP for 2FA |
| `supertokens` | SuperTokens open source auth |
| `kratos` | Ory Kratos identity management |
| `cas` | CAS authentication client |
| `saml` | Spring Security SAML |
| `spring-session` | Distributed session management |
| `recaptcha` | Google reCAPTCHA integration |

### Messaging

| Shortcut | Description |
|----------|-------------|
| `amqp` | RabbitMQ (Spring AMQP) |
| `rabbitmq` | RabbitMQ Java client (direct) |
| `kafka` | Apache Kafka |
| `pulsar` | Apache Pulsar messaging |
| `activemq` | ActiveMQ JMS messaging |
| `artemis` | Apache ActiveMQ Artemis |
| `rsocket` | RSocket reactive messaging |
| `nats` | NATS messaging system |
| `zeromq` | ZeroMQ high-performance messaging |

### I/O

| Shortcut | Description |
|----------|-------------|
| `validation` | Bean Validation |
| `mail` | Java Mail |
| `cache` | Spring Cache |
| `batch` | Spring Batch |
| `quartz` | Quartz Scheduler |
| `scheduler` | Task scheduling with @Scheduled |
| `commons-io` | Apache Commons IO |
| `jackson-datatype` | Jackson Java 8 datatypes |
| `apache-poi` | Microsoft Office file handling |
| `itext` | PDF generation and manipulation |
| `minio` | S3-compatible object storage client |
| `json-path` | JSON path query library |

### Template Engines

| Shortcut | Description |
|----------|-------------|
| `thymeleaf` | Thymeleaf templates |
| `freemarker` | FreeMarker templates |
| `mustache` | Mustache templates |

### Ops

| Shortcut | Description |
|----------|-------------|
| `actuator` | Spring Boot Actuator |
| `micrometer` | Prometheus metrics exporter |

### Observability

| Shortcut | Description |
|----------|-------------|
| `prometheus` | Prometheus metrics registry |
| `zipkin` | Distributed tracing with Zipkin |
| `opentelemetry` | OpenTelemetry distributed tracing |
| `jaeger` | Distributed tracing with Jaeger |
| `datadog` | Datadog APM metrics |
| `newrelic` | New Relic APM metrics |
| `sentry` | Sentry error tracking |
| `loki` | Grafana Loki logging |
| `grafana` | Grafana LGTM metrics |

### AI

| Shortcut | Description |
|----------|-------------|
| `openai` | OpenAI (ChatGPT, GPT-4) integration |
| `anthropic` | Anthropic Claude integration |
| `ollama` | Ollama local LLM integration |
| `azure-openai` | Azure OpenAI integration |
| `bedrock` | Amazon Bedrock AI integration |
| `vertex-ai` | Google Vertex AI Gemini integration |
| `mistral` | Mistral AI integration |
| `pgvector` | PostgreSQL vector database for AI |
| `langchain4j` | LangChain for Java (LLM orchestration) |
| `pinecone` | Pinecone vector database client |
| `weaviate` | Weaviate vector database client |
| `qdrant` | Qdrant vector database client |
| `chroma` | Chroma vector database client |
| `djl` | Deep Java Library (Amazon ML toolkit) |
| `weka` | Weka Machine Learning algorithms |
| `corenlp` | Stanford CoreNLP (NLP processing) |

### Cloud

| Shortcut | Description |
|----------|-------------|
| `aws-s3` | Amazon S3 file storage |
| `aws-sqs` | Amazon Simple Queue Service |
| `aws-ses` | Amazon Simple Email Service |
| `aws-sns` | Amazon Simple Notification Service |
| `aws-dynamodb` | Amazon DynamoDB NoSQL database |
| `aws-lambda` | Amazon Lambda serverless |
| `aws-cognito` | Amazon Cognito authentication |
| `aws-secretsmanager` | Amazon Secrets Manager |
| `aws-cloudwatch` | Amazon CloudWatch monitoring |
| `aws-kinesis` | Amazon Kinesis streaming |
| `gcp-storage` | Google Cloud Storage |
| `gcp-pubsub` | Google Cloud Pub/Sub |
| `gcp-bigquery` | Google BigQuery |
| `gcp-firestore` | Google Firestore |
| `gcp-secretmanager` | Google Secret Manager |
| `gcp-functions` | Google Cloud Functions |
| `azure-storage` | Azure Blob Storage |
| `azure-servicebus` | Azure Service Bus |
| `azure-keyvault` | Azure Key Vault |
| `azure-cosmosdb` | Azure Cosmos DB |
| `azure-functions` | Azure Functions |

### Notifications

| Shortcut | Description |
|----------|-------------|
| `twilio` | SMS, WhatsApp, Voice API |
| `sendgrid` | SendGrid email delivery |
| `mailersend` | MailerSend transactional email |
| `mailgun` | Mailgun email delivery |
| `firebase-admin` | Firebase push notifications |
| `pusher` | Pusher real-time messaging |
| `slack` | Slack API integration |
| `discord` | Discord bot integration |
| `telegram` | Telegram bot API |
| `onesignal` | OneSignal push notifications |

### Payments

| Shortcut | Description |
|----------|-------------|
| `stripe` | Stripe payment processing |
| `paypal` | PayPal payment SDK |
| `braintree` | Braintree payment processing |
| `square` | Square payment processing |
| `razorpay` | Razorpay payment gateway |
| `mollie` | Mollie payment gateway |
| `adyen` | Adyen payment platform |

### Search

| Shortcut | Description |
|----------|-------------|
| `algolia` | Algolia search-as-a-service |
| `meilisearch` | Meilisearch open-source search |
| `typesense` | Typesense typo-tolerant search |
| `solr` | Apache Solr enterprise search |
| `opensearch` | OpenSearch (Elasticsearch fork) |

### Utilities

| Shortcut | Description |
|----------|-------------|
| `okhttp` | OkHttp HTTP client |
| `retrofit` | Retrofit type-safe REST client |
| `httpclient5` | Apache HttpClient 5 |
| `jsoup` | HTML parsing and web scraping |
| `zxing` | QR code and barcode generation |
| `thumbnailator` | Image resizing and processing |
| `jasperreports` | Report generation |
| `jfreechart` | Chart generation |
| `opencsv` | CSV parsing and writing |
| `snakeyaml` | YAML parsing |
| `commonmark` | Markdown parsing |
| `flexmark` | Advanced Markdown processor |
| `protobuf` | Google Protocol Buffers |
| `avro` | Apache Avro serialization |
| `msgpack` | MessagePack binary format |
| `kryo` | Fast binary serialization |
| `jaxb` | XML binding (Java 11+) |
| `dom4j` | XML parsing |
| `xstream` | XML serialization |
| `jsch` | SSH and SFTP client |
| `sshj` | Modern SSH client |
| `libphonenumber` | Phone number validation |
| `slug` | URL slug generation |
| `emoji` | Emoji parsing |
| `joda-time` | Date/time library (legacy) |
| `jna` | Java Native Access |

### Workflow

| Shortcut | Description |
|----------|-------------|
| `camunda` | Camunda workflow automation |
| `flowable` | Flowable BPMN engine |
| `temporal` | Temporal workflow orchestration |

### Developer Tools

| Shortcut | Description |
|----------|-------------|
| `lombok` | Lombok annotations |
| `devtools` | Spring Boot DevTools |
| `mapstruct` | MapStruct bean mapping |
| `openapi` | SpringDoc OpenAPI (Swagger UI) |
| `commons-lang` | Apache Commons Lang |
| `guava` | Google Guava |
| `modelmapper` | ModelMapper |
| `config-processor` | IDE support for @ConfigurationProperties |
| `docker-compose` | Docker Compose development support |
| `native` | GraalVM native image support |
| `vavr` | Functional programming library |
| `immutables` | Immutable object generation |
| `record-builder` | Builder pattern for Java records |
| `jmolecules` | DDD architectural concepts |
| `spotbugs` | SpotBugs static analysis |
| `error-prone` | Error Prone static analysis |
| `checker-qual` | Checker Framework annotations |

### Testing

| Shortcut | Description |
|----------|-------------|
| `test` | Spring Boot Test |
| `testcontainers` | Testcontainers |
| `security-test` | Spring Security Test |
| `mockito` | Mockito mocking |
| `restdocs` | Spring REST Docs |
| `wiremock` | WireMock HTTP mock server |
| `rest-assured` | REST Assured API testing |
| `assertj` | AssertJ fluent assertions |
| `awaitility` | Awaitility async testing |
| `archunit` | ArchUnit architecture testing |
| `jsonassert` | JSONAssert comparison |
| `gatling` | Gatling load testing |
| `jmh` | JMH microbenchmark harness |
| `faker` | DataFaker test data generation |
| `greenmail` | GreenMail email testing |
| `selenium` | Selenium browser automation |
| `selenide` | Selenide concise UI tests |
| `playwright` | Microsoft Playwright browser automation |
| `cucumber` | Cucumber BDD testing |
| `hoverfly` | Hoverfly API simulation |
| `pact` | Pact consumer-driven contract testing |

### Maps

| Shortcut | Description |
|----------|-------------|
| `google-maps` | Google Maps Server SDK |
| `mapbox` | Mapbox Java SDK |
| `graphhopper` | GraphHopper routing engine |
| `h3` | Uber H3 Hexagonal Spatial Index |
| `jts` | JTS Topology Suite (geometry) |
| `geotools` | GeoTools geospatial data tools |
| `ip2location` | IP2Location geo lookup |
| `maxmind` | MaxMind GeoIP2 database |

### Media

| Shortcut | Description |
|----------|-------------|
| `ffmpeg` | FFmpeg video/audio processing |
| `opencv` | OpenCV computer vision |
| `pdfbox` | Apache PDFBox PDF handling |
| `tika` | Apache Tika content detection |
| `docx4j` | Docx4j Word/PowerPoint files |
| `openhtmltopdf` | OpenHTMLtoPDF converter |
| `batik` | Apache Batik SVG tools |
| `metadata-extractor` | Image metadata extraction |
| `imgscalr` | ImgScalr image scaling |
| `jxls` | JXLS Excel templates |
| `flying-saucer` | Flying Saucer XHTML to PDF |
| `barcode4j` | Barcode4J barcode generation |

### Fintech

| Shortcut | Description |
|----------|-------------|
| `web3j` | Web3j Ethereum integration |
| `bitcoinj` | BitcoinJ Bitcoin protocol |
| `plaid` | Plaid bank account linking |
| `xchange` | XChange crypto exchange library |
| `stellar` | Stellar network SDK |

### Social

| Shortcut | Description |
|----------|-------------|
| `facebook-sdk` | Facebook Graph API |
| `twitter-api` | Twitter/X API v2 |
| `linkedin-api` | LinkedIn API |
| `zoom` | Zoom API integration |
| `agora` | Agora real-time video/audio |
| `matrix` | Matrix.org decentralized chat |

### Data

| Shortcut | Description |
|----------|-------------|
| `spark` | Apache Spark Java API |
| `flink` | Apache Flink stream processing |
| `hadoop` | Hadoop client libraries |
| `tablesaw` | Tablesaw Java Dataframes |
| `commons-math` | Apache Commons Math |
| `nd4j` | ND4J N-Dimensional Arrays |
| `joda-money` | Joda-Money currency handling |
| `jgrapht` | JGraphT graph algorithms |

### Feature Flags

| Shortcut | Description |
|----------|-------------|
| `unleash` | Unleash feature toggle |
| `launchdarkly` | LaunchDarkly feature flags |
| `flagsmith` | Flagsmith feature flags |
| `togglz` | Togglz feature flags for Spring |
| `ff4j` | FF4J feature flipping |

### Microservices

| Shortcut | Description |
|----------|-------------|
| `spring-cloud-stream` | Spring Cloud Stream (Kafka/RabbitMQ) |
| `spring-cloud-bus` | Spring Cloud Bus state broadcast |
| `spring-cloud-function` | Spring Cloud Function serverless |
| `spring-retry` | Spring Retry logic |
| `shedlock` | ShedLock distributed scheduling lock |
| `etcd` | Etcd distributed key-value |
| `zookeeper` | Apache Zookeeper client |
| `dapr` | Dapr distributed runtime SDK |
| `grpc-server` | gRPC server starter |

### Integration

| Shortcut | Description |
|----------|-------------|
| `spring-integration` | Spring Integration EIP |
| `apache-camel` | Apache Camel routing |
| `spring-cloud-data-flow` | Spring Cloud Data Flow |
| `spring-cloud-task` | Spring Cloud Task ephemeral tasks |

### IoT

| Shortcut | Description |
|----------|-------------|
| `mqtt-paho` | Eclipse Paho MQTT client |
| `spring-integration-mqtt` | Spring Integration MQTT |
| `coap` | Californium CoAP protocol |
| `modbus` | Modbus industrial protocol |
| `pi4j` | Pi4J Raspberry Pi control |

### DevOps

| Shortcut | Description |
|----------|-------------|
| `jib` | Google Jib Docker builds |
| `docker-java` | Docker Java API |
| `kubernetes-client` | Kubernetes Java client |
| `fabric8` | Fabric8 Kubernetes client |

### Quality

| Shortcut | Description |
|----------|-------------|
| `jacoco` | JaCoCo code coverage |
| `checkstyle` | Checkstyle coding standards |
| `spotless` | Spotless code formatting |
| `pmd` | PMD source code analyzer |
| `sonar` | SonarQube integration |

### Caching

| Shortcut | Description |
|----------|-------------|
| `caffeine` | Caffeine high-performance cache |
| `ehcache` | Ehcache enterprise caching |
| `infinispan` | Infinispan data grid |
| `hazelcast-jet` | Hazelcast Jet stream processing |

### Content

| Shortcut | Description |
|----------|-------------|
| `rome` | Rome RSS/Atom feeds |
| `htmlunit` | HtmlUnit headless browser |
| `bliki` | Bliki Wikipedia parser |
| `emoji-java` | Emoji Java string handling |

### Networking

| Shortcut | Description |
|----------|-------------|
| `netty` | Netty async network framework |
| `pcap4j` | Pcap4J packet capture |
| `dns-java` | dnsjava DNS protocol |

### API

| Shortcut | Description |
|----------|-------------|
| `springdoc` | SpringDoc OpenAPI 3 |
| `netflix-dgs` | Netflix DGS GraphQL |
| `asyncapi` | AsyncAPI event-driven docs |
| `avro-serializer` | Avro Kafka serializer |
| `json-schema` | JSON Schema validation |
| `graphql-kickstart` | GraphQL Kickstart starter |

### Scheduling

| Shortcut | Description |
|----------|-------------|
| `jobrunr` | JobRunr distributed jobs |
| `db-scheduler` | DB Scheduler persistent tasks |

### Logging

| Shortcut | Description |
|----------|-------------|
| `logback` | Logback logging framework |
| `log4j2` | Apache Log4j 2 logging |
| `slf4j` | SLF4J logging facade |

## What Gets Added

### Maven Example: `haft add lombok`

```xml
<dependency>
    <groupId>org.projectlombok</groupId>
    <artifactId>lombok</artifactId>
    <scope>provided</scope>
</dependency>
```

### Gradle Example: `haft add lombok`

```groovy
// build.gradle (Groovy)
compileOnly 'org.projectlombok:lombok'
annotationProcessor 'org.projectlombok:lombok'
```

```kotlin
// build.gradle.kts (Kotlin)
compileOnly("org.projectlombok:lombok")
annotationProcessor("org.projectlombok:lombok")
```

### Maven Example: `haft add jwt`

Adds all three JJWT artifacts:

```xml
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-api</artifactId>
    <version>0.12.5</version>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-impl</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
<dependency>
    <groupId>io.jsonwebtoken</groupId>
    <artifactId>jjwt-jackson</artifactId>
    <version>0.12.5</version>
    <scope>runtime</scope>
</dependency>
```

### Maven Example: `haft add org.mapstruct:mapstruct`

Auto-fetches latest version from Maven Central:

```xml
<dependency>
    <groupId>org.mapstruct</groupId>
    <artifactId>mapstruct</artifactId>
    <version>1.5.5.Final</version>
</dependency>
```

### Gradle Example: `haft add org.mapstruct:mapstruct`

```groovy
// build.gradle (Groovy)
implementation 'org.mapstruct:mapstruct:1.5.5.Final'
```

```kotlin
// build.gradle.kts (Kotlin)
implementation("org.mapstruct:mapstruct:1.5.5.Final")
```

## Build Tool Detection

Haft automatically detects your build tool:

| File Found | Build Tool |
|------------|------------|
| `pom.xml` | Maven |
| `build.gradle.kts` | Gradle (Kotlin DSL) |
| `build.gradle` | Gradle (Groovy DSL) |

## Duplicate Detection

Haft automatically detects existing dependencies and skips them:

```
$ haft add lombok
WARN ⚠ Skipped (already exists) dependency=org.projectlombok:lombok
INFO ℹ No new dependencies added (all already exist)
```

## See Also

- [haft remove](/docs/commands/remove) — Remove dependencies
- [haft init](/docs/commands/init) — Add dependencies at project creation
- [Dependencies Guide](/docs/guides/dependencies) — Full dependency reference
