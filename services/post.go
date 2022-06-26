package services

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

//Post is the job structure.
type Post struct {
	ID              string `json:"id" gorm:"ID"`
	JobTitle        string `json:"jobTitle" gorm:"JobTitle"`
	CompanyName     string `json:"companyName" gorm:"CompanyName"`
	CompanyLocation string `json:"companyLocation" gorm:"CompanyLocation"`
	JobSnippet      string `json:"jobSnippet" gorm:"JobSnippet"`
	Date            string `json:"date" gorm:"Date"`
	URL             string `json:"url" gorm:"Url"`
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

//ParseToJSON parses an Post to a list of maps. This is used for the API.
func ParseToJSON(post []Post) []map[string]string {
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

// ParseDate returns a date in a string format to when the job post was uploaded.
func ParseDate(date string) string {
	var amount string
	Logger.Debug("Parsing date.", zap.String("Date", date))
	for _, v := range date {
		if unicode.IsDigit(v) {
			amount += string(v)
		}
	}
	intAmount, _ := strconv.Atoi(amount)
	timeMinusMinutes := time.Now().Add(-time.Minute * time.Duration(intAmount)).Format("02/01/2006 15:04")
	timeMinusHours := time.Now().Add(-time.Hour * time.Duration(intAmount)).Format("02/01/2006 15:04")
	timeMinusDays := time.Now().AddDate(0, 0, -2).Format("02/01/2006 15:04")
	switch {
	case date == "PostedPubliée à l'instant" || date == "PostedAujourd'hui":
		return time.Now().Format("02/01/2006 15:04")
	case strings.Contains(strings.ToLower(date), "minute"):
		return timeMinusMinutes
	case strings.Contains(strings.ToLower(date), "heure"):
		return timeMinusHours
	case strings.Contains(strings.ToLower(date), "jour"):
		return timeMinusDays
	}
	return fmt.Sprintf("Couldn't parse time \"%v\".", date)
}
