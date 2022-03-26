package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
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
	Logger  *zap.Logger
)

func getDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		Logger.Error("Couldn't find .env find")
	}
	return os.Getenv(key)
}

func ParseDate(str string) string {
	var amount string
	Logger.Info("Working on it!")
	for _, v := range str {
		if unicode.IsDigit(v) {
			amount += string(v)
			Logger.Debug("Working.. ", zap.String(
				"Value of V", string(v)),
			)
		} else {
			Logger.Debug("Working.. ", zap.String(
				"Value of V", string(v)))
		}
	}
	intAmount, _ := strconv.Atoi(amount)
	Logger.Warn("Careful")
	timeNow := time.Now()
	fmt.Println(timeNow)
	switch {
	case strings.Contains(str, "minutes") || strings.Contains(str, "minute"):
		fmt.Printf("%v/%v/%v\n", timeNow.Add(-time.Minute*time.Duration(intAmount)).Day(), timeNow.Add(-time.Minute*time.Duration(intAmount)).Month(), timeNow.Add(-time.Minute*time.Duration(intAmount)).Year())
		return ""
	case strings.Contains(str, "heures") || strings.Contains(str, "heure"):
		//fmt.Println(timeNow.Add(-time.Hour * time.Duration(intAmount)))
		return fmt.Sprintf("%d/%d/%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	}
	return "Couldn't find"
}

func main() {
	InitLogger()
	defer Logger.Sync()
	fmt.Println(Webdriver().SearchGoogle("Golang"))
}
