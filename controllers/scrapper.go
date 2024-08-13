package controllers

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func AmazonScrapper(encodedKeyword string, products *[]Product) {

	searchURL := fmt.Sprintf("https://www.amazon.in/s?k=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains(
			"www.amazon.in",
		),
	)

	collector.OnHTML("div.s-main-slot.s-result-list", func(e *colly.HTMLElement) {
		e.ForEach("div[data-component-type='s-search-result']", func(_ int, h *colly.HTMLElement) {
			product := Product{}

			product.Name = h.ChildText("span.a-text-normal")
			product.Price = h.ChildText("span.a-price-whole")
			product.Url = fmt.Sprintf("www.amazon.in%s", h.ChildAttr("a.a-link-normal.s-no-outline", "href"))
			product.Image = h.ChildAttr("img.s-image", "src")
			product.Source = "Amazon"

			*products = append(*products, product)
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

	collector.Visit(searchURL)

}

func FlipkartScrapper(encodedKeyword string, products *[]Product) {

	searchURL := fmt.Sprintf("https://www.flipkart.com/search?q=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains(
			"www.flipkart.com",
		),
	)

	// Set the timeout for the collector
	collector.SetRequestTimeout(30 * time.Second)

	collector.OnHTML("div._1YokD2._3Mn1Gg", func(e *colly.HTMLElement) {
		e.ForEach("div._1AtVbE", func(_ int, h *colly.HTMLElement) {
			product := Product{}

			product.Name = h.ChildText("a.IRpwTa")
			product.Price = h.ChildText("div._30jeq3")
			product.Url = fmt.Sprintf("www.flipkart.com%s", h.ChildAttr("a.IRpwTa", "href"))
			product.Image = h.ChildAttr("img._396cs4", "src")
			product.Source = "Flipkart"

			*products = append(*products, product)
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	})

	collector.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	collector.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	collector.Visit(searchURL)

}
