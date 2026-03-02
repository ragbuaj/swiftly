package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		log.Fatal(err)
	}

	products := []Product{
		{
			Name:        "Modern T-Shirt",
			Description: "Comfortable cotton t-shirt for everyday wear",
			Price:       150000,
			Stock:       100,
		},
		{
			Name:        "Slim Fit Jeans",
			Description: "High-quality denim for a perfect fit",
			Price:       450000,
			Stock:       50,
		},
		{
			Name:        "Wireless Earbuds",
			Description: "High-fidelity sound with noise cancellation",
			Price:       1200000,
			Stock:       30,
		},
	}

	log.Printf("Seeding %d products...", len(products))

	// In Supabase, the data must be a slice of maps or a struct with json tags
	_, _, err = client.From("products").Insert(products, false, "", "", "").Execute()
	if err != nil {
		log.Fatalf("Error seeding products: %v", err)
	}

	log.Println("Successfully seeded database!")
}
