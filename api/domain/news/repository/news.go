package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/news/entity"
)


// NewsRepository defines methods for news data access
type NewsRepository interface {
	GetNewses(r *http.Request) (*entity.NewsResponsePagination, error)
	GetNewsByID(newsID uint) (*entity.News, error)
	GetNews(newsID uint) (*entity.ResponseNews, error)
	CreateNews(news *entity.News, r *http.Request)  error
	UpdateNews(oldNews *entity.News, news *entity.UpdateNews, r *http.Request) error
	DeleteNews(news *entity.News, r *http.Request) error

	GetBreakingScrollingNews(r *http.Request) (*entity.ScrollNewsResponse, error)
	GetBreakingThumbnailNews(r *http.Request) (*entity.ThumbnailNewsResponse, error)
}