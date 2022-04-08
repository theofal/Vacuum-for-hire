package main

import "reflect"

//Post is the job structure.
type Post struct {
	id              string
	jobTitle        string
	companyName     string
	companyLocation string
	jobSnippet      string
	date            string
	url             string
}

//ID getter of Post struct.
func (post *Post) ID() string {
	return post.id
}

//CompanyName getter of Post struct.
func (post *Post) CompanyName() string {
	return post.companyName
}

//JobTitle getter of Post struct.
func (post *Post) JobTitle() string {
	return post.jobTitle
}

//CompanyLocation getter of Post struct.
func (post *Post) CompanyLocation() string {
	return post.companyLocation
}

//JobSnippet getter of Post struct.
func (post *Post) JobSnippet() string {
	return post.jobSnippet
}

//Date getter of Post struct.
func (post *Post) Date() string {
	return post.date
}

//URL getter of Post struct.
func (post *Post) URL() string {
	return post.url
}

//ParseStructToArray parses an interface (in our case, a Post) to a list of strings.
func ParseStructToArray(post interface{}) []string {
	dataValue := reflect.ValueOf(post)
	dataArray := make([]string, dataValue.NumField())
	for i := 0; i < dataValue.NumField(); i++ {
		dataArray[i] = dataValue.Field(i).String()
	}

	return dataArray
}
