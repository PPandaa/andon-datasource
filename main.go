package main

import (
	"andon-datasource/middleware"
	"andon-datasource/router"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const (
	envName = "demo.env"
)

func main() {
	middleware.TestDatasourceFn()
	log.Println("Activation")
	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading env file")
	}
	mongodbURL := os.Getenv("MONGODB_URL")
	mongodbDatabase := os.Getenv("MONGODB_DATABASE")
	fmt.Println("Version ->", "2021/1/7 16:07")
	fmt.Println("MongoDB ->", "URL:", mongodbURL, " Database:", mongodbDatabase)
	r := router.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
