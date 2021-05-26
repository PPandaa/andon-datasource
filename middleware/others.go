package middleware

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

const (
	envName = "demo.env"
)

func getDBInfo() (string, string, string) {
	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	database := os.Getenv("MONGODB_DATABASE")
	return username, password, database
}

func createMongoSession() *mgo.Session {
	// load .env file
	err := godotenv.Load(envName)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	addr := os.Getenv("MONGODB_URL")
	// Open the connection
	session, _ := mgo.Dial(addr)
	// return the connection
	return session
}

func closeMongoSession(Session *mgo.Session) {
	Session.Close()
}
