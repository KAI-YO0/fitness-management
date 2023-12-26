package handlers

import (
	"context"
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/services"
	"github.com/gofiber/fiber/v2"
)

type (
	// Register handler services
	handler struct {
		cacher      *cache.Cache
		userService services.UserService
	}
	// Register handler interfaces
	Handler interface {
		UserHandler
	}
)

func NewHandler(
	cacher *cache.Cache,
	userService services.UserService,
) handler {
	return handler{
		cacher:      cacher,
		userService: userService,
	}
}

type ServicePaginationFunc func(ctx context.Context, paginate database.Pagination, search string) (*database.Pagination, error)
type ServiceQueryFunc func(ctx context.Context, id int) (map[string]interface{}, error)

func (h handler) PaginationCache(ctx context.Context, key string, tags []string, paginate database.Pagination, search string, f ServicePaginationFunc) (*database.Pagination, error) {
	var (
		responseData *database.Pagination
		err          error
	)

	// Get the cached attributes object
	err = cache.Cacher.Get(ctx, key, &responseData)
	if err != nil {
		return nil, err
	}

	if responseData == nil {
		// Call service function
		responseData, err = f(ctx, paginate, search)
		if err != nil {
			return nil, err
		}

		// Set cache
		err = cache.Cacher.Tag(tags...).Set(ctx, key, &responseData)
		if err != nil {
			return nil, err
		}
	}

	return responseData, nil
}

func (h handler) QueryCache(ctx context.Context, key string, tags []string, id int, f ServiceQueryFunc) (map[string]interface{}, error) {
	var (
		responseData map[string]interface{}
		err          error
	)

	// Get the cached attributes object
	err = cache.Cacher.Get(ctx, key, &responseData)
	if err != nil {
		return nil, err
	}

	if responseData == nil {
		// Call service function
		responseData, err = f(ctx, id)
		if err != nil {
			return nil, err
		}

		// Set cache
		err = cache.Cacher.Tag(tags...).Set(ctx, key, &responseData)
		if err != nil {
			return nil, err
		}
	}

	return responseData, nil
}

// Root handlers  ------------------------------------------------------------------

func GetRootPath(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString(config.AppConfig.AppName)
}

func GetHealthCheck(c *fiber.Ctx) error {
	sqlDB, _ := database.DBConn.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("DatabaseError:", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
