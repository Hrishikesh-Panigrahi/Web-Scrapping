package controllers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// saveCSV legacy function for backward compatibility
func saveCSV(products []Product) {
	if err := saveCSVSafe(products); err != nil {
		log.Fatalf("Failed to save CSV: %v", err)
	}
}

// saveCSVSafe saves products to CSV with better error handling
func saveCSVSafe(products []Product) error {
	if len(products) == 0 {
		return fmt.Errorf("no products to save")
	}

	// Create backup of existing file if it exists
	csvPath := "products.csv"
	if _, err := os.Stat(csvPath); err == nil {
		backupPath := fmt.Sprintf("products_backup_%d.csv", time.Now().Unix())
		if err := copyFile(csvPath, backupPath); err != nil {
			log.Printf("Warning: failed to create backup: %v", err)
		} else {
			log.Printf("Created backup: %s", backupPath)
		}
	}

	// Create temporary file first
	tempPath := csvPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary CSV file: %w", err)
	}

	writer := csv.NewWriter(file)

	// Write headers
	headers := []string{
		"Url",
		"Image",
		"Name",
		"Price",
		"Source",
	}
	if err := writer.Write(headers); err != nil {
		file.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write products
	for i, product := range products {
		record := []string{
			product.Url,
			product.Image,
			product.Name,
			product.Price,
			product.Source,
		}
		if err := writer.Write(record); err != nil {
			file.Close()
			os.Remove(tempPath)
			return fmt.Errorf("failed to write product %d: %w", i, err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		file.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	if err := file.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to close CSV file: %w", err)
	}

	// Atomic move: replace the original file
	if err := os.Rename(tempPath, csvPath); err != nil {
		return fmt.Errorf("failed to replace CSV file: %w", err)
	}

	log.Printf("Successfully saved %d products to %s", len(products), csvPath)
	return nil
}

// copyFile creates a backup copy of a file
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = sourceFile.WriteTo(destFile)
	return err
}

// CleanupOldBackups removes backup files older than 7 days
func CleanupOldBackups() {
	pattern := "products_backup_*.csv"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Printf("Error finding backup files: %v", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -7) // 7 days ago

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(match); err != nil {
				log.Printf("Error removing old backup %s: %v", match, err)
			} else {
				log.Printf("Removed old backup: %s", match)
			}
		}
	}
}
