package core

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"log"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/cache"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/config"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core/database"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/http"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/render"
	html "github.com/JubaerHossain/golang-htmx-starter/static"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type App struct {
	BuildVersion string
	Echo         *echo.Echo
	HttpPort     int
	PublicFS     fs.FS
	cache        cache.CacheService
	DB           *pgxpool.Pool
	MDB          *sql.DB
	logger       *zap.Logger
	Config       *config.Config
}

func StartApp() (*App, error) {
	// Initialize environment variables
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load configuration", zap.Error(err))
		return nil, err
	}

	var pgDB *pgxpool.Pool
	var mySQLDB *sql.DB

	if cfg.DBType == "postgres" {
		db, err := InitPqDatabase(cfg)
		if err != nil {
			return nil, err
		}
		pgDB = db
	} else if cfg.DBType == "mysql" {
		mdb, err := InitMySQLDatabase(cfg)
		if err != nil {
			return nil, err
		}
		mySQLDB = mdb
	}

	cacheService, err := InitCache()
	if err != nil {
		logger.Error("failed to initialize cache", zap.Error(err))
		return nil, err
	}

	app := &App{
		HttpPort:     cfg.AppPort,
		Echo:         echo.New(),
		PublicFS:     html.PublicFS,
		BuildVersion: config.GlobalConfig.AppEnv,
		cache:        cacheService,
		DB:           pgDB,
		MDB:          mySQLDB,
		logger:       logger,
		Config:       cfg,
	}

	// Initialize template renderer
	renderer, err := render.NewRenderer(app.PublicFS)
	if err != nil {
		return nil, errors.New("failed to create template renderer: " + err.Error())
	}

	// Pimp the echo instance
	e := app.Echo
	e.Renderer = renderer
	e.StaticFS("/", html.StaticFS)
	e.Use(http.SetCacheControl)

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} ${method} ${uri} ${status} ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.SessionSecret))))

	return app, nil
}

func InitPqDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	dbService, err := database.NewPgxDatabaseService(cfg)
	if err != nil {
		return nil, err
	}
	return dbService.GetDB(), nil
}

func InitMySQLDatabase(cfg *config.Config) (*sql.DB, error) {
	dbService, err := database.NewMySQLService(cfg)
	if err != nil {
		return nil, err
	}
	return dbService.GetDB(), nil
}

// InitCache initializes the cache
func InitCache() (cache.CacheService, error) {
	ctx := context.Background()
	cacheService, err := cache.NewRedisCacheService(ctx)
	if err != nil {
		return nil, err
	}
	return cacheService, nil
}
