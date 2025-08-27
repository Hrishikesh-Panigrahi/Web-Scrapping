package middleware

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Rate limiter store
var limiters = make(map[string]*rate.Limiter)

// RateLimitMiddleware implements rate limiting per IP
func RateLimitMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ip := c.ClientIP()

		// Get or create limiter for this IP
		limiter, exists := limiters[ip]
		if !exists {
			// Allow 5 requests per minute per IP
			limiter = rate.NewLimiter(rate.Every(time.Minute/5), 1)
			limiters[ip] = limiter
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded. Please try again later.",
				"retry_after": "60 seconds",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// InputValidationMiddleware validates and sanitizes input
func InputValidationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == "POST" {
			keyword := c.PostForm("keyword")

			// Validate keyword
			if strings.TrimSpace(keyword) == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Keyword is required",
				})
				c.Abort()
				return
			}

			// Sanitize keyword (remove potentially harmful characters)
			keyword = sanitizeInput(keyword)
			if len(keyword) > 100 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Keyword too long (max 100 characters)",
				})
				c.Abort()
				return
			}

			// Check if at least one platform is selected
			amazonbutton := c.PostForm("amazonbutton")
			ebayButton := c.PostForm("EbayButton")
			searchall := c.PostForm("searchall")

			if amazonbutton == "" && ebayButton == "" && searchall == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Please select at least one platform to search",
				})
				c.Abort()
				return
			}

			// Store sanitized keyword back in form
			c.Request.PostForm.Set("keyword", keyword)
		}

		c.Next()
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.tailwindcss.com https://unpkg.com https://fonts.googleapis.com https://fonts.gstatic.com; img-src 'self' data: https:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.tailwindcss.com https://unpkg.com;")

		c.Next()
	})
}

// RequestLoggingMiddleware logs requests with additional security info
func RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s | %d | %v | %s | %s",
			method,
			path,
			ip,
			statusCode,
			latency,
			userAgent,
			c.Errors.String(),
		)
	})
}

// CORSMiddleware configures CORS
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://127.0.0.1:8080"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	config.AllowCredentials = true

	return cors.New(config)
}

// sanitizeInput removes potentially harmful characters from input
func sanitizeInput(input string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	input = re.ReplaceAllString(input, "")

	// Remove script tags content
	re = regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = re.ReplaceAllString(input, "")

	// Remove dangerous characters
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "{", "}", "[", "]", "\\", "\"", "'"}
	for _, char := range dangerousChars {
		input = strings.ReplaceAll(input, char, "")
	}

	return strings.TrimSpace(input)
}
