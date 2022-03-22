package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	chromedriverPath2 = getDotEnvVar("WEBDRIVER_PATH")
	port2, _          = strconv.Atoi(getDotEnvVar("PORT"))
)

type WebD struct {
	webdrive selenium.WebDriver
}

var web WebD

func getGoogleUrl2(termToSearch string) string {
	url := URL{
		"https://www.google.com/search?&q=",
		termToSearch,
		"&ibp=htl;jobs",
	}
	return url.Base + url.Term + url.Endpoint
}

func webdriver() *WebD {
	var opts []selenium.ServiceOption

	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromedriverPath, port, opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(service *selenium.Service) {
		err := service.Stop()
		if err != nil {
			//TODO LOG
		}
	}(service)

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--window-size=1920,1080",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
	}

	wd.SetAsyncScriptTimeout(time.Second * 10)
	wd.SetPageLoadTimeout(time.Second * 10)
	wd.SetImplicitWaitTimeout(time.Second * 10)

	web.webdrive = wd

	return &WebD{webdrive: wd}

}

// GoogleSelenium instance to scrap data from Google.
func (wd *WebD) GoogleSelenium2(termToSearch string) *[]Post {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	url := getGoogleUrl(termToSearch)

	defer wd.webdrive.Quit()

	// Navigate to the simple playground interface.
	if err := wd.webdrive.Get(url); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Consent to Google cookies.
	acceptButton, err := wd.webdrive.FindElement(selenium.ByXPATH, "//form[//span[text()=contains(., 'accepte')]]")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acceptButton.Click()
	log.Println("Clicking Chrome accept button")
	fmt.Println("Clicking Chrome accept button")

	locationButton, err := wd.webdrive.FindElement(selenium.ByXPATH, "//*[@data-value='D7fiBh9u5kdglIxow4ILBA==']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	locationButton.Click()
	fmt.Println("Setting location to Paris")

	setDateTab, err := wd.webdrive.FindElement(selenium.ByXPATH, "//*[@data-facet='date_posted' and @role='tab']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	setDateTab.Click()
	fmt.Println("Clicking on date tab")

	setDateToday, err := wd.webdrive.FindElement(selenium.ByXPATH, "//*[@data-name='today']")
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
		jobList, err = wd.webdrive.FindElements(selenium.ByXPATH, "//*[@role='treeitem']")
		if err != nil {
			break
		}
		if len(jobList) > 0 && index <= len(jobList)-1 {
			jobList[index].Click()
			jobTitleElement, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='BjJfJf PUpOsf']")
			jobTitle, _ := jobTitleElement[index].Text()
			companyElement, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='vNEEBe']")
			companyName, _ := companyElement[index].Text()
			locationElement, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='Qk80Jf'][1]")
			companyLocation, _ := locationElement[index].Text()
			jobLinkElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='pMhGee Co68jc j0vryd']")
			jobLink, _ := jobLinkElement.GetAttribute("href")

			AllJobs = append(AllJobs,
				Post{
					JobTitle:        strings.Replace(jobTitle, "<NIL>", "", 1),
					CompanyName:     strings.Replace(companyName, "<NIL>", "", 1),
					CompanyLocation: strings.Replace(companyLocation, "<NIL>", "", 1),
					Url:             strings.Replace(jobLink, "<NIL>", "", 1),
				})
			index++

		} else {
			time.Sleep(time.Second)
		}
		tmp++
	}

	return &AllJobs
}

