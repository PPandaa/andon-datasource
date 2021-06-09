package main

import (
	"DataSource/config"
	"DataSource/router"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func initGlobalVar() {
	err := godotenv.Load(config.EnvPath)
	if err != nil {
		log.Fatalf("Error loading env file")
	}

	ensaasService := os.Getenv("ENSAAS_SERVICES")
	if len(ensaasService) != 0 {
		tempReader := strings.NewReader(ensaasService)
		m, _ := simplejson.NewFromReader(tempReader)
		mongodb := m.Get("mongodb").GetIndex(0).Get("credentials").MustMap()
		config.MongodbURL = mongodb["externalHosts"].(string)
		config.MongodbDatabase = mongodb["database"].(string)
		config.MongodbUsername = mongodb["username"].(string)
		config.MongodbPassword = mongodb["password"].(string)
	} else {
		config.MongodbURL = os.Getenv("MONGODB_URL")
		config.MongodbDatabase = os.Getenv("MONGODB_DATABASE")
		config.MongodbUsername = os.Getenv("MONGODB_USERNAME")
		config.MongodbAuthSource = os.Getenv("MONGODB_AUTH_SOURCE")
		mongodbPasswordFile := os.Getenv("MONGODB_PASSWORD_FILE")
		if len(mongodbPasswordFile) != 0 {
			mongodbPassword, err := ioutil.ReadFile(mongodbPasswordFile)
			if err != nil {
				fmt.Println("MongoDB Password File", err, "->", "FilePath:", mongodbPasswordFile)
			} else {
				config.MongodbPassword = string(mongodbPassword)
			}
		}
	}

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
	if len(ensaasService) != 0 {
		config.DB = config.Session.DB(config.MongodbDatabase)
		config.DB.Login(config.MongodbUsername, config.MongodbPassword)
	} else {
		config.DB = config.Session.DB(config.MongodbAuthSource)
		config.DB.Login(config.MongodbUsername, config.MongodbPassword)
		config.DB = config.Session.DB(config.MongodbDatabase)
	}
	fmt.Println("----------", time.Now().In(config.TaipeiTimeZone), "----------")
	fmt.Println("MongoDB Connect ->", " URL:", config.MongodbURL, " Database:", config.MongodbDatabase)
	fmt.Println(config.DB.CollectionNames())
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
