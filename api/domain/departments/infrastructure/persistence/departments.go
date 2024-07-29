package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/departments/entity"
	"github.com/JubaerHossain/cn-api/domain/departments/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
)

type DepartmentRepositoryImpl struct {
	app *app.App
}

// NewDepartmentRepository returns a new instance of DepartmentRepositoryImpl
func NewDepartmentRepository(app *app.App) repository.DepartmentRepository {
	return &DepartmentRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_departments_*"); err != nil {
		return err
	}
	return nil
}

// GetAllDepartments returns all departments from the database
func (r *DepartmentRepositoryImpl) GetDepartments(req *http.Request) (*entity.DepartmentResponsePagination, error) {
	// Implement logic to get all departments
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_departments_%s", req.URL.Query().Encode()) // Encode query parameters
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		departments := &entity.DepartmentResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), departments); err != nil {
			return &entity.DepartmentResponsePagination{}, err
		}
		return departments, nil
	}

	baseQuery := "SELECT id, title, slug, created_by, updated_by, status_id FROM departments"

	// Apply filters from query parameters
	queryValues := req.URL.Query()
	var filters []string

	// Filter by search query
	if search := queryValues.Get("search"); search != "" {
		filters = append(filters, fmt.Sprintf("title LIKE '%%%s%%'", search))
	}

	// Filter by status
	if status := queryValues.Get("status"); status != "" {
		filters = append(filters, fmt.Sprintf("status_id = %s", status))
	}

	// Apply filters to query
	filterQuery := ""
	if len(filters) > 0 {
		filterQuery = " WHERE " + strings.Join(filters, " AND ")
	}

	// Sort by
	sortBy := " ORDER BY id DESC"
	if sort := queryValues.Get("sort"); sort != "" {
		sortBy = fmt.Sprintf(" ORDER BY id %s", sort)
	}

	// Pagination and limits
	pagination, limit, offset, err := utilQuery.Paginate(req, r.app, baseQuery, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("pagination error: %w", err)
	}

	// Apply pagination to query
	query := fmt.Sprintf("%s%s%s LIMIT %d OFFSET %d", baseQuery, filterQuery, sortBy, limit, offset)

	// Perform the query
	rows, err := r.app.MDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and parse the results
	departments := []*entity.ResponseDepartment{}
	for rows.Next() {
		var department entity.ResponseDepartment
		err := rows.Scan(&department.ID, &department.Title, &department.Slug, &department.CreatedBy, &department.UpdatedBy, &department.StatusID)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &department)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := entity.DepartmentResponsePagination{
		Data:       departments,
		Pagination: pagination,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.DepartmentResponsePagination{}, err
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.DepartmentResponsePagination{}, err
	}
	return &response, nil
}

// GetDepartmentByID returns a department by ID from the database
func (r *DepartmentRepositoryImpl) GetDepartmentByID(departmentID uint) (*entity.Department, error) {
	// Implement logic to get department by ID
	department := &entity.Department{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id FROM departments WHERE id = $1", departmentID).Scan(&department.ID); err != nil {
		return nil, fmt.Errorf("department not found")
	}
	return department, nil
}

// GetDepartment returns a department by ID from the database
func (r *DepartmentRepositoryImpl) GetDepartment(departmentID uint) (*entity.ResponseDepartment, error) {
	// Implement logic to get department by ID
	resDepartment := &entity.ResponseDepartment{}
	query := "SELECT id, title, slug, created_by, updated_by, status_id FROM departments WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, departmentID).Scan(&resDepartment.ID, &resDepartment.Title, &resDepartment.Slug, &resDepartment.CreatedBy, &resDepartment.UpdatedBy, &resDepartment.StatusID); err != nil {
		return nil, fmt.Errorf("department not found")
	}
	return resDepartment, nil
}

func (r *DepartmentRepositoryImpl) GetDepartmentDetails(departmentID uint) (*entity.ResponseDepartment, error) {
	// Implement logic to get department details by ID
	resDepartment := &entity.ResponseDepartment{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.title, u.slug, u.created_by, u.updated_by, u.status_id
		FROM departments u
		WHERE u.id = $1
	`, departmentID).Scan(&resDepartment.ID, &resDepartment.Title, &resDepartment.Slug, &resDepartment.CreatedBy, &resDepartment.UpdatedBy, &resDepartment.StatusID)
	if err != nil {
		return nil, fmt.Errorf("department not found")
	}
	return resDepartment, nil
}

func (r *DepartmentRepositoryImpl) CreateDepartment(department *entity.Department, req *http.Request) error {
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

	// Create the department within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO departments (title, slug, created_by, updated_by, status_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, department.Title, department.Slug, department.CreatedBy, department.UpdatedBy, department.StatusID)
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

func (r *DepartmentRepositoryImpl) UpdateDepartment(oldDepartment *entity.Department, department *entity.UpdateDepartment, req *http.Request) error {
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

	// Title     string    `json:"title" validate:"required,min=3,max=100"`
	// Slug      string    `json:"slug" validate:"required,min=3,max=100"`
	// UpdatedBy uint      `json:"updated_by"`
	// StatusID  uint      `json:"status_id"`
	// UpdatedAt time.Time `json:"updated_at"`

	query := `
		UPDATE departments
		SET title = $1, slug = $2, updated_by = $3, status_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	if _, err := tx.Exec(context.Background(), query, department.Title, department.Slug, department.UpdatedBy, department.StatusID, oldDepartment.ID); err != nil {
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

func (r *DepartmentRepositoryImpl) DeleteDepartment(department *entity.Department, req *http.Request) error {
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

	query := "DELETE FROM departments WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, department.ID); err != nil {
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
