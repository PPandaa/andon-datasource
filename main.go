package main

import (
	"DataSource/config"
	"DataSource/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func initGlobalVar() {
	err := godotenv.Load(config.EnvPath)
	if err != nil {
		log.Fatalf("Error loading env file")
	}

	config.MongodbURL = os.Getenv("MONGODB_URL")
	config.MongodbDatabase = os.Getenv("MONGODB_DATABASE")
	config.MongodbUsername = os.Getenv("MONGODB_USERNAME")
	config.MongodbPassword = os.Getenv("MONGODB_PASSWORD")

	newSession, err := mgo.Dial(config.MongodbURL)
	if err != nil {
		fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
		fmt.Println("MongoDB", err, "->", "URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
		for err != nil {
			newSession, err = mgo.Dial(config.MongodbURL)
			time.Sleep(5 * time.Second)
		}
	}
	config.Session = newSession
	config.DB = config.Session.DB(config.MongodbDatabase)
	config.DB.Login(config.MongodbUsername, config.MongodbPassword)
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("MongoDB Connect ->", " URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
}

func main() {
	initGlobalVar()
	r := router.Router()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
