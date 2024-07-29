package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JubaerHossain/cn-api/domain/news/entity"
	"github.com/JubaerHossain/cn-api/domain/news/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/cache"
	"github.com/JubaerHossain/rootx/pkg/core/config"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
)

type NewsRepositoryImpl struct {
	app *app.App
}

// NewNewsRepository returns a new instance of NewsRepositoryImpl
func NewNewsRepository(app *app.App) repository.NewsRepository {
	return &NewsRepositoryImpl{
		app: app,
	}
}

func CacheClear(req *http.Request, cache cache.CacheService) error {
	ctx := req.Context()
	if _, err := cache.ClearPattern(ctx, "get_all_newss_*"); err != nil {
		return err
	}
	if _, err := cache.ClearPattern(ctx, "get_breaking_scrolling_news*"); err != nil {
		return err
	}
	return nil
}

// GetAllNewss returns all newss from the database
func (r *NewsRepositoryImpl) GetNewses(req *http.Request) (*entity.NewsResponsePagination, error) {
	// Implement logic to get all newss
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_all_newss_%s", req.URL.Query().Encode()) // Encode query parameters
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		newss := &entity.NewsResponsePagination{}
		if err := json.Unmarshal([]byte(cachedData), newss); err != nil {
			return &entity.NewsResponsePagination{}, err
		}
		return newss, nil
	}

	baseQuery := "SELECT id, name, status, created_at FROM news" // Example SQL query

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
	newss := []*entity.ResponseNews{}
	for rows.Next() {
		var news entity.ResponseNews
		err := rows.Scan(&news.ID, &news.Name, &news.Status, &news.CreatedAt)
		if err != nil {
			return nil, err
		}
		newss = append(newss, &news)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := entity.NewsResponsePagination{
		Data:       newss,
		Pagination: pagination,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return &entity.NewsResponsePagination{}, err
	}
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), time.Duration(config.GlobalConfig.RedisExp)*time.Second); err != nil {
		return &entity.NewsResponsePagination{}, err
	}
	return &response, nil
}

// GetNewsByID returns a news by ID from the database
func (r *NewsRepositoryImpl) GetNewsByID(newsID uint) (*entity.News, error) {
	// Implement logic to get news by ID
	news := &entity.News{}
	if err := r.app.DB.QueryRow(context.Background(), "SELECT id, name, status FROM news WHERE id = $1", newsID).Scan(&news.ID, &news.Name, &news.Status); err != nil {
		return nil, fmt.Errorf("news not found")
	}
	return news, nil
}

// GetNews returns a news by ID from the database
func (r *NewsRepositoryImpl) GetNews(newsID uint) (*entity.ResponseNews, error) {
	// Implement logic to get news by ID
	resNews := &entity.ResponseNews{}
	query := "SELECT id, name, status FROM news WHERE id = $1"
	if err := r.app.DB.QueryRow(context.Background(), query, newsID).Scan(&resNews.ID, &resNews.Name, &resNews.Status); err != nil {
		return nil, fmt.Errorf("news not found")
	}
	return resNews, nil
}

func (r *NewsRepositoryImpl) GetNewsDetails(newsID uint) (*entity.ResponseNews, error) {
	// Implement logic to get news details by ID
	resNews := &entity.ResponseNews{}
	err := r.app.DB.QueryRow(context.Background(), `
		SELECT u.id, u.name, u.status
		FROM news u
		WHERE u.id = $1
	`, newsID).Scan(&resNews.ID, &resNews.Name, &resNews.Status)
	if err != nil {
		return nil, fmt.Errorf("news not found")
	}
	return resNews, nil
}

func (r *NewsRepositoryImpl) CreateNews(news *entity.News, req *http.Request) error {
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

	// Create the news within the transaction
	_, err = tx.Exec(context.Background(), `
		INSERT INTO news (name, status, created_at, updated_at) VALUES ($1, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, news.Name)
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

func (r *NewsRepositoryImpl) UpdateNews(oldNews *entity.News, news *entity.UpdateNews, req *http.Request) error {
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
		UPDATE news
		SET name = $1, status = $2
		WHERE id = $3
		RETURNING id, name, status
	`
	row := tx.QueryRow(context.Background(), query, news.Name, news.Status, oldNews.ID)
	updateNews := &entity.News{}
	err = row.Scan(&updateNews.ID, &updateNews.Name, &updateNews.Status)
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

func (r *NewsRepositoryImpl) DeleteNews(news *entity.News, req *http.Request) error {
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

	query := "DELETE FROM news WHERE id = $1"
	if _, err := tx.Exec(context.Background(), query, news.ID); err != nil {
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

// GetBreakingScrollingNews returns breaking and scrolling news
func (r *NewsRepositoryImpl) GetBreakingScrollingNews(req *http.Request) (*entity.ScrollNewsResponse, error) {
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_breaking_scrolling_news_%s", req.URL.Query().Encode()) // Encode query parameters

	// Check cache first
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		news := &entity.ScrollNewsResponse{}
		if err := json.Unmarshal([]byte(cachedData), news); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return news, nil
	}

	var filters []string
	filters = append(filters, "news.status_id = 1")
	filters = append(filters, "news.publish_status_id = 9")
	filters = append(filters, "news.breaking_scroll_news = 1")

	// Apply filters to query
	filterQuery := ""
	if len(filters) > 0 {
		filterQuery = " AND " + strings.Join(filters, " AND ")
	}


	// Get news list
	newsList, err := r.GetNewsList(req, filterQuery, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get news list: %w", err)
	}

	// Construct the response
	response := &entity.ScrollNewsResponse{
		Data: newsList,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	cacheDuration := time.Duration(r.app.Config.RedisExp) * time.Second
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), cacheDuration); err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return response, nil
}

// GetBreakingThumbnailNews returns thumbnail news
func (r *NewsRepositoryImpl) GetBreakingThumbnailNews(req *http.Request) (*entity.ThumbnailNewsResponse, error) {
	ctx := req.Context()
	cacheKey := fmt.Sprintf("get_breaking_thumbnail_news_%s", req.URL.Query().Encode()) // Encode query parameters

	// Check cache first
	if cachedData, errCache := r.app.Cache.Get(ctx, cacheKey); errCache == nil && cachedData != "" {
		news := &entity.ThumbnailNewsResponse{}
		if err := json.Unmarshal([]byte(cachedData), news); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return news, nil
	}

	var filters []string
	filters = append(filters, "news.status_id = 1")
	filters = append(filters, "news.publish_status_id = 9")
	filters = append(filters, "news.breaking_thumb_news = 1")

	// Apply filters to query
	filterQuery := ""
	if len(filters) > 0 {
		filterQuery = " AND " + strings.Join(filters, " AND ")
	}

	// Get news list
	newsList, err := r.GetNewsList(req, filterQuery, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get news list: %w", err)
	}

	// Construct the response
	response := &entity.ThumbnailNewsResponse{
		Data: newsList,
	}

	// Cache the response
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}
	cacheDuration := time.Duration(r.app.Config.RedisExp) * time.Second
	if err := r.app.Cache.Set(ctx, cacheKey, string(jsonData), cacheDuration); err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return response, nil
}
