package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	_ "github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"strconv"
	"time"
)

var (
	chromedriverPath = getDotEnvVar("WEBDRIVER_PATH")
	port, err        = strconv.Atoi(getDotEnvVar("PORT"))
)

func getGoogleUrl(termToSearch string) string {
	url := URL{
		"https://www.google.fr/search?client=firefox-b-d&q=",
		termToSearch,
		"&ibp=htl;jobs#fpstate=tldetail&htivrt=jobs&htichips=date_posted:today,city:D7fiBh9u5kdglIxow4ILBA%3D%3D&htischips=date_posted;today,city;D7fiBh9u5kdglIxow4ILBA%3D%3D:Paris_comma_ IDF",
	}
	return url.Base + url.Term + url.Endpoint
}

func GoogleSelenium(termToSearch string) []selenium.WebElement {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	url := getGoogleUrl(termToSearch)

	var opts []selenium.ServiceOption

	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromedriverPath, port, opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
	}
	defer wd.Quit()

	wd.SetAsyncScriptTimeout(time.Second * 10)
	wd.SetPageLoadTimeout(time.Second * 10)
	wd.SetImplicitWaitTimeout(time.Second * 10)

	// Navigate to the simple playground interface.
	if err := wd.Get(url); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Consent to cookies
	acceptButton, err := wd.FindElement(selenium.ByXPATH, "//form[//span[text()=contains(., 'accepte')]]")
	if err != nil {
		fmt.Println(err)
	}
	acceptButton.Click()
	fmt.Println("Clicking Chrome accept button")

	if err := wd.Get(url); err != nil {
		panic(err)
	}

	var jobList []selenium.WebElement

	tmp, listLength := 0, 1

	for listLength != tmp {
		jobList, err = wd.FindElements(selenium.ByXPATH, "//*[@role='treeitem']")
		if err != nil {
			break
		}
		if len(jobList) > 0 {
			jobList[len(jobList)-1].Click()
			time.Sleep(time.Second)
			tmp = listLength
			listLength = len(jobList)
			fmt.Printf("TMP: %d LEN: %d\n", tmp, listLength)
		} else {
			time.Sleep(time.Second)
		}
	}

	fmt.Printf("final job list size: %d\n", len(jobList))
	return jobList
}
