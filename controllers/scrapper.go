package controllers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
)

// AmazonScrapper legacy function for backward compatibility
func AmazonScrapper(encodedKeyword string, products *[]Product) {
	ctx := context.Background()
	var mu sync.Mutex
	AmazonScrapperWithContext(ctx, encodedKeyword, products, &mu)
}

// AmazonScrapperWithContext scrapes Amazon with context and mutex for concurrency
func AmazonScrapperWithContext(ctx context.Context, encodedKeyword string, products *[]Product, mu *sync.Mutex) error {
	searchURL := fmt.Sprintf("https://www.amazon.in/s?k=%s", encodedKeyword)

	// Create collector with timeout and debug logging
	collector := colly.NewCollector(
		colly.AllowedDomains("www.amazon.in"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Set timeouts
	collector.SetRequestTimeout(30 * time.Second)

	// User agent to avoid blocking
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var scrapeError error

	collector.OnHTML("div.s-main-slot.s-result-list", func(e *colly.HTMLElement) {
		e.ForEach("div[data-component-type='s-search-result']", func(_ int, h *colly.HTMLElement) {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				scrapeError = ctx.Err()
				return
			default:
			}

			product := Product{}
			product.Name = h.ChildText("span.a-text-normal")
			product.Price = h.ChildText("span.a-price-whole")
			product.Url = fmt.Sprintf("www.amazon.in%s", h.ChildAttr("a.a-link-normal.s-no-outline", "href"))
			product.Image = h.ChildAttr("img.s-image", "src")
			product.Source = "Amazon"

			if product.Name != "" && product.Price != "" && product.Url != "" {
				mu.Lock()
				*products = append(*products, product)
				mu.Unlock()
			}
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("[Amazon] Visiting: %s\n", r.URL.String())
	})

	collector.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("[Amazon] Error: %v\n", err)
		scrapeError = err
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("[Amazon] Response received from: %s (Status: %d)\n", r.Request.URL, r.StatusCode)
	})

	// Set context for the collector
	collector.OnRequest(func(r *colly.Request) {
		// Add context values to the request if needed
		r.Headers.Set("User-Agent", collector.UserAgent)
	})

	if err := collector.Visit(searchURL); err != nil {
		return fmt.Errorf("failed to visit Amazon: %w", err)
	}

	return scrapeError
}

// EbayScrapper legacy function for backward compatibility
func EbayScrapper(encodedKeyword string, products *[]Product) {
	ctx := context.Background()
	var mu sync.Mutex
	EbayScrapperWithContext(ctx, encodedKeyword, products, &mu)
}

// EbayScrapperWithContext scrapes eBay with context and mutex for concurrency
func EbayScrapperWithContext(ctx context.Context, encodedKeyword string, products *[]Product, mu *sync.Mutex) error {
	searchURL := fmt.Sprintf("https://www.ebay.com/sch/i.html?_nkw=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains("www.ebay.com"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Set timeouts
	collector.SetRequestTimeout(30 * time.Second)

	// User agent to avoid blocking
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var scrapeError error

	collector.OnHTML("li.s-item", func(e *colly.HTMLElement) {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			scrapeError = ctx.Err()
			return
		default:
		}

		product := Product{}
		product.Name = e.ChildText(".s-item__title")
		product.Price = e.ChildText(".s-item__price")
		product.Url = e.ChildAttr("a.s-item__link", "href")
		product.Image = e.ChildAttr("img.s-item__image-img", "src")
		product.Source = "Ebay"

		if product.Name != "" && product.Price != "" && product.Url != "" {
			mu.Lock()
			*products = append(*products, product)
			mu.Unlock()
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("[eBay] Visiting: %s\n", r.URL.String())
		// Add context values to the request if needed
		r.Headers.Set("User-Agent", collector.UserAgent)
	})

	collector.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("[eBay] Error: %v\n", err)
		scrapeError = err
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("[eBay] Response received from: %s (Status: %d)\n", r.Request.URL, r.StatusCode)
	})

	if err := collector.Visit(searchURL); err != nil {
		return fmt.Errorf("failed to visit eBay: %w", err)
	}

	return scrapeError
}

// WallMartScrapper legacy function for backward compatibility
func WallMartScrapper(encodedKeyword string, products *[]Product) {
	ctx := context.Background()
	var mu sync.Mutex
	WallMartScrapperWithContext(ctx, encodedKeyword, products, &mu)
}

// WallMartScrapperWithContext scrapes Walmart with context and mutex for concurrency
func WallMartScrapperWithContext(ctx context.Context, encodedKeyword string, products *[]Product, mu *sync.Mutex) error {
	searchURL := fmt.Sprintf("https://www.walmart.com/search/?query=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains("www.walmart.com"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Set timeouts
	collector.SetRequestTimeout(30 * time.Second)

	// User agent to avoid blocking
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var scrapeError error

	collector.OnHTML("div.search-result-gridview-item", func(e *colly.HTMLElement) {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			scrapeError = ctx.Err()
			return
		default:
		}

		product := Product{}
		product.Name = e.ChildText("a.product-title-link span")
		product.Price = e.ChildText("span.price-characteristic")
		product.Url = fmt.Sprintf("https://www.walmart.com%s", e.ChildAttr("a.product-title-link", "href"))
		product.Image = e.ChildAttr("img", "src")
		product.Source = "WallMart"

		if product.Name != "" && product.Price != "" && product.Url != "" {
			mu.Lock()
			*products = append(*products, product)
			mu.Unlock()
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("[Walmart] Visiting: %s\n", r.URL.String())
		// Add context values to the request if needed
		r.Headers.Set("User-Agent", collector.UserAgent)
	})

	collector.OnError(func(_ *colly.Response, err error) {
		fmt.Printf("[Walmart] Error: %v\n", err)
		scrapeError = err
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("[Walmart] Response received from: %s (Status: %d)\n", r.Request.URL, r.StatusCode)
	})

	if err := collector.Visit(searchURL); err != nil {
		return fmt.Errorf("failed to visit Walmart: %w", err)
	}

	return scrapeError
}
