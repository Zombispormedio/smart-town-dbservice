package utils

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
  
)



func Notify(date time.Time) bool {
	var result bool
    
    LIMIT := os.Getenv("NOTIFY_PUSH_LIMIT")
  
    
	hour, _ := regexp.Compile(`^(\d){1,2}h$`)
	min, _ := regexp.Compile(`^(\d){1,2}min$`)

	duration := time.Since(date)

	switch {
	case hour.MatchString(LIMIT):
		hourDuration, _ := strconv.ParseFloat(strings.Replace(LIMIT, "h", "", 1), 64)
        
        result=duration.Hours()>hourDuration
        

	case min.MatchString(LIMIT):
		minDuration, _ := strconv.ParseFloat(strings.Replace(LIMIT, "min", "", 1),64)
 
        result=duration.Minutes()>minDuration
	}

	return result

}
