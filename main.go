package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hrishikesh-Panigrahi/Web-Scrapping/config"
	"github.com/Hrishikesh-Panigrahi/Web-Scrapping/controllers"
	"github.com/Hrishikesh-Panigrahi/Web-Scrapping/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.New()

	// Add middleware
	r.Use(middleware.RequestLoggingMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	// Rate limiting disabled for development
	// r.Use(middleware.RateLimitMiddleware())

	// Load HTML templates
	r.LoadHTMLGlob("ui/templates/*")

	// Static assets (JS, CSS)
	r.Static("/static", "./ui/static")

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// Landing page
	r.GET("/", controllers.Landing)
	r.POST("/download", controllers.Download)

	// App dashboard (Instagram Reels + YouTube MP4 download)
	r.GET("/app", controllers.Dashboard)
	r.POST("/crawl/instagram", controllers.InstagramCrawl)
	r.POST("/download/youtube", controllers.YouTubeDownload)

	// Legacy product search
	r.GET("/search", controllers.Index)

	// Apply input validation for POST routes
	searchGroup := r.Group("/")
	searchGroup.Use(middleware.InputValidationMiddleware())
	{
		searchGroup.POST("/web-crawler", controllers.WebScrapper)
	}

	r.GET("/web-crawler", controllers.ShowResults)

	// Create server with timeouts
	srv := &http.Server{
		Addr:           ":" + cfg.ServerPort,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.RequestTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.RequestTimeout) * time.Second,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s in %s mode", cfg.ServerPort, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
