package handlers

import (
	"fmt"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/cache"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/services"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/utils"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	UserHandler interface {
		// User handlers
		GetUsers(c *fiber.Ctx) error
		GetUser(c *fiber.Ctx) error
		CreateUser(c *fiber.Ctx) error
		UpdateUser(c *fiber.Ctx) error
		DeleteUser(c *fiber.Ctx) error
	}
)

func (h handler) GetUsers(c *fiber.Ctx) error {
	var (
		ctx, span    = tracing.Tracer.Start(c.Context(), "GetUsersHandler", trace.WithAttributes(attribute.String("handler", "GetUsers")))
		responseData *database.Pagination
	)

	// Get paginate values
	paginate := database.Pagination{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 20),
	}
	search := c.Query("search")

	// Make cache key
	cacheTags := []string{"users"}
	cacheKey := fmt.Sprintf("GetUsers_%d_%d", paginate.Page, paginate.Limit)
	if search != "" {
		cacheKey = fmt.Sprintf(`%s_%s`, cacheKey, search)
	}

	responseData, err := h.PaginationCache(ctx, cacheKey, cacheTags, paginate, search, h.userService.GetUsers)
	if err != nil {
		utils.HandleErrors(err)
		return fiber.ErrInternalServerError
	}

	span.End()
	return c.JSON(responseData)
}

func (h handler) GetUser(c *fiber.Ctx) error {
	var (
		id, _        = c.ParamsInt("id")
		ctx, span    = tracing.Tracer.Start(c.Context(), "GetUserHandler", trace.WithAttributes(attribute.String("handler", "GetUser"), attribute.Int("id", id)))
		responseData map[string]interface{}
	)

	// Make cache key
	cacheTags := []string{"users"}
	cacheKey := fmt.Sprintf("GetUser_%d", id)

	responseData, err := h.QueryCache(ctx, cacheKey, cacheTags, id, h.userService.GetUser)
	if err != nil {
		utils.HandleErrors(err)
		return fiber.ErrInternalServerError
	}

	span.End()
	return c.JSON(responseData)
}

func (h handler) CreateUser(c *fiber.Ctx) error {
	var (
		ctx, span = tracing.Tracer.Start(c.Context(), "CreateUserHandler", trace.WithAttributes(attribute.String("handler", "CreateUser")))
	)

	// Create data transfer object
	userDto := new(services.UserDto)

	// Parse HTTP request body to struct variable
	if err := c.BodyParser(userDto); err != nil {
		return err
	}

	// Form request validation
	errors := utils.Validate(*userDto)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Call service function
	err := h.userService.CreateUser(ctx, userDto)
	if err != nil {
		utils.HandleErrors(err)
		return err
	}

	// Clear user cache
	cache.Cacher.Tag("users").Flush(ctx)

	span.End()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    "0",
		"message": "OK",
	})
}

func (h handler) UpdateUser(c *fiber.Ctx) error {
	var (
		id, _     = c.ParamsInt("id")
		ctx, span = tracing.Tracer.Start(c.Context(), "UpdateUserHandler", trace.WithAttributes(attribute.String("handler", "UpdateUser"), attribute.Int("id", id)))
	)

	// Create data transfer object
	userDto := new(services.UserDto)

	// Parse HTTP request body to struct variable
	if err := c.BodyParser(userDto); err != nil {
		return err
	}

	// Form request validation
	errors := utils.Validate(*userDto)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Call service function
	err := h.userService.UpdateUser(ctx, id, userDto)
	if err != nil {
		utils.HandleErrors(err)
		return err
	}

	// Clear user cache
	cache.Cacher.Tag("users").Flush(ctx)

	span.End()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    "0",
		"message": "OK",
	})
}

func (h handler) DeleteUser(c *fiber.Ctx) error {
	var (
		id, _     = c.ParamsInt("id")
		ctx, span = tracing.Tracer.Start(c.Context(), "DeleteUserHandler", trace.WithAttributes(attribute.String("handler", "DeleteUser"), attribute.Int("id", id)))
	)

	// Call service function
	err := h.userService.DeleteUser(ctx, id)
	if err != nil {
		utils.HandleErrors(err)
		return err
	}

	// Clear user cache
	cache.Cacher.Tag("users").Flush(ctx)

	span.End()
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"code":    "0",
		"message": "OK",
	})
}
