package controllers

import (
	"encoding/csv"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Url, Image, Name, Price, Source string
}

// Index is a controller to render the index.html file
func Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

// WebScrapper is a controller to crawl the web
// and it saves the data in products.csv file
// and redirects to /web-crawler route
func WebScrapper(c *gin.Context) {
	if c.Request.Method == "POST" {

		var products []Product

		keyword := c.PostForm("keyword")
		amazonbutton := c.PostForm("amazonbutton")
		flipkartbutton := c.PostForm("flipkartbutton")
		searchall := c.PostForm("searchall")

		encodedKeyword := url.QueryEscape(keyword)

		if amazonbutton == "amazon" {
			AmazonScrapper(encodedKeyword, &products)
		} else if flipkartbutton == "flipkart" {
			FlipkartScrapper(encodedKeyword, &products)
		} else if searchall == "searchall" {
			AmazonScrapper(encodedKeyword, &products)
			FlipkartScrapper(encodedKeyword, &products)
		}

		saveCSV(products)

		c.Redirect(http.StatusMovedPermanently, "/web-crawler")
	}

}

// ShowResults is a controller to render the web_crawler.html file
// and it reads the data from products.csv file
// and renders the data in the web_crawler.html file
func ShowResults(c *gin.Context) {
	if c.Request.Method == "GET" {
		type Productbody struct {
			Url, Image, Name, Price, Source string
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
		})
	}
}
