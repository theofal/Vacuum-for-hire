package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type URL struct {
	Base     string
	Term     string
	Endpoint string
}

type Post struct {
	JobTitle        string
	CompanyName     string
	CompanyLocation string
	JobSnippet      string
	Date            string
	Url             string
}

var (
	AllJobs []Post
)

func getDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	//fmt.Println(IndeedScrap(GetIndeedUrl(), GetIndeedJobs))
	//fmt.Println("AAAAAAA", len(GoogleSelenium("golang").ElementsList))
	GoogleSelenium("golang")
}
