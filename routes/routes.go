package routes

import (
	"github.com/Stream-I-T-Consulting/stream-http-service-go/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/handlers"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/microservices"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/repositories"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func HTTPRootMiddleware(ms *microservices.Microservice) {
	// Default middleware config
	fiberRequestID := requestid.New()
	fiberLogger := logger.New(config.LoggerConfig)
	fiberRecover := recover.New()
	fiberLimiter := limiter.New(config.LimiterConfig)
	fiberETag := etag.New(config.ETagConfig)
	fiberCors := cors.New(config.CorsConfig)
	fiberFavicon := favicon.New(config.FaviconConfig)

	ms.Use(fiberRequestID)
	ms.Use(fiberLogger)
	ms.Use(fiberRecover)
	ms.Use(fiberLimiter)
	ms.Use(fiberETag)
	ms.Use(fiberCors)
	ms.Use(fiberFavicon)
}

func HTTPRootRoute(ms *microservices.Microservice) {
	HTTPRootMiddleware(ms)

	ms.GET("/", func(c *fiber.Ctx) error {
		return handlers.GetRootPath(c)
	})
	ms.GET("/health", func(c *fiber.Ctx) error {
		return handlers.GetHealthCheck(c)
	})
	ms.GET("/monitor", monitor.New(monitor.Config{Title: "Fiber"}))
}

func HTTPRoutes(ms *microservices.Microservice) {
	// Initialize repositories, services, and handlers
	userRepo := repositories.NewUserRepository(database.DBConn)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	handler := handlers.NewHandler(
		cache.Cacher,
		userService,
	)

	// REST API endpoint ------------------------------------------------------------------

	// Example route grouping
	api := ms.Group("api")
	apiV1 := api.Group("v1")

	// User service routes
	apiV1.Get("/users", func(c *fiber.Ctx) error { return handler.GetUsers(c) })
	apiV1.Get("/users/:id", func(c *fiber.Ctx) error { return handler.GetUser(c) })
	apiV1.Post("/users", func(c *fiber.Ctx) error { return handler.CreateUser(c) })
	apiV1.Put("/users/:id", func(c *fiber.Ctx) error { return handler.UpdateUser(c) })
	apiV1.Delete("/users/:id", func(c *fiber.Ctx) error { return handler.DeleteUser(c) })
}
