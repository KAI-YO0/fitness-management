package exceptions

import (
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/utils/color"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

func SentryInitialize() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.AppConfig.SentryDSN,
		EnableTracing:    config.AppConfig.SentryEnableTracing,
		TracesSampleRate: config.AppConfig.SentryTracesSampleRate,
	})

	// If there is an error, do not continue.
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// Check if Sentry DSN is already set
	if config.AppConfig.SentryDSN != "" {
		if !fiber.IsChild() {
			log.Println("Sentry: Error handling is", color.Format(color.GREEN, "on!"))
			if config.AppConfig.SentryEnableTracing {
				log.Println("Sentry: Tracing is", color.Format(color.GREEN, "on!"))
				log.Println("Sentry: Tracing sample rate is", config.AppConfig.SentryTracesSampleRate)
			}
		}
	}
}
