package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Url, Image, Name, Price, Source string
}

type ErrorResponse struct {
	Error   string    `json:"error"`
	Code    string    `json:"code"`
	Time    time.Time `json:"timestamp"`
	Details string    `json:"details,omitempty"`
}

// Index renders the main search page
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Smart Web Crawler - Find Best Deals",
	})
}

// WebScrapper handles the search request with improved error handling and concurrency
func WebScrapper(c *gin.Context) {
	// Create context with timeout for the entire operation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	var products []Product
	var mu sync.Mutex // Protect products slice
	var wg sync.WaitGroup

	// Get form parameters
	keyword := c.PostForm("keyword")
	amazonbutton := c.PostForm("amazonbutton")
	EbayButton := c.PostForm("EbayButton")
	searchall := c.PostForm("searchall")

	encodedKeyword := url.QueryEscape(keyword)

	// Channel to collect errors from goroutines
	errChan := make(chan error, 3)

	// Execute scrapers concurrently
	if amazonbutton == "amazon" || searchall == "searchall" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := AmazonScrapperWithContext(ctx, encodedKeyword, &products, &mu); err != nil {
				errChan <- fmt.Errorf("Amazon scraper error: %w", err)
			}
		}()
	}

	if EbayButton == "Ebay" || searchall == "searchall" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := EbayScrapperWithContext(ctx, encodedKeyword, &products, &mu); err != nil {
				errChan <- fmt.Errorf("eBay scraper error: %w", err)
			}
		}()
	}

	if searchall == "searchall" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := WallMartScrapperWithContext(ctx, encodedKeyword, &products, &mu); err != nil {
				errChan <- fmt.Errorf("Walmart scraper error: %w", err)
			}
		}()
	}

	// Wait for all scrapers to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Collect any errors
	var errors []string
	for err := range errChan {
		if err != nil {
			log.Printf("Scraper error: %v", err)
			errors = append(errors, err.Error())
		}
	}

	// Check if context was cancelled
	if ctx.Err() != nil {
		c.JSON(http.StatusRequestTimeout, ErrorResponse{
			Error:   "Request timeout",
			Code:    "TIMEOUT",
			Time:    time.Now(),
			Details: "The search operation took too long to complete",
		})
		return
	}

	// Save results even if some scrapers failed
	if len(products) > 0 {
		if err := saveCSVSafe(products); err != nil {
			log.Printf("Error saving CSV: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to save search results",
				Code:    "SAVE_ERROR",
				Time:    time.Now(),
				Details: err.Error(),
			})
			return
		}
	}

	// Log any scraper errors but continue
	if len(errors) > 0 {
		log.Printf("Some scrapers encountered errors: %v", errors)
	}

	// Redirect to results page
	c.Redirect(http.StatusSeeOther, "/web-crawler")
}

// ShowResults displays the search results with error handling
func ShowResults(c *gin.Context) {
	type Productbody struct {
		Url, Image, Name, Price, Source string
	}

	var productbody []Productbody

	// Try to open CSV file
	file, err := os.Open("products.csv")
	if err != nil {
		if os.IsNotExist(err) {
			// No search results yet, show empty state
			c.HTML(http.StatusOK, "web_crawler.html", gin.H{
				"Productbody": productbody,
				"Message":     "No search results found. Try performing a search first.",
			})
			return
		}

		log.Printf("Error opening CSV file: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error":   "Unable to load search results",
			"Details": "There was an error accessing the search results file",
		})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error":   "Unable to read search results",
			"Details": "The search results file appears to be corrupted",
		})
		return
	}

	// Process records (skip header)
	for i, record := range records {
		if i == 0 || len(record) < 5 {
			continue
		}

		singleproduct := Productbody{
			Url:    record[0],
			Image:  record[1],
			Name:   record[2],
			Price:  record[3],
			Source: record[4],
		}
		productbody = append(productbody, singleproduct)
	}

	c.HTML(http.StatusOK, "web_crawler.html", gin.H{
		"Productbody": productbody,
		"Title":       "Search Results",
	})
}
