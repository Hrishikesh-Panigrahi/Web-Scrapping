package main

import (
	"github.com/Hrishikesh-Panigrahi/Web-Scrapping/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// get request for the user url
	r.GET("/", controllers.Index)

	// post request for the user url
	r.POST("/api/web-crawler", controllers.WebCrawler)
	r.GET("/api/web-crawler", controllers.WebCrawler)

	r.Run()
}
