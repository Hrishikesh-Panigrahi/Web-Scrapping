package controllers

import (
	"encoding/csv"
	"log"
	"os"
)

func saveCSV(products []Product) {
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
}
