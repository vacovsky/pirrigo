package helpers

import (
	"fmt"
)

func convertSQLDayToDOW(daynum int) string {
	result := ""
	if daynum == 0 {
		result = "Sunday"
	} else if daynum == 1 {
		result = "Monday"
	} else if daynum == 2 {
		result = "Tuesday"
	} else if daynum == 3 {
		result = "Wednesday"
	} else if daynum == 4 {
		result = "Thursday"
	} else if daynum == 5 {
		result = "Friday"
	} else if daynum == 6 {
		result = "Saturday"
	} else {
		fmt.Println("Number provided does not correspond to a day of the week.  Valid inputs are 0-6.")
	}
	return result
}
