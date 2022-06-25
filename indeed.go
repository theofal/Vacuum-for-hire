package app

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//NOT USED IN PROJECT.

//GetIndeedUrl.
func _() URL {
	return URL{
		Base:     "https://fr.indeed.com/jobs?q=",
		Term:     os.Args[1],
		Endpoint: "&l=France&sort=date&limit=50&fromage=1",
	}
}

// IndeedScrap.
func _(url URL, f func(doc *goquery.Document) []Post) []Post {
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

//GetIndeedJobs.
func _(selection *goquery.Document) []Post {
	var post Post
	var allJobs []Post
	selection.Find(".result").Each(func(i int, s *goquery.Selection) {
		url, isVisible := s.Attr("href")
		if !isVisible {
			fmt.Println(fmt.Errorf("couldn't find url %v", isVisible))
		}
		post.URL = parseIndeedURL(url)
		post.JobTitle = s.Find("h2.jobTitle>span").Text()
		post.CompanyName = s.Find(".companyName").Text()
		post.CompanyLocation = s.Find(".companyLocation").Text()
		post.JobSnippet = s.Find(".job-snippet>ul>li").Text()
		post.Date = ParseDate(s.Find(".date").Text())
		allJobs = append(allJobs, post)
	})
	return allJobs
}

func parseIndeedURL(url string) string {
	if url == "" {
		fmt.Println("No URL found")
	}
	if !strings.Contains(url, "fccid") {
		return url
	}
	url = url[:strings.Index(url, "fccid")-1]
	if strings.Contains(url, "/rc/clk") {
		return "https://fr.indeed.com/viewjob" + url[7:]
	}
	return "https://fr.indeed.com/viewjob?jk=" + url[len(url)-16:]
}
