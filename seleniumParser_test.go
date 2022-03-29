package main

import (
	"fmt"
	"github.com/tebeka/selenium"
	"os"
	"testing"
)

func TestGetGoogleURL(t *testing.T) {
	got := getGoogleURL("Test")
	want := "https://www.google.com/search?&q=Test&ibp=htl;jobs"
	if got != want {
		t.Errorf("TestGetGoogleURL FAILED : want %v, got %v.\n", want, got)
	}
}

func TestWebdriver(t *testing.T) {
	a := Webdriver()
	defer a.Service.Stop() //nolint:errcheck
	defer a.Driver.Quit()  //nolint:errcheck

	err := a.Driver.Get("https://theofal.github.io")
	if err != nil {
		fmt.Println(err)
	}
	name, err := a.Driver.FindElement(selenium.ByXPATH, "//*[@class='navbar-brand']")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	got, err := name.Text()
	if err != nil {
		fmt.Println(err)
	}

	want := "Mon site"
	if got != want {
		t.Errorf("TestWebdriver FAILED : want %v, got %v.\n", want, got)
	}
}

func TestParseString(t *testing.T) {
	got := parseString("DAYUM...<NIL>")
	want := "DAYUM"
	if got != want {
		t.Errorf("TestWebdriver FAILED : want %v, got %v.\n", want, got)
	}

	got = parseString("<NIL>...<NIL>DAYUM...<NIL>...")
	want = "DAYUM"
	if got != want {
		t.Errorf("TestWebdriver FAILED : want %v, got %v.\n", want, got)
	}
}

func TestSearchGoogle(t *testing.T) {
	TermToSearch = "Golang"
	Webdriver().SearchGoogle(TermToSearch)
	if len(AllJobs) <= 0 || AllJobs == nil {
		t.Errorf("Empty or nil list of jobs")
	}
	if AllJobs[0].JobTitle == "<NIL>" {
		t.Errorf("Got <NIL> jobtitle from the first element of AllJobs")
	}
}
