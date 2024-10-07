package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Open the CSV file
	file, err := os.Open("new_cost.csv")
	if err != nil {
		log.Fatal("Unable to open CSV file", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Unable to read CSV file", err)
	}

	// Set up the database connection using environment variables
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal("Unable to connect to the database", err)
	}
	defer db.Close()

	// Loop through the CSV records starting from row 1 (assuming row 0 is the header)
	for i, record := range records[1:] {
		// Get the cost (Column D) and SKU (Column E)

		// Remove commas from cost string and then parse
		costStr := strings.ReplaceAll(record[3], ",", "")
		cost, err := strconv.ParseFloat(costStr, 64)
		if err != nil {
			log.Printf("Error parsing cost in row %d: %v", i+2, err)
			continue
		}
		sku := record[4]

		// Prepare the SQL query
		query := "UPDATE store_product SET cost = ? WHERE sku = ?"

		// Execute the query
		_, err = db.Exec(query, cost, sku)
		if err != nil {
			log.Printf("Error updating product with SKU %s: %v", sku, err)
		} else {
			fmt.Printf("Updated product with SKU %s: set cost to %.2f\n", sku, cost)
		}
	}
}
