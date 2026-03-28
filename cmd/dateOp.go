/*
Copyright © 2025 ganlinden@gmail.com
*/
package cmd

import (
	"math"
	"time"
)

func date2String(date time.Time) string {
	return date.Format("1/2")
}

/*
Possible input str: 8/3, 1/17, 12/9, 04/24, 12/31
*/
func string2Date(str string, year int) time.Time {
	res, err := time.Parse("1/2", str)
	if err != nil {
		panic(err)
	}
	res = res.AddDate(year, 0, 0)
	return res
}

/*
Deduce year based on the assumption that:
The date is in the range of today plus minus half a year.
Example 1: if today is 2025/8/1, then we assume 1/1 as 2026/1/1,
and assume 3/1 as 2025/3/1.
Example 2: if today is 2025/2/1, then we assume 11/1 as 2024/11/1,
and assume 7/1 as 2025/7/1.
*/
const halfYear time.Duration = 182 * 24 * time.Hour
func deduceYear(date string) int {
	year := time.Now().Year()
	tentative := string2Date(date, year)
    if time.Since(tentative) > halfYear {
		year++
	} else if time.Until(tentative) > halfYear {
		year--
	}
	return year
}

/*
Map a slice of mmdd strings to a slice of time.Time.
Example input: [12/25, 12/30, 1/2, 1/6] (say today is 12/31)
*/
func strings2Dates(input []string) []time.Time {
	res := make([]time.Time, len(input))
	if len(input) == 0 {
		return res
	}
	year := deduceYear(input[len(input) - 1])
    // Reversely traverse input and update year if needed
	i := len(input) - 1
	res[i] = string2Date(input[i], year)
	for i = len(input) - 2; i >= 0; i-- {
		currDate := string2Date(input[i], year)
		laterDate := res[i + 1]
		// If currDate is 12/31, laterDate is 1/1, then we know the year
		// becomes the previous year.
		if currDate.After(laterDate) {
			year--
			currDate = string2Date(input[i], year)
		}
		res[i] = currDate
	}
	return res
}

func dates2Strings(input []time.Time) []string {
	res := make([]string, len(input))
	for i, date := range input {
		res[i] = date2String(date)
	}
	return res
}

func countDays(a, b time.Time) int {
	// Normalize to midnight
	a = string2Date(date2String(a), a.Year())
	b = string2Date(date2String(b), b.Year())
	if a.After(b) {
		a, b = b, a
	}
	diff := b.Sub(a)
	days := math.Round(float64(diff) / float64(time.Hour) / 24.0)
	return int(days)
}
