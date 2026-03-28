package cmd

import (
	"testing"
	"time"
)

func TestCountDays(t *testing.T) {
	cases := [][]time.Time{
		{string2Date("12/12", 2025), string2Date("12/12", 2025)},
		{string2Date("12/14", 2025), string2Date("12/12", 2025)},
		{string2Date("12/12", 2025), string2Date("12/14", 2025)},
		// Start of daylight saving
		{string2Date("3/8", 2025), string2Date("3/10", 2025)},
		{string2Date("3/9", 2025), string2Date("3/10", 2025)},
		// End of daylight saving
		{string2Date("11/1", 2025), string2Date("11/3", 2025)},
		{string2Date("11/2", 2025), string2Date("11/3", 2025)},
		// Test different years
		{string2Date("12/30", 2025), string2Date("1/1", 2026)},
	}
	expected := []int {
		0,
		2,
		2,
		// Start of daylight saving
		2,
		1,
		// End of daylight saving
		2,
		1,
		// Test different years
		2,
	}
	for i := range cases {
		actual := countDays(cases[i][0], cases[i][1])
		if expected[i] != actual {
			t.Errorf("Test case %d failed. Expected: %d. Actual: %d", i, expected[i], actual)
		}
	}
}

// func TestStrings2Dates(t *testing.T) {
// 	cases := [][]string{
// 		{},
// 	}
// }
