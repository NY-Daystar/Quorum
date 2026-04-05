package utils

import (
	"fmt"
	"regexp"
	"time"
)

// ParseDate process date from mail's headers
func ParseDate(dateStr string) (string, error) {
	matchTz, _ := regexp.MatchString(`-|\+`, dateStr)

	re := regexp.MustCompile(`(.*)\s(-|\+)`)
	if !matchTz {
		re = regexp.MustCompile(`(.*)\s\(?[A-Z]*`)
	}

	matches := re.FindStringSubmatch(dateStr)
	if len(matches) < 2 {
		return "", fmt.Errorf("Parsing Error with this date %v", dateStr)
	}
	dateTime := re.FindStringSubmatch(dateStr)[1]

	t0, err0 := time.Parse("2 Jan 2006 15:04:05", dateTime)
	t1, err1 := time.Parse("02 Jan 2006 15:04:05", dateTime)
	t2, err2 := time.Parse("Mon, 2 Jan 2006 15:04:05", dateTime)
	t3, err3 := time.Parse("Mon, 02 Jan 2006 15:04:05", dateTime)

	if err0 != nil && err1 != nil && err2 != nil && err3 != nil {
		err := fmt.Sprintf("%v\n%v\n%v\n%v\n", err0, err1, err2, err3)
		fmt.Printf("err0: %v\n", err0)
		fmt.Printf("err1: %v\n", err1)
		fmt.Printf("err2: %v\n", err2)
		fmt.Printf("err3: %v\n", err3)
		return "", fmt.Errorf("%v", err)
	}

	var t time.Time
	if err0 == nil {
		t = t0
	} else if err1 == nil {
		t = t1
	} else if err2 == nil {
		t = t2
	} else if err3 == nil {
		t = t3
	}

	layout := "2006-01-02_15-04-05"
	return t.Format(layout), nil
}
