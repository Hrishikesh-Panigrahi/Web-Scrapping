package controllers

import (
	"fmt"

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

			if product.Name != "" && product.Price != "" && product.Url != "" {
				*products = append(*products, product)
			}
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

func EbayScrapper(encodedKeyword string, products *[]Product) {

	searchURL := fmt.Sprintf("https://www.ebay.com/sch/i.html?_nkw=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains(
			"www.ebay.com",
		),
	)

	collector.OnHTML("li.s-item", func(e *colly.HTMLElement) {

		product := Product{}

		product.Name = e.ChildText(".s-item__title")
		product.Price = e.ChildText(".s-item__price")
		product.Url = fmt.Sprintf("www.ebay.com%s", e.ChildAttr("a.s-item__link", "href"))
		product.Image = e.ChildAttr("img.s-item__image-img", "src")
		product.Source = "Ebay"

		if product.Name != "" && product.Price != "" && product.Url != "" {
			*products = append(*products, product)
		}

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

func WallMartScrapper(encodedKeyword string, products *[]Product) {

	searchURL := fmt.Sprintf("https://www.walmart.com/search/?query=%s", encodedKeyword)

	collector := colly.NewCollector(
		colly.AllowedDomains(
			"www.walmart.com",
		),
	)

	collector.OnHTML("div.search-result-gridview-item", func(e *colly.HTMLElement) {

		product := Product{}

		product.Name = e.ChildText("a.product-title-link span")
		product.Price = e.ChildText("span.price-characteristic")
		product.Url = fmt.Sprintf("https://www.walmart.com%s", e.ChildAttr("a.product-title-link", "href"))
		product.Image = e.ChildAttr("img", "src")
		product.Source = "WallMart"

		if product.Name != "" && product.Price != "" && product.Url != "" {
			*products = append(*products, product)
		}

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
