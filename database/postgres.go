package database

import (
	"log"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/utils/color"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	DBConn *gorm.DB
)

func Initialize() (DBConn *gorm.DB) {
	DBConn, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.AppConfig.DatabaseDSN,
			},
		),
		&gorm.Config{},
	)
	DBConn.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{},
			Replicas:          []gorm.Dialector{},
			Policy:            nil,
			TraceResolverMode: false,
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(config.AppConfig.DatabaseMaxIdleConns).
			SetMaxOpenConns(config.AppConfig.DatabaseMaxOpenConns),
	)

	if err != nil {
		log.Printf("Cannot connect to database")
		log.Fatal("DatabaseError:", err)
	}

	if !fiber.IsChild() {
		log.Println("Database connected", color.Format(color.GREEN, "successfully!"))
	}

	return
}

func Migrate() {
	// Migrate the schema
	migrator, err := migrate.New(
		"file://database/migrations",
		config.AppConfig.DatabaseURL)

	// Check error when create new migrate instance
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	// Run the migrations
	log.Println("Migrating the schema...")
	if err := migrator.Up(); err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}

	log.Println("Migration completed, Close the database connection...")

	// Close the database connection
	migrator.Close()
}

func Rollback() {
	// Migrate the schema
	migrator, err := migrate.New(
		"file://database/migrations",
		config.AppConfig.DatabaseURL)

	// Check error when create new migrate instance
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	// Run the migrations
	log.Println("Rollback the schema...")
	if err := migrator.Down(); err != nil {
		sentry.CaptureException(err)
		log.Println(err)
	}

	log.Println("Rollback completed, Close the database connection...")

	// Close the database connection
	migrator.Close()
}
