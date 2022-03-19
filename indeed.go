package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func GetIndeedUrl() URL {
	return URL{
		"https://fr.indeed.com/jobs?q=",
		os.Args[1],
		"&l=France&sort=date&limit=50&fromage=1",
	}
}

func IndeedScrap(url URL, f func(doc *goquery.Document) []Post) []Post {
	// Request the HTML page.
	client := &http.Client{Timeout: time.Second * 20}
	res, err := client.Get(url.Base + url.Term + url.Endpoint)

	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return f(doc)
}

func GetJobs(selection *goquery.Document) []Post {
	var post Post
	selection.Find(".result").Each(func(i int, s *goquery.Selection) {
		url, isVisible := s.Attr("href")
		if isVisible == false {
			fmt.Println(fmt.Errorf("couldn't find url %v", isVisible))
		}
		post.Url = ParseIndeedUrl(url)
		post.JobTitle = s.Find("h2.jobTitle>span").Text()
		post.CompanyName = s.Find(".companyName").Text()
		post.CompanyLocation = s.Find(".companyLocation").Text()
		post.JobSnippet = s.Find(".job-snippet>ul>li").Text()
		post.Date = ParseIndeedDate(s.Find(".date").Text())
		AllJobs = append(AllJobs, post)
	})
	return AllJobs
}

func ParseIndeedUrl(url string) string {
	if url == "" {
		fmt.Println("No URL found")
	}
	if strings.Contains(url, "fccid") == false {
		return url
	}
	url = url[:strings.Index(url, "fccid")-1]
	if strings.Contains(url, "/rc/clk") {
		return "https://fr.indeed.com/viewjob" + url[7:]
	}
	return "https://fr.indeed.com/viewjob?jk=" + url[len(url)-16:]
}

func ParseIndeedDate(date string) string {
	if date == "PostedPubliée à l'instant" || date == "PostedAujourd'hui" {
		return fmt.Sprintf("%d/%d/%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	}
	date = strings.Replace(date, "Postedil y a ", "", 1)
	date = strings.Replace(date, " jour", "", 1)
	day, err := strconv.Atoi(date)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("%d/%d/%d", time.Now().Day()-day, time.Now().Month(), time.Now().Year())
}
