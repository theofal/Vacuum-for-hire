package main

type Post struct {
	jobTitle        string
	companyName     string
	companyLocation string
	jobSnippet      string
	date            string
	url             string
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
