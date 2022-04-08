package main

import "reflect"

type Post struct {
	id              string
	jobTitle        string
	companyName     string
	companyLocation string
	jobSnippet      string
	date            string
	url             string
}

func (post *Post) ID() string {
	return post.id
}

func (post *Post) CompanyName() string {
	return post.companyName
}

func (post *Post) JobTitle() string {
	return post.jobTitle
}

func (post *Post) CompanyLocation() string {
	return post.companyLocation
}

func (post *Post) JobSnippet() string {
	return post.jobSnippet
}

func (post *Post) Date() string {
	return post.date
}

func (post *Post) URL() string {
	return post.url
}

//ParseStructToArray parses an interface to a list of strings.
func ParseStructToArray(post interface{}) []string {
	dataValue := reflect.ValueOf(post)
	dataArray := make([]string, dataValue.NumField())
	for i := 0; i < dataValue.NumField(); i++ {
		dataArray[i] = dataValue.Field(i).String()
	}

	return dataArray
}
