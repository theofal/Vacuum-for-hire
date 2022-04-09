package main

import (
	"testing"
)

func TestParseStructToArray(t *testing.T) {
	post := Post{
		Id:              "1",
		JobTitle:        "a",
		CompanyName:     "b",
		CompanyLocation: "c",
		JobSnippet:      "d",
		Date:            "e",
		Url:             "f",
	}
	want := [7]string{"1", "a", "b", "c", "d", "e", "f"}
	got := ParseStructToArray(post)
	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Errorf("TestParseStructToArray FAILED : want %v, got %v.\n", want, got)
		}
	}
}
