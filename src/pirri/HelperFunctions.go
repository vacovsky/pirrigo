package main

//	"fmt"

func ConvertSqlDayToDOW(daynum int) string {
	result := ""
	if daynum == 0 {
		result = "Sunday"
	} else if daynum == 1 {
		result = "Monday"
	} else if daynum == 1 {
		result = "Tuesday"
	} else if daynum == 1 {
		result = "Wednesday"
	} else if daynum == 1 {
		result = "Thursday"
	} else if daynum == 1 {
		result = "Friday"
	} else if daynum == 1 {
		result = "Saturday"
	}
	return result
}
