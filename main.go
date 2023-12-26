package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database/seeders"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/microservices"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/exceptions"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/routes"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Migration flag
	dbMigratePtr := flag.Bool("db-migrate", false, "a bool")
	dbRollbackPtr := flag.Bool("db-rollback", false, "a bool")
	dbSeederPtr := flag.Bool("db-seed", false, "a bool")
	flag.Parse()

	// TODO: Load environment variables
	config.LoadConfig()

	// Run database migration, if -db-migrate flag is set to TRUE. and exit the application
	if *dbMigratePtr {
		database.Migrate()
		os.Exit(0)
	}
	// Run database rollback, if -db-rollback flag is set to TRUE. and exit the application
	if *dbRollbackPtr {
		database.Rollback()
		os.Exit(0)
	}
	// Run database seeder, if -db-seed flag is set to TRUE. and exit the application
	if *dbSeederPtr {
		seeders.RunSeed()
		os.Exit(0)
	}

	// Initialize connection to database and cache
	database.DBConn = database.Initialize()
	cache.Client = cache.Initialize()
	cache.Cacher = cache.NewCache(
		cache.WithPrefix(config.AppConfig.CachePrefix),
		cache.WithExpired(time.Minute*time.Duration(config.AppConfig.CacheMinuteDuration)),
	)

	// Initialize Sentry client for error logging and tracing
	exceptions.SentryInitialize()
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	// Initialize OpenTelemetry tracing
	tracing.TraceProvider = tracing.InitTracer()
	defer tracing.Cleanup()

	// Create microservice instance
	ms := microservices.NewMicroservice()

	// HTTP Routes setup
	routes.HTTPRootRoute(ms)

	// Start http server
	startHTTP(ms)

	// Microservice start up
	ms.Start()
}

func startHTTP(ms *microservices.Microservice) {
	zone, _ := time.Now().Zone()
	if !fiber.IsChild() {
		log.Println("HTTP service is running")
		log.Println("[Timezone]:", zone)
	}
	routes.HTTPRoutes(ms)
}
