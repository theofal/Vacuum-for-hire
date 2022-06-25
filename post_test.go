package main

import (
	"reflect"
	"testing"
)

func TestParseStructToArray(t *testing.T) {
	post := Post{
		ID:              "1",
		JobTitle:        "a",
		CompanyName:     "b",
		CompanyLocation: "c",
		JobSnippet:      "d",
		Date:            "e",
		URL:             "f",
	}
	want := [7]string{"1", "a", "b", "c", "d", "e", "f"}
	got := ParseStructToArray(post)
	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Errorf("TestParseStructToArray FAILED : want %v, got %v.\n", want, got)
		}
	}
}

func TestParseToJson(t *testing.T) {
	jobList := make([]Post, 1)
	testPost := Post{
		ID:              "a",
		JobTitle:        "b",
		CompanyName:     "c",
		CompanyLocation: "d",
		JobSnippet:      "e",
		Date:            "f",
		URL:             "g",
	}
	jobList[0] = testPost

	mapJob := map[string]string{"ID": "a", "JobTitle": "b", "CompanyName": "c", "CompanyLocation": "d", "JobSnippet": "e", "Date": "f", "URL": "g"}
	want := make([]map[string]string, 1)
	want[0] = mapJob
	got := ParseToJson(jobList)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("TestParseStructToArray FAILED : want %v, got %v.\n", want, got)
	}
}
