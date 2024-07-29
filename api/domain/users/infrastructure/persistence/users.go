package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/users/entity"
	"github.com/JubaerHossain/cn-api/domain/users/repository"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
)

type UserRepositoryImpl struct {
	app *app.App
}

// NewUserRepository returns a new instance of UserRepositoryImpl
func NewUserRepository(app *app.App) repository.UserRepository {
	return &UserRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_users_*"); err != nil {
		return err
	}
	return nil
}

// GetAllUsers returns all users from the database
func (r *UserRepositoryImpl) GetUsers(req *http.Request) (*entity.UserResponsePagination, error) {
	// Implement logic to get all users
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_users_%s", req.URL.Query().Encode()) // Encode query parameters
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		users := &entity.UserResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), users); err != nil {
			return &entity.UserResponsePagination{}, err
		}
		return users, nil
	}

	
	baseQuery := "SELECT id, name, status, created_at FROM users" // Example SQL query

	// Apply filters from query parameters
	queryValues := req.URL.Query()
	var filters []string

	// Filter by search query
	if search := queryValues.Get("search"); search != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE '%%%s%%'", search))
	}

	// Filter by status
	if status := queryValues.Get("status"); status != "" {
		filters = append(filters, fmt.Sprintf("status = %s", status))
	}

	// Apply filters to query
	filterQuery := ""
	if len(filters) > 0 {
		filterQuery = " WHERE " + strings.Join(filters, " AND ")
	}

	// sort by
	sortBy := " ORDER BY bf.id DESC"
	if sort := queryValues.Get("sort"); sort != "" {
		sortBy = fmt.Sprintf(" ORDER BY bf.id %s", sort)
	}

	// Pagination and limits
	pagination, limit, offset, err := utilQuery.Paginate(req, r.app, baseQuery, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("pagination error: %w", err)
	}

	// Apply pagination to query
	query := fmt.Sprintf("%s%s%s LIMIT %d OFFSET %d", baseQuery, filterQuery, sortBy, limit, offset)

	// Get database connection from pool
	conn, err := r.app.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Perform the query
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and parse the results
	users := []*entity.ResponseUser{}
	for rows.Next() {
		var user entity.ResponseUser
		err := rows.Scan(&user.ID, &user.Name, &user.Status, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := entity.UserResponsePagination{
		Data: users,
		Pagination: pagination,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.UserResponsePagination{}, err
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.UserResponsePagination{}, err
	}
	return &response, nil
}


// GetUserByID returns a user by ID from the database
func (r *UserRepositoryImpl) GetUserByID(userID uint) (*entity.User, error) {
	// Implement logic to get user by ID
	user := &entity.User{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, name, status FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Status ); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUser returns a user by ID from the database
func (r *UserRepositoryImpl) GetUser(userID uint) (*entity.ResponseUser, error) {
	// Implement logic to get user by ID
	resUser := &entity.ResponseUser{}
	query := "SELECT id, name, status FROM users WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, userID).Scan(&resUser.ID, &resUser.Name, &resUser.Status); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return resUser, nil
}

func (r *UserRepositoryImpl) GetUserDetails(userID uint) (*entity.ResponseUser, error) {
	// Implement logic to get user details by ID
	resUser := &entity.ResponseUser{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.name, u.status
		FROM users u
		WHERE u.id = $1
	`, userID).Scan(&resUser.ID, &resUser.Name, &resUser.Status)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return resUser, nil
}

func (r *UserRepositoryImpl) CreateUser(user *entity.User, req *http.Request) error {
	// Begin a transaction
	tx, err := r.app.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic and rollback the transaction
			tx.Rollback(context.Background())
		} else if err := tx.Commit(context.Background()); err != nil {
			// Commit the transaction if no error occurred, otherwise rollback
			tx.Rollback(context.Background())
		}
	}()

	// Create the user within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO users (name, status, created_at, updated_at) VALUES ($1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, user.Name)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Clear cache
	if err := CacheClear(req, r.app.Cache); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) UpdateUser(oldUser *entity.User, user *entity.UpdateUser, req *http.Request)  error {
	tx, err := r.app.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		} else if err := tx.Commit(context.Background()); err != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := `
		UPDATE users
		SET name = $1, status = $2
		WHERE id = $3
		RETURNING id, name, status
	`
	row := tx.QueryRow(context.Background(), query, user.Name,  user.Status, oldUser.ID)
	updateUser := &entity.User{}
	err = row.Scan(&updateUser.ID, &updateUser.Name, &updateUser.Status)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Clear cache
	if err := CacheClear(req, r.app.Cache); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) DeleteUser(user *entity.User, req *http.Request) error {
	tx, err := r.app.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(context.Background())
		} else if err := tx.Commit(context.Background()); err != nil {
			tx.Rollback(context.Background())
		}
	}()

	query := "DELETE FROM users WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, user.ID); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Clear cache
	if err := CacheClear(req, r.app.Cache); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}

