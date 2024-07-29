package persistence

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/news/entity"
)

func (r *NewsRepositoryImpl) GetNewsList(req *http.Request, where string, limit uint) ([]*entity.ScrollNews, error) {
	ctx := req.Context()

	// Base SQL query
	baseQuery := `
	SELECT DISTINCT 
	    news.id, 
	    news_translations.title, 
	    news_translations.slug, 
	    news.type, 
	    news_translations.sub_title, 
	    news_translations.tags, 
	    news_translations.content, 
	    news_translations.meta_title, 
	    news_translations.meta_description, 
	    news_translations.meta_keywords, 
	    news.updated_at, 
	    users.name as author, 
	    news.path_small, 
	    news.path_medium, 
	    news.path_large, 
	    news.status_id, 
	    news_categories.title
	FROM news
	JOIN news_translations ON news.id = news_translations.news_id
	JOIN assign_categories ON news.id = assign_categories.news_id
	JOIN news_categories ON assign_categories.news_category_id = news_categories.id
	JOIN users ON news.created_by = users.id
	WHERE news_translations.locale = 'en' %s
	ORDER BY news.id DESC
	LIMIT ?;
	`

	// Combine base query with where clause
	query := fmt.Sprintf(baseQuery, where)

	rows, err := r.app.MDB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Parse the results
	var newsList []*entity.ScrollNews
	for rows.Next() {
		var (
			news   entity.ScrollNews
			tags   string
			meta_keywords string
		)
		if err := rows.Scan(
			&news.ID,
			&news.Title,
			&news.Slug,
			&news.Type,
			&news.SubTitle,
			&tags,
			&news.Content,
			&news.MetaTitle,
			&news.MetaDesc,
			&meta_keywords,
			&news.UpdatedAt,
			&news.Author,
			&news.PathSmall,
			&news.PathMedium,
			&news.PathLarge,
			&news.Status,
			&news.Category,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if err := json.Unmarshal([]byte(tags), &news.Tags); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
		if err := json.Unmarshal([]byte(meta_keywords), &news.MetaKeywords); err != nil {
			return nil, fmt.Errorf("failed to unmarshal meta keywords: %w", err)
		}
		news.URL = fmt.Sprintf("%s/news/%s", r.app.Config.Domain, news.Slug)
		news.Loading = fmt.Sprintf("%s/uploads/%s", r.app.Config.Domain, "default/loading.png")
		newsList = append(newsList, &news)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return newsList, nil
}
