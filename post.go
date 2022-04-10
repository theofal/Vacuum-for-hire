package main

import (
	"reflect"
)

//Post is the job structure.
type Post struct {
	ID              string `json:"ID" gorm:"ID"`
	JobTitle        string `json:"JobTitle" gorm:"JobTitle"`
	CompanyName     string `json:"CompanyName" gorm:"CompanyName"`
	CompanyLocation string `json:"CompanyLocation" gorm:"CompanyLocation"`
	JobSnippet      string `json:"JobSnippet" gorm:"JobSnippet"`
	Date            string `json:"Date" gorm:"Date"`
	URL             string `json:"Url" gorm:"Url"`
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

//ParseToJson parses an Post to a list of maps. This is used for the API.
func ParseToJson(post []Post) []map[string]string {
	jsonList := make([]map[string]string, 0)
	for i := 0; i < len(post); i++ {
		j := map[string]string{
			"ID":              post[i].ID,
			"JobTitle":        post[i].JobTitle,
			"CompanyName":     post[i].CompanyName,
			"CompanyLocation": post[i].CompanyLocation,
			"JobSnippet":      post[i].JobSnippet,
			"Date":            post[i].Date,
			"URL":             post[i].URL,
		}
		jsonList = append(jsonList, j)
	}
	return jsonList
}
