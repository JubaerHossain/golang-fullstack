package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/news/entity"
	"github.com/JubaerHossain/cn-api/domain/news/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/news/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"go.uber.org/zap"
)

type Service struct {
	app  *app.App
	repo repository.NewsRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewNewsRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetNewses(r *http.Request) (*entity.NewsResponsePagination, error) {
	// Call repository to get all news
	news, newsErr := s.repo.GetNewses(r)
	if newsErr != nil {
		s.app.Logger.Error("Error getting news", zap.Error(newsErr))
		return nil, newsErr
	}
	return news, nil
}



// CreateNews creates a new news
func (s *Service) CreateNews(news *entity.News, r *http.Request)  error {
	// Add any validation or business logic here before creating the news
    if err := s.repo.CreateNews(news, r); err != nil {
		s.app.Logger.Error("Error creating news", zap.Error(err))
        return err
    }
	return nil
}

func (s *Service) GetNewsByID(r *http.Request) (*entity.News, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid news ID")
	}
	news, newsErr := s.repo.GetNewsByID(uint(id))
	if newsErr != nil {
		s.app.Logger.Error("Error getting news by ID", zap.Error(newsErr))
		return nil, newsErr
	}
	return news, nil
}

// GetNewsDetails retrieves a news by ID
func (s *Service) GetNewsDetails(r *http.Request) (*entity.ResponseNews, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid news ID")
	}
	news, newsErr := s.repo.GetNews(uint(id))
	if newsErr != nil {
		s.app.Logger.Error("Error getting news details", zap.Error(newsErr))
		return nil, newsErr
	}
	return news, nil
}

// UpdateNews updates an existing news
func (s *Service) UpdateNews(r *http.Request, news *entity.UpdateNews)  error {
	// Call repository to update news
	oldNews, err := s.GetNewsByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateNews(oldNews, news, r)
	if err2 != nil {
		s.app.Logger.Error("Error updating news", zap.Error(err2))
		return err2
	}
	return  nil
}

// DeleteNews deletes a news by ID
func (s *Service) DeleteNews(r *http.Request) error {
	// Call repository to delete news
	news, err := s.GetNewsByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteNews(news, r)
	if err2 != nil {
		s.app.Logger.Error("Error deleting news", zap.Error(err2))
		return err2
	}

	return nil
}

// GetBreakingScrollingNews retrieves breaking and scrolling news
func (s *Service) GetBreakingScrollingNews(r *http.Request) (*entity.ScrollNewsResponse, error) {
	// Call repository to get breaking and scrolling news
	news, newsErr := s.repo.GetBreakingScrollingNews(r)
	if newsErr != nil {
		s.app.Logger.Error("Error getting breaking and scrolling news", zap.Error(newsErr))
		return nil, newsErr
	}
	return news, nil
}

// GetThumbnailNews retrieves thumbnail news
func (s *Service) GetBreakingThumbnailNews(r *http.Request) (*entity.ThumbnailNewsResponse, error) {
	// Call repository to get thumbnail news
	news, newsErr := s.repo.GetBreakingThumbnailNews(r)
	if newsErr != nil {
		s.app.Logger.Error("Error getting thumbnail news", zap.Error(newsErr))
		return nil, newsErr
	}
	return news, nil
}
