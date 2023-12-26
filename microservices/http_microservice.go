package microservices

import (
	"context"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// startHTTP will start HTTP service, this function will block thread
func (ms *Microservice) startHTTP(exitChannel chan bool) error {
	// Caller can exit by sending value to exitChannel
	go func() {
		<-exitChannel
		ms.stopHTTP()
	}()

	// Default middleware config
	fiberRequestID := requestid.New()
	fiberLogger := logger.New(config.LoggerConfig)
	fiberRecover := recover.New()
	fiberLimiter := limiter.New(config.LimiterConfig)
	fiberETag := etag.New(config.ETagConfig)
	fiberCors := cors.New(config.CorsConfig)
	fiberCsrf := csrf.New(config.CsrfConfig)
	fiberFavicon := favicon.New(config.FaviconConfig)

	ms.fiber.Use(fiberRequestID)
	ms.fiber.Use(fiberLogger)
	ms.fiber.Use(fiberRecover)
	ms.fiber.Use(fiberLimiter)
	ms.fiber.Use(fiberETag)
	ms.fiber.Use(fiberCors)
	ms.fiber.Use(fiberCsrf)
	ms.fiber.Use(fiberFavicon)

	return ms.fiber.Listen(":" + config.AppConfig.Port)
}

func (ms *Microservice) stopHTTP() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ms.fiber.Shutdown()
}

// For Fiber middlewares
func (ms *Microservice) Use(args ...interface{}) {
	ms.fiber.Use(args...)
}

// For Fiber route grouping
func (ms *Microservice) Group(prefix string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Group(prefix, handlers...)
}

// GET register service endpoint for HTTP GET
func (ms *Microservice) GET(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Get(path, handlers...)
}

// POST register service endpoint for HTTP POST
func (ms *Microservice) POST(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Post(path, handlers...)
}

// PUT register service endpoint for HTTP PUT
func (ms *Microservice) PUT(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Put(path, handlers...)
}

// PATCH register service endpoint for HTTP PATCH
func (ms *Microservice) PATCH(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Patch(path, handlers...)
}

// DELETE register service endpoint for HTTP DELETE
func (ms *Microservice) DELETE(path string, handlers ...func(*fiber.Ctx) error) fiber.Router {
	return ms.fiber.Delete(path, handlers...)
}
