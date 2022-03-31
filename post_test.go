package main

import (
	"testing"
)

func TestCompanyName(t *testing.T) {
	post := Post{companyName: "Company Name Test"}
	want := "Company Name Test"
	got := post.CompanyName()
	if got != want {
		t.Errorf("TestCompanyName FAILED : want %v, got %v.\n", want, got)
	}
}

func TestJobTitle(t *testing.T) {
	post := Post{jobTitle: "Job Title Test"}
	want := "Job Title Test"
	got := post.JobTitle()
	if got != want {
		t.Errorf("TestJobTitle FAILED : want %v, got %v.\n", want, got)
	}
}

func TestCompanyLocation(t *testing.T) {
	post := Post{companyLocation: "Company Location Test"}
	want := "Company Location Test"
	got := post.CompanyLocation()
	if got != want {
		t.Errorf("TestCompanyLocation FAILED : want %v, got %v.\n", want, got)
	}
}

func TestJobSnippet(t *testing.T) {
	post := Post{jobSnippet: "Job Snippet Test"}
	want := "Job Snippet Test"
	got := post.JobSnippet()
	if got != want {
		t.Errorf("TestJobSnippet FAILED : want %v, got %v.\n", want, got)
	}
}

func TestDate(t *testing.T) {
	post := Post{date: "Date Test"}
	want := "Date Test"
	got := post.Date()
	if got != want {
		t.Errorf("TestDate FAILED : want %v, got %v.\n", want, got)
	}
}

func TestURL(t *testing.T) {
	post := Post{url: "URL Test"}
	want := "URL Test"
	got := post.URL()
	if got != want {
		t.Errorf("TestURL FAILED : want %v, got %v.\n", want, got)
	}
}
