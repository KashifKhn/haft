---
sidebar_position: 12
---

# haft generate scheduler

Generate Spring Boot scheduled task classes with `@Scheduled` annotations.

## Usage

```bash
haft generate scheduler <name> [flags]
haft g sch <name> [flags]
```

## Aliases

- `scheduler`
- `sch`
- `scheduled`
- `task`

## Description

The scheduler generator creates:

- A `@Component` class with `@Scheduled` method
- `SchedulingConfig.java` (if not exists) to enable scheduling via `@EnableScheduling`

## Schedule Types

| Type | Description | Example |
|------|-------------|---------|
| `cron` | Run at specific times using cron syntax | `0 0 8 * * *` (8 AM daily) |
| `fixedRate` | Run every N milliseconds from method start | `300000` (every 5 min) |
| `fixedDelay` | Run N milliseconds after previous completion | `60000` (1 min after done) |

## Common Cron Expressions

| Expression | Description |
|------------|-------------|
| `0 0 * * * *` | Every hour |
| `0 0 8 * * *` | Every day at 8 AM |
| `0 0 0 * * MON` | Every Monday at midnight |
| `0 */15 * * * *` | Every 15 minutes |
| `0 0 8 * * MON-FRI` | Weekdays at 8 AM |
| `0 0 2 * * *` | Every day at 2 AM |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--cron` | | Cron expression (e.g., `"0 0 * * * *"`) |
| `--rate` | | Fixed rate in milliseconds |
| `--delay` | | Fixed delay in milliseconds |
| `--initial` | | Initial delay (used with `--delay`) |
| `--package` | `-p` | Override base package |
| `--no-interactive` | | Skip interactive wizard |
| `--refresh` | | Force re-detection of project profile |
| `--json` | | Output result as JSON |

## Examples

### Interactive Mode

```bash
# Opens wizard to configure scheduler
haft generate scheduler cleanup

# Short alias
haft g sch report
```

### With Cron Expression

```bash
# Run every day at 2 AM
haft generate scheduler cleanup --cron "0 0 2 * * *"

# Run weekdays at 8 AM
haft generate scheduler report --cron "0 0 8 * * MON-FRI"

# Run every 15 minutes
haft generate scheduler sync --cron "0 */15 * * * *"
```

### With Fixed Rate

```bash
# Run every 5 minutes (300000ms)
haft generate scheduler sync --rate 300000

# Run every 30 seconds
haft generate scheduler heartbeat --rate 30000
```

### With Fixed Delay

```bash
# Run 1 minute after previous execution completes
haft generate scheduler process --delay 60000

# With initial delay of 5 seconds
haft generate scheduler batch --delay 60000 --initial 5000
```

### Non-Interactive Mode

```bash
haft generate scheduler report \
  --cron "0 0 8 * * MON-FRI" \
  --package com.example.app \
  --no-interactive
```

### JSON Output (for CI/CD)

```bash
haft generate scheduler cleanup --cron "0 0 2 * * *" --no-interactive --json
```

## Generated Files

### Layered Architecture

```
src/main/java/com/example/app/
├── config/
│   └── SchedulingConfig.java    # @EnableScheduling
└── scheduler/
    └── CleanupTask.java         # @Scheduled task
```

### Feature Architecture

```
src/main/java/com/example/app/
└── common/
    ├── config/
    │   └── SchedulingConfig.java
    └── scheduler/
        └── CleanupTask.java
```

### Hexagonal/Clean Architecture

```
src/main/java/com/example/app/
└── infrastructure/
    ├── config/
    │   └── SchedulingConfig.java
    └── scheduler/
        └── CleanupTask.java
```

## Generated Code Example

### SchedulingConfig.java

```java
package com.example.app.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.annotation.EnableScheduling;

@Configuration
@EnableScheduling
public class SchedulingConfig {
}
```

### CleanupTask.java (with Lombok)

```java
package com.example.app.scheduler;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@RequiredArgsConstructor
@Component
public class CleanupTask {

    @Scheduled(cron = "0 0 2 * * *")
    public void execute() {
        log.info("CleanupTask started");
        try {
            // TODO: Implement your scheduled task logic here
            
        } catch (Exception e) {
            log.error("CleanupTask failed", e);
        }
        log.info("CleanupTask completed");
    }
}
```

### CleanupTask.java (without Lombok)

```java
package com.example.app.scheduler;

import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Component
public class CleanupTask {

    private static final Logger log = LoggerFactory.getLogger(CleanupTask.class);

    @Scheduled(cron = "0 0 2 * * *")
    public void execute() {
        log.info("CleanupTask started");
        try {
            // TODO: Implement your scheduled task logic here
            
        } catch (Exception e) {
            log.error("CleanupTask failed", e);
        }
        log.info("CleanupTask completed");
    }
}
```

## Tips

1. **Cron format**: Spring uses 6-field cron expressions (second, minute, hour, day, month, weekday)

2. **Time zones**: Configure timezone in application properties:
   ```yaml
   spring:
     task:
       scheduling:
         pool:
           size: 5
   ```

3. **Error handling**: Always wrap task logic in try-catch to prevent task failures

4. **Logging**: Use structured logging to track task execution

5. **Long-running tasks**: Consider using `fixedDelay` instead of `fixedRate` to prevent overlapping executions
