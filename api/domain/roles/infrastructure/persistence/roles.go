package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/roles/entity"
	"github.com/JubaerHossain/cn-api/domain/roles/repository"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
)

type RoleRepositoryImpl struct {
	app *app.App
}

// NewRoleRepository returns a new instance of RoleRepositoryImpl
func NewRoleRepository(app *app.App) repository.RoleRepository {
	return &RoleRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_roles_*"); err != nil {
		return err
	}
	return nil
}

// GetAllRoles returns all roles from the database
func (r *RoleRepositoryImpl) GetRoles(req *http.Request) (*entity.RoleResponsePagination, error) {
	// Implement logic to get all roles
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_roles_%s", req.URL.Query().Encode()) // Encode query parameters
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		roles := &entity.RoleResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), roles); err != nil {
			return &entity.RoleResponsePagination{}, err
		}
		return roles, nil
	}

	
	baseQuery := "SELECT id, name, status, created_at FROM roles" // Example SQL query

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
	roles := []*entity.ResponseRole{}
	for rows.Next() {
		var role entity.ResponseRole
		err := rows.Scan(&role.ID, &role.Name, &role.Status, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := entity.RoleResponsePagination{
		Data: roles,
		Pagination: pagination,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.RoleResponsePagination{}, err
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.RoleResponsePagination{}, err
	}
	return &response, nil
}


// GetRoleByID returns a role by ID from the database
func (r *RoleRepositoryImpl) GetRoleByID(roleID uint) (*entity.Role, error) {
	// Implement logic to get role by ID
	role := &entity.Role{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, name, status FROM roles WHERE id = $1", roleID).Scan(&role.ID, &role.Name, &role.Status ); err != nil {
		return nil, fmt.Errorf("role not found")
	}
	return role, nil
}

// GetRole returns a role by ID from the database
func (r *RoleRepositoryImpl) GetRole(roleID uint) (*entity.ResponseRole, error) {
	// Implement logic to get role by ID
	resRole := &entity.ResponseRole{}
	query := "SELECT id, name, status FROM roles WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, roleID).Scan(&resRole.ID, &resRole.Name, &resRole.Status); err != nil {
		return nil, fmt.Errorf("role not found")
	}
	return resRole, nil
}

func (r *RoleRepositoryImpl) GetRoleDetails(roleID uint) (*entity.ResponseRole, error) {
	// Implement logic to get role details by ID
	resRole := &entity.ResponseRole{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.name, u.status
		FROM roles u
		WHERE u.id = $1
	`, roleID).Scan(&resRole.ID, &resRole.Name, &resRole.Status)
	if err != nil {
		return nil, fmt.Errorf("role not found")
	}
	return resRole, nil
}

func (r *RoleRepositoryImpl) CreateRole(role *entity.Role, req *http.Request) error {
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

	// Create the role within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO roles (name, status, created_at, updated_at) VALUES ($1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, role.Name)
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

func (r *RoleRepositoryImpl) UpdateRole(oldRole *entity.Role, role *entity.UpdateRole, req *http.Request)  error {
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
		UPDATE roles
		SET name = $1, status = $2
		WHERE id = $3
		RETURNING id, name, status
	`
	row := tx.QueryRow(context.Background(), query, role.Name,  role.Status, oldRole.ID)
	updateRole := &entity.Role{}
	err = row.Scan(&updateRole.ID, &updateRole.Name, &updateRole.Status)
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

func (r *RoleRepositoryImpl) DeleteRole(role *entity.Role, req *http.Request) error {
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

	query := "DELETE FROM roles WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, role.ID); err != nil {
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

