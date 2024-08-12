package main

import (
	"github.com/Hrishikesh-Panigrahi/Web-Scrapping/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	
	r.GET("/", controllers.Index)
	r.POST("/web-crawler", controllers.WebCrawler)
	r.GET("/web-crawler", controllers.ShowResults)

	r.Run()
}
