package gox

import (
	"testing"
)

func TestTimeUnitStr(t *testing.T) {
	// Test cases: milliseconds value and expected string representation
	cases := []struct {
		ms       int64
		expected string
	}{
		{500, "500ms"},      // Less than 1 second
		{1500, "1.5s"},      // Seconds
		{60000, "1.0m"},     // Exactly 1 minute
		{90000, "1.5m"},     // 1.5 minutes
		{3600000, "1.0h"},   // Exactly 1 hour
		{5400000, "1.5h"},   // 1.5 hours
		{86400000, "1.0d"},  // Exactly 1 day
		{129600000, "1.5d"}, // 1.5 days
		{0, "0ms"},          // Zero value
		{3599999, "60.0m"},  // Just under 1 hour
		{86399999, "23.9h"}, // Just under 1 day
	}

	for i, tc := range cases {
		result := TimeUnitStr(tc.ms)
		if result != tc.expected {
			t.Errorf("Test case %d failed: for %dms expected '%s', got '%s'", i, tc.ms, tc.expected, result)
		}
	}
}
