package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/database"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/models"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{db: db}
}

func (r userRepository) GetUserPaginate(ctx context.Context, pagination database.Pagination, search string) (*database.Pagination, error) {
	var (
		_, childSpan = tracing.Tracer.Start(ctx, "GetUserPaginateRepository", trace.WithAttributes(attribute.String("repository", "GetUserPaginate"), attribute.String("search", search)))
		users        []models.User
		err          error
	)

	// Pagination query
	if search != "" {
		if err = r.db.Scopes(database.Paginate(users, &pagination, r.db)).
			Where(`email LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Or(`first_name LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Or(`last_name LIKE ?`, fmt.Sprintf(`%%%s%%`, search)).
			Find(&users).Error; err != nil {
			log.Println(err)
			return nil, errors.New("GetUserPaginateError")
		}
	} else {
		if err = r.db.Scopes(database.Paginate(users, &pagination, r.db)).
			Find(&users).Error; err != nil {
			return nil, err
		}
	}

	// Set data
	pagination.Data = users

	childSpan.End()

	return &pagination, nil
}

func (r userRepository) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var (
		_, childSpan = tracing.Tracer.Start(ctx, "GetUserByIDRepository", trace.WithAttributes(attribute.String("repository", "GetUserByID")))
		user         models.User
		err          error
	)

	// Query
	if err = r.db.First(&user, id).Error; err != nil {
		return user, err
	}

	childSpan.End()

	return user, nil
}

func (r userRepository) CreateUser(ctx context.Context, user *models.User) error {
	var (
		_, childSpan = tracing.Tracer.Start(ctx, "CreateUserRepository", trace.WithAttributes(attribute.String("repository", "CreateUser")))
		err          error
	)

	// Execute
	if err = r.db.Create(&user).Error; err != nil {
		return err
	}

	childSpan.End()

	return nil
}

func (r userRepository) UpdateUser(ctx context.Context, id int, user *models.User) error {
	var (
		_, childSpan = tracing.Tracer.Start(ctx, "UpdateUserRepository", trace.WithAttributes(attribute.String("repository", "UpdateUser")))
		existUser    *models.User
		err          error
	)

	// Get model
	r.db.First(&existUser)

	// Set attributes
	existUser.FirstName = user.FirstName
	existUser.LastName = user.LastName
	existUser.Email = user.Email

	// Execute
	if err = r.db.Save(&existUser).Error; err != nil {
		return err
	}

	childSpan.End()

	return nil
}

func (r userRepository) DeleteUser(ctx context.Context, id int) error {
	var (
		_, childSpan = tracing.Tracer.Start(ctx, "DeleteUserRepository", trace.WithAttributes(attribute.String("repository", "DeleteUser")))
		err          error
	)

	// Execute
	if err = r.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	childSpan.End()

	return nil
}
