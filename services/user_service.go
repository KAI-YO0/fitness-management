package services

import (
	"context"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/models"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/repositories"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	userService struct {
		userRepository repositories.UserRepository
	}
)

func NewUserService(
	userRepo repositories.UserRepository,
) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (s userService) GetUsers(ctx context.Context, paginate database.Pagination, search string) (*database.Pagination, error) {
	ctx, childSpan := tracing.Tracer.Start(ctx, "GetUsersService", trace.WithAttributes(attribute.String("service", "GetUsers")))
	result, err := s.userRepository.GetUserPaginate(ctx, paginate, search)
	childSpan.End()

	return result, err
}

func (s userService) GetUser(ctx context.Context, id int) (map[string]interface{}, error) {
	ctx, childSpan := tracing.Tracer.Start(ctx, "GetUserService", trace.WithAttributes(attribute.String("service", "GetUser")))
	user, err := s.userRepository.GetUserByID(ctx, id)
	childSpan.End()

	return map[string]interface{}{"data": user}, err
}

func (s userService) CreateUser(ctx context.Context, userDto *UserDto) error {
	ctx, childSpan := tracing.Tracer.Start(ctx, "CreateUserService", trace.WithAttributes(attribute.String("service", "CreateUser")))
	user := new(models.User)

	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email

	childSpan.End()

	return s.userRepository.CreateUser(ctx, user)
}

func (s userService) UpdateUser(ctx context.Context, id int, userDto *UserDto) error {
	ctx, childSpan := tracing.Tracer.Start(ctx, "UpdateUserService", trace.WithAttributes(attribute.String("service", "UpdateUser")))
	user := new(models.User)

	user.FirstName = userDto.FirstName
	user.LastName = userDto.LastName
	user.Email = userDto.Email

	childSpan.End()

	return s.userRepository.UpdateUser(ctx, id, user)
}

func (s userService) DeleteUser(ctx context.Context, id int) error {
	ctx, childSpan := tracing.Tracer.Start(ctx, "DeleteUserService", trace.WithAttributes(attribute.String("service", "DeleteUser")))
	err := s.userRepository.DeleteUser(ctx, id)
	childSpan.End()

	return err
}
