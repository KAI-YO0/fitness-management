# Stream HTTP Service Template

[Stream] Microservices REST API template using Go Fiber Framework + exception tracking by Sentry.io

## Extra feature includes

1. Object relational mapping (ORM)
2. In-memory Caching
3. Error exception handling
4. Observability (tracing)

## Documentation

- Fiber Framework (<https://docs.gofiber.io/>)
- GORM (<https://gorm.io/>)
- Golang Migrate (<https://github.com/golang-migrate/migrate>)
- Go Redis (<https://redis.uptrace.dev/guide/go-redis.html>)
- Sentry (<https://docs.sentry.io/platforms/go/>)
- OpenTelemetry (<https://opentelemetry.io/docs/instrumentation/go/>)

## System requirements

- Golang 1.21
- PostgreSQL Database 15.x
- Redis 6

## Install dependencies

```bash
go install
```

## Update dependencies

```bash
go mod tidy
```

---

## Start development server

```bash
go run .
```

The server will listen by default at <http://localhost:8000>

---

## Build go executable file

```bash
go build -o main
```

## Run database migration, rollback or seeder

```bash
go run . --db-migrate
go run . --db-rollback
go run . --db-seed
```

Run with executable file

```bash
./main --db-migrate
./main --db-rollback
./main --db-seed
```

---

## Environment Variables

```env
ENV="local"
APP_NAME="Stream - HTTP service"
SERVICE_NAME="UserService"
PORT=8000

FIBER_PREFORK=true

DATABASE_HOST="localhost"
DATABASE_PORT=5432
DATABASE_NAME="postgres"
DATABASE_USER="postgres"
DATABASE_PASSWORD="mysecretpassword"
DATABASE_MAX_IDLE_CONNS=10
DATABASE_MAX_OPEN_CONNS=20

REDIS_ADDR="127.0.0.1:6379"
REDIS_USERNAME="default"
REDIS_PASSWORD=""
CACHE_MINUTE_DURATION=15

SENTRY_DSN=""
SENTRY_ERROR_TRACING=false
SENTRY_TRACES_SAMPLE_RATE=0.2

OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
OTEL_INSECURE_MODE=true

OAUTH_PUBLIC_KEY=""
 ```

---
