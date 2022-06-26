package crawler

import (
	"github.com/tebeka/selenium"
	"github.com/theofal/Vacuum-for-hire/services"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// SearchGoogle : selenium steps to scrap data from Google.
func (wd *WebDriver) SearchGoogle(termToSearch string) ([]services.Post, error) {

	var jobList []selenium.WebElement
	var allJobs []services.Post

	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	services.Logger.Debug("Getting URL infos.")
	url := getGoogleURL(termToSearch) //termToSearch)

	defer func(service *selenium.Service) {
		err := service.Stop()
		if err != nil {
			services.Logger.Error("Error while stopping the service.", zap.Error(err))
		}
	}(&wd.Service)

	defer func(Webdriver selenium.WebDriver) {
		err := Webdriver.Quit()
		if err != nil {
			services.Logger.Error("Error while stopping the service.", zap.Error(err))
		}
	}(wd.Driver)

	// Navigate to the Google jobs website.
	services.Logger.Info("Navigating to Chrome website.")
	if err := wd.Driver.Get(url); err != nil {
		services.Logger.Error("Couldn't get the web page.", zap.Error(err))
		return nil, services.ErrTimedOut
	}

	// Consent to Google cookies.
	services.Logger.Debug("Trying to find Google cookies consent accept button")
	acceptButton, err := wd.Driver.FindElement(selenium.ByXPATH, "//form[//span[text()=contains(., 'accepte')]]")
	if err != nil {
		services.Logger.Error("Couldn't find the Google cookies consent accept button element.", zap.Error(err))
		return nil, services.ErrTimedOut
	}
	services.Logger.Debug("Clicking Google cookies consent accept button")
	err = acceptButton.Click()
	if err != nil {
		services.Logger.Error("Couldn't click on the Google cookies consent accept button element.", zap.Error(err))
		return nil, services.ErrTimedOut
	}

	services.Logger.Debug("Trying to find Google jobs location button")
	locationButton, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-value='D7fiBh9u5kdglIxow4ILBA==']")
	if err != nil {
		services.Logger.Error("Couldn't find Google jobs location button", zap.Error(err))
		return nil, services.ErrTimedOut
	}
	services.Logger.Debug("Clicking on Google jobs location button", zap.String("Location", "Paris"))
	err = locationButton.Click()
	if err != nil {
		services.Logger.Error("Couldn't click on Google jobs location button", zap.Error(err))
		return nil, services.ErrTimedOut
	}

	services.Logger.Debug("Getting date tab.")
	setDateTab, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-facet='date_posted' and @role='tab']")
	if err != nil {
		services.Logger.Error("Couldn't find date tab element.", zap.Error(err))
		return nil, services.ErrTimedOut
	}
	services.Logger.Debug("Clicking on date tab")
	err = setDateTab.Click()
	if err != nil {
		services.Logger.Error("Couldn't click on date tab", zap.Error(err))
		return nil, services.ErrTimedOut
	}

	services.Logger.Debug("Setting date range to today.")
	setDateToday, err := wd.Driver.FindElement(selenium.ByXPATH, "//*[@data-name='today']")
	if err != nil {
		services.Logger.Error("Couldn't find date button element.", zap.Error(err))
		return nil, services.ErrTimedOut
	}
	services.Logger.Debug("Clicking on \"Today\" date button")
	err = setDateToday.Click()
	if err != nil {
		services.Logger.Error("Couldn't click on \"Today\" date button", zap.Error(err))
		return nil, services.ErrTimedOut
	}

	// Get the list of jobs as WebElements
	tmp, index := -1, 0
	services.Logger.Debug("Getting the list of jobs.")
	for index != tmp {
		services.Logger.Debug("Trying to find job element.",
			zap.String("elementIndex", strconv.Itoa(index)),
			zap.String("tmp value", strconv.Itoa(tmp)),
		)
		jobList, err = wd.Driver.FindElements(selenium.ByXPATH, "//*[@role='treeitem']")
		if err != nil {
			services.Logger.Error("Couldn't find job element",
				zap.String("elementIndex", strconv.Itoa(index)),
				zap.Error(err))
			break
		}
		if len(jobList) > 0 && index <= len(jobList)-1 {
			err := jobList[index].Click()
			if err != nil {
				services.Logger.Warn("couldn't click on job element", zap.String("elementIndex", strconv.Itoa(index)), zap.Error(err))
			}
			jobTitleElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//*[@id='tl_ditsc']//*[@class='KLsYvd']")
			jobTitle, _ := jobTitleElement.Text()
			jobDateElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//*[@id='tl_ditsc']//*[@class='LL4CDc']")
			jobDate, _ := jobDateElement.GetAttribute("aria-label")
			companyElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//*[@id='tl_ditsc']//*[@class='nJlQNd sMzDkb']")
			companyName, _ := companyElement.Text()
			locationElement, _ := jobList[index].FindElements(selenium.ByXPATH, "//*[@class='Qk80Jf'][1]")
			companyLocation, _ := locationElement[index].Text()
			jobLinkElement, _ := jobList[index].FindElement(selenium.ByXPATH, "//*[@id='tl_ditsc']//*[@class='pMhGee Co68jc j0vryd']")
			jobLink, _ := jobLinkElement.GetAttribute("href")

			allJobs = append(allJobs,
				services.Post{
					JobTitle:        parseString(jobTitle),
					Date:            services.ParseDate(jobDate),
					CompanyName:     parseString(companyName),
					CompanyLocation: parseString(companyLocation),
					URL:             parseString(jobLink),
				})
			index++

		} else {
			time.Sleep(time.Second)
		}
		tmp++
	}
	services.Logger.Info("Done finding all the jobs !", zap.String("numberOfJobs", strconv.Itoa(len(allJobs))))
	return allJobs, err
}
