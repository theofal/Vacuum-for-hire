package crawler

import (
	"fmt"
	"github.com/theofal/Vacuum-for-hire/services"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.uber.org/zap"
)

var (
	chromedriverPath = os.Getenv("WEBDRIVER_PATH")
	port, _          = strconv.Atoi(os.Getenv("PORT"))
)

//WebDriver object structure.
type WebDriver struct {
	Driver  selenium.WebDriver
	Service selenium.Service
}

// getGoogleURL returns a URL string to search jobs in.
func getGoogleURL(termToSearch string) string {
	if termToSearch == "" {
		services.Logger.Error("Empty termToSearch parameter: Disabled all jobs search as there would be too many elements.", zap.String("termToSearch", termToSearch))
		os.Exit(1)
	}
	url := services.URL{
		Base:     "https://www.google.com/search?&q=",
		Term:     termToSearch,
		Endpoint: "&ibp=htl;jobs",
	}
	return url.Base + url.Term + url.Endpoint
}

// Webdriver instance.
func Webdriver() *WebDriver {
	var opts []selenium.ServiceOption

	services.Logger.Info("Starting WebDriver creation.")

	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromedriverPath, port, opts...)
	if err != nil {
		services.Logger.Error("Unable to create new chromeDriver instance.", zap.Error(err))
		os.Exit(1)
	}

	// Connect to the WebDriver instance running locally.
	services.Logger.Debug("Setting browser name to Chrome.")
	caps := selenium.Capabilities{"browserName": "chrome"}

	services.Logger.Debug("Initialising Chrome settings.")
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			//"--start-fullscreen",
			"--window-size=2560,1440",
			"--disable-dev-shm-usage",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7",
		},
	}

	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		services.Logger.Error("Couldn't create new selenium remote client.", zap.Error(err))
		os.Exit(1)
	}

	err = wd.SetAsyncScriptTimeout(time.Second * 10)
	if err != nil {
		services.Logger.Warn("Couldn't set async scripts timeout.", zap.Error(err))
	}
	err = wd.SetPageLoadTimeout(time.Second * 10)
	if err != nil {
		services.Logger.Warn("Couldn't set page load timeout.", zap.Error(err))
	}
	err = wd.SetImplicitWaitTimeout(time.Second * 10)
	if err != nil {
		services.Logger.Warn("Couldn't set elements searching timeout.", zap.Error(err))
	}

	services.Logger.Info("WebDriver created.")
	return &WebDriver{Driver: wd, Service: *service}
}

// parseString removes recurrent unneeded substrings in Post strings.
func parseString(str string) string {
	str = strings.Replace(str, "<NIL>", "", -1)
	str = strings.Replace(str, "...", "", -1)
	return str
}
