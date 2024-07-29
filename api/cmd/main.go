package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "github.com/JubaerHossain/cn-api/docs"
	"github.com/JubaerHossain/cn-api/pkg/api"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/health"
	"github.com/JubaerHossain/rootx/pkg/core/middleware"
	"github.com/JubaerHossain/rootx/pkg/core/monitor"
	"github.com/JubaerHossain/rootx/pkg/utils"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title           News API
// @version         1.0
// @description     This is the API documentation for the News API
// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   jubaer01.cse@gmail.com
// @host            localhost:9009
// @BasePath        /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Initialize the application
	application, err := app.StartApp()
	if err != nil {
		log.Fatalf("❌ Failed to start application: %v", err)
	}

	// Initialize HTTP server
	httpServer := initHTTPServer(application)

	go func() {
		if err := startHTTPServer(application, httpServer); err != nil {
			log.Printf("❌ %v", err)
			log.Println("🔄 Trying to start the server on another port...")
			if err := startHTTPServerOnAvailablePort(application, httpServer); err != nil {
				log.Fatalf("❌ Failed to start server on another port: %v", err)
			}
		}
	}()

	baseURL := fmt.Sprintf("http://localhost:%d", application.Config.AppPort)
	log.Printf("🌐 API base URL: %s", baseURL)

	// Open Swagger URL in browser if in development environment
	if application.Config.AppEnv == "development" {
		openBrowser(baseURL)
	}

	// Graceful shutdown
	gracefulShutdown(httpServer, 5*time.Second)
}

func initHTTPServer(application *app.App) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", application.Config.AppPort),
		Handler: setupRoutes(application),
	}
}

func startHTTPServer(application *app.App, server *http.Server) error {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		application.Logger.Error("Could not start server: %v", zap.Error(err))
		return fmt.Errorf("could not start server: %v", err)
	}
	return nil
}

func startHTTPServerOnAvailablePort(application *app.App, server *http.Server) error {
	for i := application.Config.AppPort + 1; i <= application.Config.AppPort+10; i++ {
		newAddr := fmt.Sprintf(":%d", i)
		server.Addr = newAddr
		log.Printf("Trying to start server on port %d...", i)
		err := startHTTPServer(application, server)
		if err == nil {
			log.Printf("✅ Server started on port %d", i)
			return nil
		}
	}
	return errors.New("could not find available port to start server")
}

func setupRoutes(application *app.App) http.Handler {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register health check endpoint
	mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(health.HealthCheckHandler())))

	// Register monitoring endpoint
	mux.Handle("/metrics", monitor.MetricsHandler())

	// Register Swagger routes
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
		httpSwagger.BeforeScript(`const UrlMutatorPlugin = (system) => ({
			rootInjects: {
				setScheme: (scheme) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const schemes = Array.isArray(scheme) ? scheme : [scheme];
				const newJsonSpec = Object.assign({}, jsonSpec, { schemes });

				return system.specActions.updateJsonSpec(newJsonSpec);
				},
				setHost: (host) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { host });

				return system.specActions.updateJsonSpec(newJsonSpec);
				},
				setBasePath: (basePath) => {
				const jsonSpec = system.getState().toJSON().spec.json;
				const newJsonSpec = Object.assign({}, jsonSpec, { basePath });

				return system.specActions.updateJsonSpec(newJsonSpec);
				}
			}
		});`),
		httpSwagger.Plugins([]string{"UrlMutatorPlugin"}),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": fmt.Sprintf(`() => {
			window.ui.setScheme('%s');
			window.ui.setHost('%s');
			window.ui.setBasePath('%s');
		}`, "http", "localhost:9008", "/api/v1"),
		}),
	))

	// Register API routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api.APIRouter(application)))
	mux.Handle("/api/public/v1/", http.StripPrefix("/api/public/v1", api.PublicAPIRouter(application)))

	// Register file uploads route
	mux.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir("storage"))))

	// Add security headers
	mux.Handle("/", middleware.LimiterMiddleware(middleware.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "Welcome to the API"})
	}))))

	return middleware.PrometheusMiddleware(mux, monitor.RequestsTotal(), monitor.RequestDuration())
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	case "darwin":
		cmd = "open"
	default:
		return
	}
	args = append(args, url)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("⚙️ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Could not gracefully shutdown the server: %v", err)
	}

	log.Printf("✅ Server gracefully stopped")
}
