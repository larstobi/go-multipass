package multipass

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	testCases := []struct {
		prevVersion      string
		mostRecentVersion string
		expectedResult   int
	}{
		{"v1.10.0-rc", "v1.11.0-rc", -1}, // prevVersion is less than mostRecentVersion
		{"v1.11.0-rc", "v1.11.0-rc", 0},  // prevVersion is equal to mostRecentVersion
		{"v1.11.0-rc", "v1.10.0-rc", 1},  // prevVersion is greater than mostRecentVersion
		{"0,5", "v1.11.0-rc", -1}, // prevVersion is less than mostRecentVersion
		{"2018.12.1", "v1.11.0-rc", -1}, // prevVersion is less than mostRecentVersion
	}

	for _, tc := range testCases {
		got := CompareVersion(tc.prevVersion, tc.mostRecentVersion)
		if got != tc.expectedResult {
			t.Errorf("Comparison failed for versions %s and %s. Expected: %d, Got: %d",
				tc.prevVersion, tc.mostRecentVersion, tc.expectedResult, got)
		}
	}

}