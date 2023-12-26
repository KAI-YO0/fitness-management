package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/joho/godotenv"
)

type config struct {
	// App settings
	AppName     string
	Env         string
	Port        string
	ServiceName string
	// Fiber settings
	IsPrefork        string
	CorsAllowOrigins string
	RateLimit        int
	// Database
	DatabaseDSN          string
	DatabaseURL          string
	DatabaseHostname     string
	DatabasePort         string
	DatabaseName         string
	DatabaseUser         string
	DatabasePassword     string
	DatabaseMaxIdleConns int
	DatabaseMaxOpenConns int
	// Cache driver
	RedisAddr           string
	RedisUsername       string
	RedisPassword       string
	CachePrefix         string
	CacheMinuteDuration int
	// Sentry.io
	SentryDSN              string
	SentryEnableTracing    bool
	SentryTracesSampleRate float64
	// OpenTelemetry settings
	OtelExporterOTLPEndpoint string
	OtelInsecureMode         bool
	// OAuth Public Key
	OAuthPublicKey string
}

var (
	AppConfig     *config
	FiberConfig   fiber.Config
	ETagConfig    etag.Config
	CorsConfig    cors.Config
	CsrfConfig    csrf.Config
	LoggerConfig  logger.Config
	FaviconConfig favicon.Config
	LimiterConfig limiter.Config
	CacheConfig   cache.Config
)

// Cache environment config file
func LoadConfig() {
	godotenv.Load()

	AppConfig = &config{
		// App settings
		Env:              os.Getenv("ENV"),
		AppName:          os.Getenv("APP_NAME"),
		ServiceName:      os.Getenv("SERVICE_NAME"),
		Port:             os.Getenv("PORT"),
		IsPrefork:        os.Getenv("FIBER_PREFORK"),
		CorsAllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
		// Database
		DatabaseHostname: os.Getenv("DATABASE_HOST"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		// Redis
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		// Cache settings
		CachePrefix: os.Getenv("CACHE_PREFIX"),
		// Sentry.io
		SentryDSN: os.Getenv("SENTRY_DSN"),
		// OpenTelemetry settings
		OtelExporterOTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		// OAuth Public Key
		OAuthPublicKey: os.Getenv("OAUTH_PUBLIC_KEY"),
	}

	// Build database DSN
	AppConfig.DatabaseDSN = fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok`,
		AppConfig.DatabaseHostname,
		AppConfig.DatabaseUser,
		AppConfig.DatabasePassword,
		AppConfig.DatabaseName,
		AppConfig.DatabasePort,
	)

	// Build database URL
	AppConfig.DatabaseURL = fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		AppConfig.DatabaseUser,
		AppConfig.DatabasePassword,
		AppConfig.DatabaseHostname,
		AppConfig.DatabasePort,
		AppConfig.DatabaseName,
	)

	rateLimiter, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err == nil {
		AppConfig.RateLimit = rateLimiter
	} else {
		// Default cache duration is 60 seconds
		AppConfig.RateLimit = 60
	}

	databaseMaxIdleConns, err := strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONNS"))
	if err == nil {
		AppConfig.DatabaseMaxIdleConns = databaseMaxIdleConns
	} else {
		// Default max idle conns is 10
		AppConfig.DatabaseMaxIdleConns = 10
	}

	databaseMaxOpenConns, err := strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONNS"))
	if err == nil {
		AppConfig.DatabaseMaxOpenConns = databaseMaxOpenConns
	} else {
		// Default max open conns is 20
		AppConfig.DatabaseMaxIdleConns = 20
	}

	cacheMinuteDuration, err := strconv.Atoi(os.Getenv("CACHE_MINUTE_DURATION"))
	if err == nil {
		AppConfig.CacheMinuteDuration = cacheMinuteDuration
	} else {
		// Default cache duration is 5 minutes
		AppConfig.CacheMinuteDuration = 5
	}

	var isPrefork bool
	if AppConfig.IsPrefork == "true" {
		isPrefork = true
	} else {
		isPrefork = false
	}

	FiberConfig = fiber.Config{
		Prefork:                 isPrefork,
		CaseSensitive:           true,
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
		ServerHeader:            "Fiber",
		AppName:                 AppConfig.AppName,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
	}

	ETagConfig = etag.Config{
		Weak: true,
	}

	LoggerConfig = logger.Config{
		Format: "[${ip}]:${port} (${pid}) ${locals:requestid} ${status} - ${method} ${path}\n",
	}

	FaviconConfig = favicon.Config{
		Next: nil,
		File: "",
	}

	LimiterConfig = limiter.Config{
		Max: AppConfig.RateLimit,
		Next: func(c *fiber.Ctx) bool {
			return c.Query("loadtest") == "true"
		},
	}

	CacheConfig = cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		CacheHeader:  "X-Cache",
		Expiration:   time.Duration(AppConfig.CacheMinuteDuration) * time.Minute,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path())
		},
	}

	// Set CORS default to allow all
	if AppConfig.CorsAllowOrigins == "" {
		AppConfig.CorsAllowOrigins = "*"
	}

	CorsConfig = cors.Config{
		Next:             nil,
		AllowOrigins:     AppConfig.CorsAllowOrigins,
		AllowCredentials: true,
	}

	sentryEnableTracing, err := strconv.ParseBool(os.Getenv("SENTRY_ERROR_TRACING"))
	if err == nil {
		AppConfig.SentryEnableTracing = sentryEnableTracing
	} else {
		// Default is false
		AppConfig.SentryEnableTracing = false
	}

	sentryTracesSampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_TRACES_SAMPLE_RATE"), 64)
	if err == nil {
		AppConfig.SentryTracesSampleRate = sentryTracesSampleRate
	} else {
		// Default traces sample rate is 0.2
		AppConfig.SentryTracesSampleRate = 0.2
	}

	otelInsecureMode, err := strconv.ParseBool(os.Getenv("OTEL_INSECURE_MODE"))
	if err == nil {
		AppConfig.OtelInsecureMode = otelInsecureMode
	} else {
		// Default is false
		AppConfig.OtelInsecureMode = false
	}
}
