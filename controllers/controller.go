package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

func Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func WebCrawler(c *gin.Context) {

	if c.Request.Method == "POST" {

		var products []Product

		url := c.PostForm("url")
		fmt.Println(url)

		collector := colly.NewCollector(
			colly.AllowedDomains("www.scrapingcourse.com"),
		)

		collector.OnHTML("li.product", func(e *colly.HTMLElement) {
			product := Product{}
			product.Url = e.ChildAttr("a", "href")
			product.Image = e.ChildAttr("img", "src")
			product.Name = e.ChildText(".product-name")
			product.Price = e.ChildText(".price")

			products = append(products, product)
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
			log.Printf("Number of products collected: %d", len(products))

			// open the CSV file
			file, err := os.Create("products.csv")
			if err != nil {
				log.Fatalln("Failed to create output CSV file", err)
			}
			defer file.Close()

			// initialize a file writer
			writer := csv.NewWriter(file)

			// write the CSV headers
			headers := []string{
				"Url",
				"Image",
				"Name",
				"Price",
			}
			writer.Write(headers)

			for _, product := range products {
				// convert a Product to an array of strings
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

		collector.Visit("https://www.scrapingcourse.com/ecommerce")
	}

	if c.Request.Method == "GET" {
		type body struct {
			Url, Image, Name, Price string
		}

		var bodyreq []body
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
			singleproduct := body{
				Url:   record[0],
				Image: record[1],
				Name:  record[2],
				Price: record[3],
			}
			bodyreq = append(bodyreq, singleproduct)

		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"bodyreq": bodyreq,
		})
	}

}
