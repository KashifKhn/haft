---
sidebar_position: 2
title: Dependencies
description: Available Spring Boot dependencies and shortcuts
---

# Dependencies

Haft provides access to all official Spring Boot starters plus 330+ additional dependency shortcuts for popular libraries across 35 categories.

## Build Tool Support

Haft automatically detects your build tool and modifies the appropriate file:

| File Found | Build Tool | Dependency Format |
|------------|------------|-------------------|
| `pom.xml` | Maven | XML dependencies |
| `build.gradle.kts` | Gradle (Kotlin DSL) | Kotlin DSL syntax |
| `build.gradle` | Gradle (Groovy DSL) | Groovy DSL syntax |

## Spring Initializr Dependencies

These are the official Spring Boot starters available during `haft init`:

### Developer Tools

| Dependency | Description |
|------------|-------------|
| `devtools` | Spring Boot DevTools for fast restarts |
| `lombok` | Java annotation library for boilerplate reduction |
| `configuration-processor` | Generate metadata for configuration properties |
| `docker-compose` | Docker Compose support |

### Web

| Dependency | Description |
|------------|-------------|
| `web` | Spring Web MVC for RESTful applications |
| `webflux` | Reactive web framework |
| `graphql` | Spring GraphQL support |
| `rest-docs` | Document RESTful APIs |
| `hateoas` | Hypermedia-driven RESTful services |
| `web-services` | SOAP web services |
| `jersey` | JAX-RS with Jersey |
| `vaadin` | Vaadin Flow web framework |
| `thymeleaf` | Thymeleaf template engine |
| `freemarker` | FreeMarker template engine |
| `mustache` | Mustache template engine |
| `groovy-templates` | Groovy template engine |

### SQL Databases

| Dependency | Description |
|------------|-------------|
| `data-jpa` | Spring Data JPA with Hibernate |
| `data-jdbc` | Spring Data JDBC |
| `jdbc` | Spring JDBC |
| `h2` | H2 Database (embedded) |
| `mysql` | MySQL Driver |
| `postgresql` | PostgreSQL Driver |
| `mariadb` | MariaDB Driver |
| `oracle` | Oracle Driver |
| `sqlserver` | Microsoft SQL Server Driver |
| `flyway` | Flyway database migrations |
| `liquibase` | Liquibase database migrations |
| `jooq` | jOOQ SQL DSL |

### NoSQL Databases

| Dependency | Description |
|------------|-------------|
| `data-mongodb` | Spring Data MongoDB |
| `data-mongodb-reactive` | Reactive MongoDB |
| `data-redis` | Spring Data Redis |
| `data-redis-reactive` | Reactive Redis |
| `data-elasticsearch` | Spring Data Elasticsearch |
| `data-cassandra` | Spring Data Cassandra |
| `data-cassandra-reactive` | Reactive Cassandra |
| `data-couchbase` | Spring Data Couchbase |
| `data-couchbase-reactive` | Reactive Couchbase |
| `data-neo4j` | Spring Data Neo4j |

### Security

| Dependency | Description |
|------------|-------------|
| `security` | Spring Security |
| `oauth2-client` | OAuth2 Client |
| `oauth2-resource-server` | OAuth2 Resource Server |
| `oauth2-authorization-server` | OAuth2 Authorization Server |
| `ldap` | LDAP authentication |

### Messaging

| Dependency | Description |
|------------|-------------|
| `kafka` | Apache Kafka |
| `kafka-streams` | Kafka Streams |
| `amqp` | RabbitMQ (AMQP) |
| `artemis` | Apache ActiveMQ Artemis |
| `pulsar` | Apache Pulsar |
| `pulsar-reactive` | Reactive Pulsar |
| `websocket` | WebSocket support |
| `rsocket` | RSocket |

### Cloud

| Dependency | Description |
|------------|-------------|
| `cloud-config-client` | Config Client |
| `cloud-config-server` | Config Server |
| `cloud-eureka-client` | Eureka Discovery Client |
| `cloud-eureka-server` | Eureka Discovery Server |
| `cloud-gateway` | Spring Cloud Gateway |
| `cloud-gateway-mvc` | Spring Cloud Gateway MVC |
| `cloud-circuitbreaker-resilience4j` | Circuit Breaker (Resilience4j) |
| `cloud-starter-consul-discovery` | Consul Discovery |
| `cloud-starter-vault-config` | Vault Configuration |
| `cloud-starter-zookeeper-discovery` | Zookeeper Discovery |

### Observability

| Dependency | Description |
|------------|-------------|
| `actuator` | Spring Boot Actuator |
| `datadog` | Datadog metrics |
| `graphite` | Graphite metrics |
| `influx` | InfluxDB metrics |
| `new-relic` | New Relic metrics |
| `prometheus` | Prometheus metrics |
| `wavefront` | Wavefront metrics |
| `zipkin` | Zipkin distributed tracing |
| `distributed-tracing` | Micrometer Tracing |

### Testing

| Dependency | Description |
|------------|-------------|
| `testcontainers` | Testcontainers support |
| `cloud-contract-verifier` | Spring Cloud Contract Verifier |
| `cloud-contract-stub-runner` | Spring Cloud Contract Stub Runner |
| `restdocs` | Spring REST Docs |

### I/O

| Dependency | Description |
|------------|-------------|
| `batch` | Spring Batch |
| `integration` | Spring Integration |
| `mail` | Java Mail Sender |
| `quartz` | Quartz Scheduler |
| `cache` | Spring Cache |
| `validation` | Bean Validation (Hibernate Validator) |

---

## Extended Shortcuts (330+)

Beyond Spring Initializr starters, `haft add` supports 330+ shortcuts for popular libraries:

### Web & HTTP

| Shortcut | Description |
|----------|-------------|
| `feign` | Declarative REST client (OpenFeign) |
| `resilience4j` | Fault tolerance (circuit breaker) |
| `restclient` | Synchronous HTTP client |
| `webclient` | Reactive HTTP client |
| `okhttp` | OkHttp HTTP client |
| `retrofit` | Retrofit type-safe REST client |
| `httpclient5` | Apache HttpClient 5 |
| `session-redis` | Distributed sessions with Redis |
| `session-jdbc` | Distributed sessions with JDBC |
| `grpc` | gRPC framework for Spring Boot |
| `bucket4j` | Rate limiting library |
| `spring-cloud-gateway` | API Gateway for microservices |
| `eureka` | Service discovery with Eureka |
| `consul` | Service discovery with Consul |
| `config-server` | Centralized configuration server |
| `config-client` | Config server client |

### SQL Databases (Extended)

| Shortcut | Description |
|----------|-------------|
| `jpa` | Spring Data JPA (alias) |
| `mybatis` | MyBatis SQL mapping framework |
| `r2dbc` | Reactive database access |
| `r2dbc-postgresql` | Reactive PostgreSQL driver |
| `clickhouse` | ClickHouse column-oriented OLAP |
| `cockroachdb` | CockroachDB distributed SQL |
| `timescaledb` | TimescaleDB time-series |
| `voltdb` | VoltDB in-memory SQL |
| `sqlite` | SQLite embedded database |
| `hsqldb` | HSQLDB embedded database |
| `derby` | Apache Derby embedded database |
| `duckdb` | DuckDB analytical database |
| `questdb` | QuestDB time-series database |

### NoSQL (Extended)

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

### Security (Extended)

| Shortcut | Description |
|----------|-------------|
| `jwt` | JJWT library (api + impl + jackson) |
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

### Messaging (Extended)

| Shortcut | Description |
|----------|-------------|
| `rabbitmq` | RabbitMQ Java client (direct) |
| `activemq` | ActiveMQ JMS messaging |
| `nats` | NATS messaging system |
| `zeromq` | ZeroMQ high-performance messaging |

### AI & Machine Learning

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

### Cloud Providers

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

### Observability (Extended)

| Shortcut | Description |
|----------|-------------|
| `opentelemetry` | OpenTelemetry distributed tracing |
| `jaeger` | Distributed tracing with Jaeger |
| `sentry` | Sentry error tracking |
| `loki` | Grafana Loki logging |
| `grafana` | Grafana LGTM metrics |
| `micrometer` | Prometheus metrics exporter |

### Developer Tools (Extended)

| Shortcut | Description |
|----------|-------------|
| `mapstruct` | MapStruct bean mapping |
| `openapi` | SpringDoc OpenAPI (Swagger UI) |
| `commons-lang` | Apache Commons Lang |
| `guava` | Google Guava |
| `modelmapper` | ModelMapper |
| `config-processor` | IDE support for @ConfigurationProperties |
| `native` | GraalVM native image support |
| `vavr` | Functional programming library |
| `immutables` | Immutable object generation |
| `record-builder` | Builder pattern for Java records |
| `jmolecules` | DDD architectural concepts |
| `spotbugs` | SpotBugs static analysis |
| `error-prone` | Error Prone static analysis |
| `checker-qual` | Checker Framework annotations |

### Testing (Extended)

| Shortcut | Description |
|----------|-------------|
| `test` | Spring Boot Test |
| `security-test` | Spring Security Test |
| `mockito` | Mockito mocking |
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

### Utilities

| Shortcut | Description |
|----------|-------------|
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
| `commons-io` | Apache Commons IO |
| `jackson-datatype` | Jackson Java 8 datatypes |
| `apache-poi` | Microsoft Office file handling |
| `itext` | PDF generation and manipulation |
| `minio` | S3-compatible object storage client |
| `json-path` | JSON path query library |

### Workflow

| Shortcut | Description |
|----------|-------------|
| `camunda` | Camunda workflow automation |
| `flowable` | Flowable BPMN engine |
| `temporal` | Temporal workflow orchestration |

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

### Caching

| Shortcut | Description |
|----------|-------------|
| `caffeine` | Caffeine high-performance cache |
| `ehcache` | Ehcache enterprise caching |
| `infinispan` | Infinispan data grid |
| `hazelcast-jet` | Hazelcast Jet stream processing |

### Scheduling

| Shortcut | Description |
|----------|-------------|
| `jobrunr` | JobRunr distributed jobs |
| `db-scheduler` | DB Scheduler persistent tasks |
| `scheduler` | Task scheduling with @Scheduled |

### API Documentation

| Shortcut | Description |
|----------|-------------|
| `springdoc` | SpringDoc OpenAPI 3 |
| `netflix-dgs` | Netflix DGS GraphQL |
| `asyncapi` | AsyncAPI event-driven docs |
| `json-schema` | JSON Schema validation |
| `graphql-kickstart` | GraphQL Kickstart starter |

### Logging

| Shortcut | Description |
|----------|-------------|
| `logback` | Logback logging framework |
| `log4j2` | Apache Log4j 2 logging |
| `slf4j` | SLF4J logging facade |

### Quality

| Shortcut | Description |
|----------|-------------|
| `jacoco` | JaCoCo code coverage |
| `checkstyle` | Checkstyle coding standards |
| `spotless` | Spotless code formatting |
| `pmd` | PMD source code analyzer |
| `sonar` | SonarQube integration |

### DevOps

| Shortcut | Description |
|----------|-------------|
| `jib` | Google Jib Docker builds |
| `docker-java` | Docker Java API |
| `kubernetes-client` | Kubernetes Java client |
| `fabric8` | Fabric8 Kubernetes client |

### Feature Flags

| Shortcut | Description |
|----------|-------------|
| `unleash` | Unleash feature toggle |
| `launchdarkly` | LaunchDarkly feature flags |
| `flagsmith` | Flagsmith feature flags |
| `togglz` | Togglz feature flags for Spring |
| `ff4j` | FF4J feature flipping |

### Maps & Geo

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

### Data Processing

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

### IoT

| Shortcut | Description |
|----------|-------------|
| `mqtt-paho` | Eclipse Paho MQTT client |
| `spring-integration-mqtt` | Spring Integration MQTT |
| `coap` | Californium CoAP protocol |
| `modbus` | Modbus industrial protocol |
| `pi4j` | Pi4J Raspberry Pi control |

### Ops

| Shortcut | Description |
|----------|-------------|
| `actuator` | Spring Boot Actuator for monitoring |
| `micrometer` | Micrometer Prometheus metrics exporter |

### Template Engines

| Shortcut | Description |
|----------|-------------|
| `thymeleaf` | Thymeleaf server-side HTML templating |
| `freemarker` | FreeMarker template engine |
| `mustache` | Mustache template engine |

### Content

| Shortcut | Description |
|----------|-------------|
| `rome` | Rome RSS and Atom feed parser |
| `htmlunit` | HtmlUnit headless browser |
| `bliki` | Bliki Wikipedia syntax parser |
| `emoji-java` | Emoji Java string handling |

### Networking

| Shortcut | Description |
|----------|-------------|
| `netty` | Netty async event-driven network framework |
| `pcap4j` | Pcap4J packet capture library |
| `dns-java` | dnsjava DNS protocol implementation |

### Integration

| Shortcut | Description |
|----------|-------------|
| `spring-integration` | Spring Integration EIP patterns |
| `apache-camel` | Apache Camel integration framework |
| `spring-cloud-data-flow` | Spring Cloud Data Flow orchestration |
| `spring-cloud-task` | Spring Cloud Task short-lived microservices |

### API

| Shortcut | Description |
|----------|-------------|
| `springdoc` | SpringDoc OpenAPI 3 generation |
| `asyncapi` | AsyncAPI event-driven documentation |
| `netflix-dgs` | Netflix DGS GraphQL framework |
| `avro-serializer` | Avro serialization for Kafka schemas |
| `json-schema` | JSON Schema validation |
| `graphql-kickstart` | GraphQL Kickstart starter |

---

## Using Dependencies

### With `haft add`

```bash
# Interactive search picker
haft add

# Browse by category
haft add --browse

# Add by shortcut
haft add lombok jpa validation

# Add multiple at once
haft add jwt openapi mapstruct

# Add by Maven coordinates (auto-verified against Maven Central)
haft add org.mapstruct:mapstruct

# Add with specific version
haft add io.jsonwebtoken:jjwt-api:0.12.5

# Add with scope override
haft add h2 --scope test

# List all shortcuts
haft add --list
```

### Adding Any Maven Central Dependency

You can add **any dependency from Maven Central** using Maven coordinates format:

```bash
# Format: groupId:artifactId
haft add org.apache.commons:commons-collections4

# Format: groupId:artifactId:version
haft add com.google.code.gson:gson:2.10.1

# Multiple coordinates
haft add org.modelmapper:modelmapper io.github.cdimascio:dotenv-java
```

Haft automatically:
- Verifies the dependency exists on Maven Central
- Fetches the latest version if not specified
- Adds to your `pom.xml` or `build.gradle` with correct format

```bash
# This will fetch latest version automatically
$ haft add org.mapstruct:mapstruct
SUCCESS ✓ Added dependency=org.mapstruct:mapstruct:1.5.5.Final

# Error if dependency doesn't exist
$ haft add com.fake:nonexistent
ERROR ✗ dependency 'com.fake:nonexistent' not found on Maven Central
```

### With `haft remove`

```bash
# Interactive picker
haft remove

# Remove by artifact name
haft remove lombok

# Remove by suffix (jpa matches spring-boot-starter-data-jpa)
haft rm jpa web

# Remove by coordinates
haft rm org.projectlombok:lombok
```

### In Interactive Mode (`haft init`)

1. Navigate to the Dependencies step
2. Use `/` to search
3. Use `Tab` or `0-9` for categories
4. Press `Space` to select
5. Press `Enter` to confirm

### In Non-Interactive Mode

```bash
haft init my-app --deps web,data-jpa,lombok,validation
```

---

## Popular Combinations

### REST API

```bash
haft init my-api --deps web,data-jpa,validation,lombok
haft add jwt openapi
```

### Microservice

```bash
haft init my-service --deps web,actuator,cloud-config-client,cloud-eureka-client
haft add resilience4j opentelemetry
```

### Reactive Application

```bash
haft init my-reactive --deps webflux,data-mongodb-reactive,security
```

### AI Application

```bash
haft init my-ai-app --deps web,lombok
haft add openai langchain4j pgvector
```

### Batch Processing

```bash
haft init my-batch --deps batch,data-jpa,postgresql
haft add jobrunr
```

### Full-Stack with Auth

```bash
haft init my-fullstack --deps web,data-jpa,security,thymeleaf,lombok
haft add jwt passay
```
