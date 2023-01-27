package app

import (
	"testing"
	"time"

	"github.com/awcjack/samknows-backend-code-test/types"
)

var app = NewApplication(mockReader{}, mockWriter{})

func TestFindOptimalUnit(t *testing.T) {
	type testcase struct {
		name       string
		input      float64
		result     string
		resultTime int
	}

	testcases := []testcase{
		{
			name:       "Bits per second",
			input:      1,
			result:     "Bits per second",
			resultTime: 0,
		},
		{
			name:       "Kilobits per second",
			input:      1000,
			result:     "Kilobits per second",
			resultTime: 1,
		},
		{
			name:       "Megabits per second",
			input:      1000000,
			result:     "Megabits per second",
			resultTime: 2,
		},
		{
			name:       "Gigabit per second",
			input:      1000000000,
			result:     "Gigabits per second",
			resultTime: 3,
		},
		{
			name:       "Terabits per second",
			input:      1000000000000,
			result:     "Terabits per second",
			resultTime: 4,
		},
		{
			name:       "Petabits per second",
			input:      1000000000000000,
			result:     "Petabits per second",
			resultTime: 5,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			unit, time := app.findOptimalUnit(tc.input)
			if unit != tc.result {
				t.Errorf("Expected get %s, but got %s", tc.result, unit)
			}
			if time != tc.resultTime {
				t.Errorf("Expected get %d, but got %d", tc.resultTime, time)
			}
		})
	}
}

func TestFindMinMaxMean(t *testing.T) {
	type testcase struct {
		name  string
		input []types.Mesurement
		min   float64
		max   float64
		mean  float64
	}

	testcases := []testcase{
		{
			name: "Even element",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 4,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
			},
			min:  1,
			max:  4,
			mean: 2.5,
		},
		{
			name: "Odd element",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
			},
			min:  1,
			max:  3,
			mean: 2,
		},

		{
			name:  "Empty",
			input: []types.Mesurement{},
			min:   0,
			max:   0,
			mean:  0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			min, max, mean := app.findMinMaxMean(tc.input)
			if min != tc.min {
				t.Errorf("Expected get %v, but got %v", tc.min, min)
			}
			if max != tc.max {
				t.Errorf("Expected get %v, but got %v", tc.max, max)
			}
			if mean != tc.mean {
				t.Errorf("Expected get %v, but got %v", tc.mean, mean)
			}
		})
	}
}

func TestFindMedianFirstQuartileIQR(t *testing.T) {
	type testcase struct {
		name          string
		input         []types.Mesurement
		median        float64
		firstQuartile float64
		iqr           float64
	}

	testcases := []testcase{
		{
			name: "Even element",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 4,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
			},
			median:        2.5,
			firstQuartile: 1.5,
			iqr:           2,
		},
		{
			name: "Odd element",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
			},
			median:        2,
			firstQuartile: 1,
			iqr:           2,
		},
		{
			name:          "Empty",
			input:         []types.Mesurement{},
			median:        0,
			firstQuartile: 0,
			iqr:           0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			median, firstQuartile, iqr := app.findMedianFirstQuartileIQR(tc.input)
			if median != tc.median {
				t.Errorf("Expected get %v, but got %v", tc.median, median)
			}
			if firstQuartile != tc.firstQuartile {
				t.Errorf("Expected get %v, but got %v", tc.firstQuartile, firstQuartile)
			}
			if iqr != tc.iqr {
				t.Errorf("Expected get %v, but got %v", tc.iqr, iqr)
			}
		})
	}
}

func TestFindMinMaxDate(t *testing.T) {
	type testcase struct {
		name    string
		input   []types.Mesurement
		minDate time.Time
		maxDate time.Time
	}

	day1, _ := time.Parse("2006-01-02", "2006-01-01")
	day2, _ := time.Parse("2006-01-02", "2006-01-02")
	day3, _ := time.Parse("2006-01-02", "2006-01-03")
	testcases := []testcase{
		{
			name: "Normal",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: day1},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: day3},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: day2},
				},
			},
			minDate: day1,
			maxDate: day3,
		},
		{
			name:    "Empty",
			input:   []types.Mesurement{},
			minDate: time.Time{},
			maxDate: time.Time{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			minDate, maxDate := app.findMinMaxDate(tc.input)
			if !minDate.Equal(tc.minDate) {
				t.Errorf("Expected get %v, but got %v", tc.minDate, minDate)
			}
			if !maxDate.Equal(tc.maxDate) {
				t.Errorf("Expected get %v, but got %v", tc.maxDate, maxDate)
			}
		})
	}
}

func TestFindUnderPerformance(t *testing.T) {
	type testcase struct {
		name          string
		input         []types.Mesurement
		firstQuartile float64
		iqr           float64
		result        []time.Time
	}

	day1, _ := time.Parse("2006-01-02", "2006-01-01")
	day2, _ := time.Parse("2006-01-02", "2006-01-02")
	day3, _ := time.Parse("2006-01-02", "2006-01-03")
	testcases := []testcase{
		{
			name: "Normal",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: day1},
				},
				{
					MetricValue: 2,
					Dtime:       types.JSONTime{Time: day2},
				},
				{
					MetricValue: 3,
					Dtime:       types.JSONTime{Time: day3},
				},
			},
			firstQuartile: 1,
			iqr:           2,
			result:        []time.Time{},
		},
		{
			name: "Outliers",
			input: []types.Mesurement{
				{
					MetricValue: 1,
					Dtime:       types.JSONTime{Time: day1},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: day2},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: day3},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
				{
					MetricValue: 30000,
					Dtime:       types.JSONTime{Time: time.Time{}},
				},
			},
			firstQuartile: 30000,
			iqr:           0,
			result:        []time.Time{day1},
		},
		{
			name:          "Empty",
			input:         []types.Mesurement{},
			firstQuartile: 0,
			iqr:           0,
			result:        []time.Time{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			times := app.findUnderPerformance(tc.input, tc.firstQuartile, tc.iqr)
			if len(times) != len(tc.result) {
				t.Errorf("Expected get %v, but got %v", tc.result, times)
			}
			for i, time := range times {
				if !time.Equal(tc.result[i]) {
					t.Errorf("Expected get %v, but got %v", tc.result, times)
				}
			}
		})
	}
}

func TestDateArrayConcatString(t *testing.T) {
	type testcase struct {
		name   string
		input  []time.Time
		result []string
	}

	day1, _ := time.Parse("2006-01-02", "2006-01-01")
	day2, _ := time.Parse("2006-01-02", "2006-01-02")
	day3, _ := time.Parse("2006-01-02", "2006-01-03")
	day4, _ := time.Parse("2006-01-02", "2006-01-04")
	day5, _ := time.Parse("2006-01-02", "2006-01-05")
	testcases := []testcase{
		{
			name: "Normal",
			input: []time.Time{
				day1,
				day2,
				day3,
				day4,
				day5,
			},
			result: []string{"between 2006-01-01 and 2006-01-05"},
		},
		{
			name: "Normal",
			input: []time.Time{
				day1,
				day2,
				day4,
				day5,
			},
			result: []string{"between 2006-01-01 and 2006-01-02", "between 2006-01-04 and 2006-01-05"},
		},

		{
			name: "Normal1",
			input: []time.Time{
				day1,
				day2,
				day5,
			},
			result: []string{"between 2006-01-01 and 2006-01-02", "2006-01-05"},
		},
		{
			name:   "Empty",
			input:  []time.Time{},
			result: []string{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			dateStrings := app.DateArrayConcatString(tc.input)
			if len(dateStrings) != len(tc.result) {
				t.Errorf("Expected get %v, but got %v", tc.result, dateStrings)
			}
			for i, dateString := range dateStrings {
				if dateString != tc.result[i] {
					t.Errorf("Expected get %v, but got %v", tc.result, dateStrings)
				}
			}
		})
	}
}

type mockReader struct{}

func (r mockReader) GetInputs() ([]types.InputFormat, error) {
	return nil, nil
}

func (r mockReader) GetInput(name string) (types.InputFormat, error) {
	return types.InputFormat{}, nil
}

type mockWriter struct{}

func (w mockWriter) WriteMultipleOutput(outputs []types.OutputFormat) error {
	return nil
}

func (w mockWriter) WriteOutput(name string, content []byte) error {
	return nil
}
