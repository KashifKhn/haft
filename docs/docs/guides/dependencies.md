---
sidebar_position: 2
title: Dependencies
description: Available Spring Boot dependencies
---

# Dependencies

Haft provides access to all official Spring Boot starters from Spring Initializr, organized by category.

## Categories

### Developer Tools

Development-time tools and utilities.

| Dependency | Description |
|------------|-------------|
| `devtools` | Spring Boot DevTools for fast restarts |
| `lombok` | Java annotation library for boilerplate reduction |
| `configuration-processor` | Generate metadata for configuration properties |
| `docker-compose` | Docker Compose support |

### Web

Web application frameworks and tools.

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

Relational database support.

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

Non-relational database support.

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

Authentication and authorization.

| Dependency | Description |
|------------|-------------|
| `security` | Spring Security |
| `oauth2-client` | OAuth2 Client |
| `oauth2-resource-server` | OAuth2 Resource Server |
| `oauth2-authorization-server` | OAuth2 Authorization Server |
| `ldap` | LDAP authentication |

### Messaging

Message brokers and queues.

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

Cloud-native features.

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

Monitoring and metrics.

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

Testing frameworks and utilities.

| Dependency | Description |
|------------|-------------|
| `testcontainers` | Testcontainers support |
| `cloud-contract-verifier` | Spring Cloud Contract Verifier |
| `cloud-contract-stub-runner` | Spring Cloud Contract Stub Runner |
| `restdocs` | Spring REST Docs |

### I/O

Input/output and integrations.

| Dependency | Description |
|------------|-------------|
| `batch` | Spring Batch |
| `integration` | Spring Integration |
| `mail` | Java Mail Sender |
| `quartz` | Quartz Scheduler |
| `cache` | Spring Cache |
| `validation` | Bean Validation (Hibernate Validator) |

## Using Dependencies

### With `haft add`

The easiest way to add dependencies to an existing project:

```bash
# Interactive search picker
haft add

# Browse by category
haft add --browse

# Add by shortcut
haft add lombok jpa validation

# Add by Maven coordinates (auto-verified)
haft add org.mapstruct:mapstruct

# List all shortcuts
haft add --list
```

### With `haft remove`

Remove dependencies from your project:

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

Use the dependency ID (shown in the left column above).

## Popular Combinations

### REST API

```bash
--deps web,data-jpa,validation,lombok
```

### Microservice

```bash
--deps web,actuator,cloud-config-client,cloud-eureka-client
```

### Reactive Application

```bash
--deps webflux,data-mongodb-reactive,security
```

### Batch Processing

```bash
--deps batch,data-jpa,postgresql
```
