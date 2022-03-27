package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	chromedriverPath = getDotEnvVar("WEBDRIVER_PATH")
	port, _          = strconv.Atoi(getDotEnvVar("PORT"))
	Web              WebDriver
)

type WebDriver struct {
	Driver  selenium.WebDriver
	Service selenium.Service
}

// getGoogleURL returns a URL string to search jobs in.
func getGoogleURL(termToSearch string) string {
	if termToSearch == "" {
		Logger.Fatal("Empty termToSearch parameter: Disabled all jobs search as there would be too many elements.", zap.Any("termToSearch", termToSearch))
	}
	url := URL{
		"https://www.google.com/search?&q=",
		termToSearch,
		"&ibp=htl;jobs",
	}
	return url.Base + url.Term + url.Endpoint
}

// Webdriver instance.
func Webdriver() *WebDriver {
	var opts []selenium.ServiceOption

	Logger.Info("Starting WebDriver creation.")

	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromedriverPath, port, opts...)
	if err != nil {
		Logger.Error("Unable to create new chromeDriver instance.", zap.Error(err))
		os.Exit(1)
	}

	// Connect to the WebDriver instance running locally.
	Logger.Debug("Setting browser name to Chrome.")
	caps := selenium.Capabilities{"browserName": "chrome"}

	Logger.Debug("Initialising Chrome settings.")
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
		Logger.Error("Couldn't create new selenium remote client.", zap.Error(err))
		os.Exit(1)
	}

	err = wd.SetAsyncScriptTimeout(time.Second * 10)
	if err != nil {
		Logger.Warn("Couldn't set async scripts timeout.", zap.Error(err))
	}
	err = wd.SetPageLoadTimeout(time.Second * 10)
	if err != nil {
		Logger.Warn("Couldn't set page load timeout.", zap.Error(err))
	}
	err = wd.SetImplicitWaitTimeout(time.Second * 10)
	if err != nil {
		Logger.Warn("Couldn't set elements searching timeout.", zap.Error(err))
	}

	Web.Driver = wd
	Web.Service = *service

	Logger.Info("WebDriver created.")
	return &WebDriver{Driver: wd, Service: *service}

}

// SearchGoogle : selenium steps to scrap data from Google.
func (wd *WebDriver) SearchGoogle(termToSearch string) (*[]Post, error) {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	Logger.Debug("Getting URL infos.")
	url := getGoogleURL(termToSearch) //termToSearch)

	defer func(service *selenium.Service) {
		err := service.Stop()
		if err != nil {
			Logger.Error("Error while stopping the service.", zap.Error(err))
		}
	}(&wd.Service)

	defer func(Webdriver selenium.WebDriver) {
		err := Webdriver.Quit()
		if err != nil {
			Logger.Error("Error while stopping the service.", zap.Error(err))
		}
	}(wd.Driver)

	// Navigate to the Google jobs website.
	Logger.Info("Navigating to Chrome website.")
	if err := wd.Driver.Get(url); err != nil {
		Logger.Error("Couldn't get the web page.", zap.Error(err))
		return nil, ErrTimedOut
	}

	// Consent to Google cookies.
	Logger.Debug("Trying to find Google cookies consent accept button")
	acceptButton, err := wd.Driver.FindElement(selenium.ByXPATH, "//form[//span[text()=contains(., 'accepte')]]")
	if err != nil {
		Logger.Error("Couldn't find the Google cookies consent accept button element.", zap.Error(err))
		return nil, ErrTimedOut
	}
	Logger.Debug("Clicking Google cookies consent accept button")
	err = acceptButton.Click()
	if err != nil {
		Logger.Error("Couldn't click on the Google cookies consent accept button element.", zap.Error(err))
		return nil, ErrTimedOut
	}

	Logger.Debug("Trying to find Google jobs location button")
	locationButton, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-value='D7fiBh9u5kdglIxow4ILBA==']")
	if err != nil {
		Logger.Error("Couldn't find Google jobs location button", zap.Error(err))
		return nil, ErrTimedOut
	}
	Logger.Debug("Clicking on Google jobs location button", zap.String("Location", "Paris"))
	err = locationButton.Click()
	if err != nil {
		Logger.Error("Couldn't click on Google jobs location button", zap.Error(err))
		return nil, ErrTimedOut
	}

	Logger.Debug("Getting date tab.")
	setDateTab, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-facet='date_posted' and @role='tab']")
	if err != nil {
		Logger.Error("Couldn't find date tab element.", zap.Error(err))
		return nil, ErrTimedOut
	}
	Logger.Debug("Clicking on date tab")
	err = setDateTab.Click()
	if err != nil {
		Logger.Error("Couldn't click on date tab", zap.Error(err))
		return nil, ErrTimedOut
	}

	Logger.Debug("Setting date range to today.")
	setDateToday, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-name='today']")
	if err != nil {
		Logger.Error("Couldn't find date button element.", zap.Error(err))
		return nil, ErrTimedOut
	}
	Logger.Debug("Clicking on \"Today\" date button")
	err = setDateToday.Click()
	if err != nil {
		Logger.Error("Couldn't click on \"Today\" date button", zap.Error(err))
		return nil, ErrTimedOut
	}

	var jobList []selenium.WebElement

	tmp, index := -1, 0

	// Get the list of jobs as WebElements
	Logger.Debug("Getting the list of jobs.")
	for index != tmp {
		Logger.Debug("Trying to find job element.",
			zap.String("elementIndex", strconv.Itoa(index)),
			zap.String("tmp value", strconv.Itoa(tmp)),
		)
		jobList, err = wd.Driver.FindElements(selenium.ByXPATH, "//*[@role='treeitem']")
		if err != nil {
			Logger.Error("Couldn't find job element",
				zap.String("elementIndex", strconv.Itoa(index)),
				zap.Error(err))
			break
		}
		if len(jobList) > 0 && index <= len(jobList)-1 {
			err := jobList[index].Click()
			if err != nil {
				Logger.Warn("couldn't click on job element", zap.String("elementIndex", strconv.Itoa(index)), zap.Error(err))
			}
			jobTitleElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='KLsYvd']")
			jobTitle, _ := jobTitleElement.Text()
			jobDateElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='LL4CDc']")
			jobDate, _ := jobDateElement.GetAttribute("aria-label")
			companyElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='nJlQNd sMzDkb']")
			companyName, _ := companyElement.Text()
			locationElement, _ := jobList[index].FindElements(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@class='Qk80Jf'][1]")
			companyLocation, _ := locationElement[index].Text()
			jobLinkElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//body[*[*[div[@class='gb_Fc gb_Dc gb_Kc']]]]//*[@id='tl_ditsc']//*[@class='pMhGee Co68jc j0vryd']")
			jobLink, _ := jobLinkElement.GetAttribute("href")

			AllJobs = append(AllJobs,
				Post{
					JobTitle:        ParseString(jobTitle),
					Date:            ParseDate(jobDate),
					CompanyName:     ParseString(companyName),
					CompanyLocation: ParseString(companyLocation),
					URL:             ParseString(jobLink),
				})
			index++

		} else {
			time.Sleep(time.Second)
		}
		tmp++
	}
	Logger.Info("Done finding all the jobs !", zap.String("numberOfJobs", strconv.Itoa(len(AllJobs))))
	return &AllJobs, err
}

// ParseString removes recurrent unneeded substrings in Post strings.
func ParseString(str string) string {
	str = strings.Replace(str, "<NIL>", "", -1)
	str = strings.Replace(str, "...", "", -1)
	return str
}
