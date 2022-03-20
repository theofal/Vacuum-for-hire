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
	port, _          = strconv.Atoi(getDotEnvVar("PORT"))
)

const (
	url  = "https://www.google.com/search?&q=golang&ibp=htl;jobs"
	url2 = "https://www.google.com/search?=&q=golang&ibp=htl;jobs#htivrt=jobs&fpstate=tldetail&htichips=city:D7fiBh9u5kdglIxow4ILBA%3D%3D&htischips=city;D7fiBh9u5kdglIxow4ILBA%3D%3D:Paris_comma_%20IDF&htidocid=CbS7UKGsjJ4AAAAAAAAAAA%3D%3D"
)

type WebElementList struct {
	ElementsList []selenium.WebElement
}

func getGoogleUrl(termToSearch string) string {
	url := URL{
		"https://www.google.com/search?&q=",
		termToSearch,
		"&ibp=htl;jobs",
	}
	return url.Base + url.Term + url.Endpoint
}

// GoogleSelenium : Selenium instance to scrap data from Google.
func GoogleSelenium(termToSearch string) *WebElementList {
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
			//"--headless",
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

	// Consent to Google cookies.
	acceptButton, err := wd.FindElement(selenium.ByXPATH, "//form[//span[text()=contains(., 'accepte')]]")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acceptButton.Click()
	fmt.Println("Clicking Chrome accept button")

	locationButton, err := wd.FindElement(selenium.ByXPATH, "//*[@data-value='D7fiBh9u5kdglIxow4ILBA==']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	locationButton.Click()
	fmt.Println("Setting location to Paris")

	setDateTab, err := wd.FindElement(selenium.ByXPATH, "//*[@data-facet='date_posted' and @role='tab']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	setDateTab.Click()
	fmt.Println("Clicking on date tab")

	setDateToday, err := wd.FindElement(selenium.ByXPATH, "//*[@data-name='today']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	setDateToday.Click()
	fmt.Println("Setting the day to today")

	var jobList []selenium.WebElement

	tmp, index := -1, 0

	// Get the list of jobs as WebElements
	for index != tmp {
		jobList, err = wd.FindElements(selenium.ByXPATH, "//*[@role='treeitem']")
		if err != nil {
			break
		}
		if len(jobList) > 0 && index <= len(jobList)-1 {
			jobList[index].Click()
			jobTitle, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='BjJfJf PUpOsf']")
			company, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='vNEEBe']")
			location, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='Qk80Jf'][1]")
			jobLink, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='pMhGee Co68jc j0vryd']")
			fmt.Println(jobTitle[index].Text())
			fmt.Println(company[index].Text())
			fmt.Println(location[index].Text())
			fmt.Println(jobLink.GetAttribute("href"))
			//time.Sleep(time.Second)
			index++
		} else {
			time.Sleep(time.Second)
		}
		tmp++
	}

	return &WebElementList{ElementsList: jobList}
}
