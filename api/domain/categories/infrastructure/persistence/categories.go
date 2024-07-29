package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/categories/entity"
	"github.com/JubaerHossain/cn-api/domain/categories/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
)

type CategoryRepositoryImpl struct {
	app *app.App
}

// NewCategoryRepository returns a new instance of CategoryRepositoryImpl
func NewCategoryRepository(app *app.App) repository.CategoryRepository {
	return &CategoryRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_categorys_*"); err != nil {
		return err
	}
	return nil
}

// getChildCategories retrieves the child categories for a given parent ID.
func (r *CategoryRepositoryImpl) GetCategories(req *http.Request) (*entity.CategoryResponsePagination, error) {
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_categories_%s", req.URL.Query().Encode())
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		categories := &entity.CategoryResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), categories); err != nil {
			return &entity.CategoryResponsePagination{}, fmt.Errorf("cache unmarshal error: %w", err)
		}
		return categories, nil
	}

	baseQuery := "SELECT id, title, slug, `order`, status_id, parent_id FROM news_categories"
	queryValues := req.URL.Query()
	var filters []string

	if search := queryValues.Get("search"); search != "" {
		filters = append(filters, fmt.Sprintf("title LIKE '%%%s%%'", search))
	}

	if status := queryValues.Get("status"); status != "" {
		filters = append(filters, fmt.Sprintf("status_id = %s", status))
	}

	filters = append(filters, "is_featured = 1")
	filters = append(filters, "parent_id IS NULL")

	filterQuery := ""
	if len(filters) > 0 {
		filterQuery = " WHERE " + strings.Join(filters, " AND ")
	}

	sortBy := " ORDER BY `order` ASC"
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

	rows, err := r.app.MDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	categories := []*entity.ResponseCategory{}
	for rows.Next() {
		var category entity.ResponseCategory
		err := rows.Scan(&category.ID, &category.Title, &category.Slug, &category.Order, &category.StatusID, &category.ParentID)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		childCategories, err := r.getChildCategories(category.ID)
		if err != nil {
			return nil, fmt.Errorf("get child categories error: %w", err)
		}
		category.ChildCategory = childCategories
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	response := entity.CategoryResponsePagination{
		Data:       categories,
		Pagination: pagination,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.CategoryResponsePagination{}, fmt.Errorf("response marshal error: %w", err)
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.CategoryResponsePagination{}, fmt.Errorf("cache set error: %w", err)
	}
	return &response, nil
}



// GetCategoryByID returns a category by ID from the database
func (r *CategoryRepositoryImpl) GetCategoryByID(categoryID uint) (*entity.Category, error) {
	// Implement logic to get category by ID
	category := &entity.Category{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, title, status_id FROM news_categories WHERE id = $1", categoryID).Scan(&category.ID, &category.Title, &category.StatusID); err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

// GetCategory returns a category by ID from the database
func (r *CategoryRepositoryImpl) GetCategory(categoryID uint) (*entity.ResponseCategory, error) {
	// Implement logic to get category by ID
	resCategory := &entity.ResponseCategory{}
	query := "SELECT id, title, status_id FROM news_categories WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, categoryID).Scan(&resCategory.ID, &resCategory.Title, &resCategory.StatusID); err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return resCategory, nil
}

func (r *CategoryRepositoryImpl) GetCategoryDetails(categoryID uint) (*entity.ResponseCategory, error) {
	// Implement logic to get category details by ID
	resCategory := &entity.ResponseCategory{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.title, u.status_id
		FROM news_categories u
		WHERE u.id = $1
	`, categoryID).Scan(&resCategory.ID, &resCategory.Title, &resCategory.StatusID)
	if err != nil {
		return nil, fmt.Errorf("category not found")
	}
	return resCategory, nil
}

func (r *CategoryRepositoryImpl) CreateCategory(category *entity.Category, req *http.Request) error {
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

	// Create the category within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO news_categories (title, status_id, created_at, updated_at) VALUES ($1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, category.Title)
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

func (r *CategoryRepositoryImpl) UpdateCategory(oldCategory *entity.Category, category *entity.UpdateCategory, req *http.Request) error {
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
		UPDATE news_categories
		SET title = $1, status_id = $2
		WHERE id = $3
	`
	if _, err := tx.Exec(context.Background(), query, category.Title, category.StatusID, oldCategory.ID); err != nil {
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

func (r *CategoryRepositoryImpl) DeleteCategory(category *entity.Category, req *http.Request) error {
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

	query := "DELETE FROM news_categories WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, category.ID); err != nil {
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
