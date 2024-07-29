package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/designations/entity"
	"github.com/JubaerHossain/cn-api/domain/designations/repository"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
)

type DesignationRepositoryImpl struct {
	app *app.App
}

// NewDesignationRepository returns a new instance of DesignationRepositoryImpl
func NewDesignationRepository(app *app.App) repository.DesignationRepository {
	return &DesignationRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_designations_*"); err != nil {
		return err
	}
	return nil
}

// GetAllDesignations returns all designations from the database
func (r *DesignationRepositoryImpl) GetDesignations(req *http.Request) (*entity.DesignationResponsePagination, error) {
	// Implement logic to get all designations
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_designations_%s", req.URL.Query().Encode()) // Encode query parameters
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		designations := &entity.DesignationResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), designations); err != nil {
			return &entity.DesignationResponsePagination{}, err
		}
		return designations, nil
	}

	
	baseQuery := "SELECT id, name, status, created_at FROM designations" // Example SQL query

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
	designations := []*entity.ResponseDesignation{}
	for rows.Next() {
		var designation entity.ResponseDesignation
		err := rows.Scan(&designation.ID, &designation.Name, &designation.Status, &designation.CreatedAt)
		if err != nil {
			return nil, err
		}
		designations = append(designations, &designation)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := entity.DesignationResponsePagination{
		Data: designations,
		Pagination: pagination,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.DesignationResponsePagination{}, err
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.DesignationResponsePagination{}, err
	}
	return &response, nil
}


// GetDesignationByID returns a designation by ID from the database
func (r *DesignationRepositoryImpl) GetDesignationByID(designationID uint) (*entity.Designation, error) {
	// Implement logic to get designation by ID
	designation := &entity.Designation{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, name, status FROM designations WHERE id = $1", designationID).Scan(&designation.ID, &designation.Name, &designation.Status ); err != nil {
		return nil, fmt.Errorf("designation not found")
	}
	return designation, nil
}

// GetDesignation returns a designation by ID from the database
func (r *DesignationRepositoryImpl) GetDesignation(designationID uint) (*entity.ResponseDesignation, error) {
	// Implement logic to get designation by ID
	resDesignation := &entity.ResponseDesignation{}
	query := "SELECT id, name, status FROM designations WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, designationID).Scan(&resDesignation.ID, &resDesignation.Name, &resDesignation.Status); err != nil {
		return nil, fmt.Errorf("designation not found")
	}
	return resDesignation, nil
}

func (r *DesignationRepositoryImpl) GetDesignationDetails(designationID uint) (*entity.ResponseDesignation, error) {
	// Implement logic to get designation details by ID
	resDesignation := &entity.ResponseDesignation{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.name, u.status
		FROM designations u
		WHERE u.id = $1
	`, designationID).Scan(&resDesignation.ID, &resDesignation.Name, &resDesignation.Status)
	if err != nil {
		return nil, fmt.Errorf("designation not found")
	}
	return resDesignation, nil
}

func (r *DesignationRepositoryImpl) CreateDesignation(designation *entity.Designation, req *http.Request) error {
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

	// Create the designation within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO designations (name, status, created_at, updated_at) VALUES ($1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, designation.Name)
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

func (r *DesignationRepositoryImpl) UpdateDesignation(oldDesignation *entity.Designation, designation *entity.UpdateDesignation, req *http.Request)  error {
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
		UPDATE designations
		SET name = $1, status = $2
		WHERE id = $3
		RETURNING id, name, status
	`
	row := tx.QueryRow(context.Background(), query, designation.Name,  designation.Status, oldDesignation.ID)
	updateDesignation := &entity.Designation{}
	err = row.Scan(&updateDesignation.ID, &updateDesignation.Name, &updateDesignation.Status)
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

func (r *DesignationRepositoryImpl) DeleteDesignation(designation *entity.Designation, req *http.Request) error {
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

	query := "DELETE FROM designations WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, designation.ID); err != nil {
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

