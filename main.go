package main

import (
	"andon-datasource/middleware"
	"andon-datasource/router"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

const (
	envName = "demo.env"
)

func TestFn() {
	// middleware.GetWorkOrderDetail("", "")
}

func main() {

	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading env file")
	} else {
		log.Println("using .env for ", envName)
	}

	middleware.Start()

	TestFn()

	r := router.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
