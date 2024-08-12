package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

// Index is a controller to render the index.html file
func Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

// WebCrawler is a controller to crawl the web
// and it saves the data in products.csv file
// and redirects to /web-crawler route
func WebCrawler(c *gin.Context) {
	if c.Request.Method == "POST" {

		var products []Product

		keyword := c.PostForm("keyword")
		encodedKeyword := url.QueryEscape(keyword)

		searchURL := fmt.Sprintf("https://www.amazon.in/s?k=%s", encodedKeyword)

		collector := colly.NewCollector(
			colly.AllowedDomains("www.scrapingcourse.com", "www.amazon.in", "www.flipkart.com"),
		)

		collector.OnHTML("div.s-result-list.s-search-results.sg-row", func(e *colly.HTMLElement) {
			e.ForEach("div.a-section.a-spacing-base", func(_ int, h *colly.HTMLElement) {
				product := Product{}

				product.Name = h.ChildText("span.a-size-base-plus.a-color-base.a-text-normal")
				product.Price = h.ChildText("span.a-price-whole")
				product.Url = h.ChildAttr("a.a-link-normal.s-no-outline", "href")
				product.Image = h.ChildAttr("img.s-image", "src")

				products = append(products, product)
			})
		})

		collector.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		collector.OnError(func(_ *colly.Response, err error) {
			fmt.Println("Something went wrong: ", err)
		})

		collector.OnResponse(func(r *colly.Response) {
			fmt.Println("Page visited: ", r.Request.URL)
		})

		collector.OnScraped(func(r *colly.Response) {
			file, err := os.Create("products.csv")
			if err != nil {
				log.Fatalln("Failed to create output CSV file", err)
			}
			defer file.Close()

			writer := csv.NewWriter(file)

			headers := []string{
				"Url",
				"Image",
				"Name",
				"Price",
			}
			writer.Write(headers)

			for _, product := range products {
				record := []string{
					product.Url,
					product.Image,
					product.Name,
					product.Price,
				}
				writer.Write(record)
			}
			defer writer.Flush()
			for _, list := range products {
				fmt.Println(list)
			}
		})

		collector.Visit(searchURL)

		c.Redirect(http.StatusMovedPermanently, "/web-crawler")
	}

}

// ShowResults is a controller to render the web_crawler.html file
// and it reads the data from products.csv file
// and renders the data in the web_crawler.html file
func ShowResults(c *gin.Context) {
	if c.Request.Method == "GET" {
		type Productbody struct {
			Url, Image, Name, Price string
		}

		var productbody []Productbody
		file, err := os.Open("products.csv")
		if err != nil {
			log.Fatalln("Failed to open CSV file", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)

		records, err := reader.ReadAll()
		if err != nil {
			log.Fatalln("Failed to read CSV file", err)
		}

		for i, record := range records {
			if i == 0 {
				continue
			}
			singleproduct := Productbody{
				Url:   record[0],
				Image: record[1],
				Name:  record[2],
				Price: record[3],
			}
			productbody = append(productbody, singleproduct)
		}
		c.HTML(http.StatusOK, "web_crawler.html", gin.H{
			"Productbody": productbody,
		})
	}
}
