package app

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/awcjack/samknows-backend-code-test/infrastructure/reader"
	"github.com/awcjack/samknows-backend-code-test/infrastructure/writer"
	"github.com/awcjack/samknows-backend-code-test/types"
)

type Application struct {
	reader reader.Reader
	writer writer.Writer
}

// function to make new application with reader and writer (using interfae to provide flexibility to switch to other reader or writer like database easily)
func NewApplication(reader reader.Reader, writer writer.Writer) Application {
	return Application{
		reader: reader,
		writer: writer,
	}
}

// Always use ____bits per second unit to prevent confussion
func (a Application) findOptimalUnit(min float64) (string, int) {
	// bytes per second to bits per second
	min = min * 8
	result := "Bits per second"
	time := 0

	if min > 1000 {
		min /= 1000
		result = "Kilobits per second"
		time = 1
	}

	if min > 1000 {
		min /= 1000
		result = "Megabits per second"
		time = 2
	}

	if min > 1000 {
		min /= 1000
		result = "Gigabits per second"
		time = 3
	}

	if min > 1000 {
		min /= 1000
		result = "Terabits per second"
		time = 4
	}

	if min > 1000 {
		min /= 1000
		result = "Petabits per second"
		time = 5
	}

	return result, time
}

// function to find min, max, average from dataset
func (a Application) findMinMaxMean(input []types.Mesurement) (float64, float64, float64) {
	var min float64 = 0
	var max float64 = 0
	var sum float64 = 0

	if len(input) == 0 {
		return 0, 0, 0
	}

	for _, mesurement := range input {
		if mesurement.MetricValue < min || min == 0 {
			min = mesurement.MetricValue
		}

		if mesurement.MetricValue > max {
			max = mesurement.MetricValue
		}

		sum += mesurement.MetricValue
	}

	return min, max, (sum / float64(len(input)))
}

// function to find median, first quartile and IQR from dataset
func (a Application) findMedianFirstQuartileIQR(input []types.Mesurement) (float64, float64, float64) {
	floatArray := make([]float64, 0, len(input))

	for _, mesurement := range input {
		floatArray = append(floatArray, mesurement.MetricValue)
	}

	sort.Float64s(floatArray)

	var median float64
	var firstQuartile float64
	var thirdQuartile float64
	l := len(floatArray)
	if l == 0 {
		return 0, 0, 0
	} else if l%2 == 0 {
		median = (floatArray[l/2-1] + floatArray[l/2]) / 2
		firstQuartile = (floatArray[l/4-1] + floatArray[l/4]) / 2
		thirdQuartile = (floatArray[3*l/4-1] + floatArray[3*l/4]) / 2
	} else {
		median = floatArray[l/2]
		firstQuartile = floatArray[l/4]
		thirdQuartile = floatArray[3*l/4]
	}

	return median, firstQuartile, thirdQuartile - firstQuartile
}

// function to find min date and max date from data set (order may not preserved in production)
func (a Application) findMinMaxDate(input []types.Mesurement) (time.Time, time.Time) {
	if len(input) == 0 {
		return time.Time{}, time.Time{}
	}

	timeArray := make([]time.Time, 0, len(input))

	for _, mesurement := range input {
		timeArray = append(timeArray, mesurement.Dtime.Time)
	}

	sort.Slice(timeArray, func(i, j int) bool {
		return timeArray[i].Before(timeArray[j])
	})

	return timeArray[0], timeArray[len(timeArray)-1]
}

// function to find period that under performance (based on Q1 - 1.5 * IQR rules to find outlier)
func (a Application) findUnderPerformance(input []types.Mesurement, firstQuartile float64, IQR float64) []time.Time {
	result := make([]time.Time, 0)

	for _, mesurement := range input {
		if mesurement.MetricValue < firstQuartile-1.5*IQR {
			result = append(result, mesurement.Dtime.Time)
		}
	}

	return result
}

// function to convert time slice to string slice that concat continuous date into period
func (a Application) DateArrayConcatString(times []time.Time) []string {
	if len(times) == 0 {
		return nil
	}

	if len(times) == 1 {
		return []string{times[0].Format("2006-01-02")}
	}

	result := make([]string, 0)

	startCursor := times[0]
	cursor := times[0]
	for i := 1; i < len(times); i++ {
		nextCursor := times[i]

		fmt.Println("cursor", cursor)
		fmt.Println("nextCursor", nextCursor)
		fmt.Println("!cursor.AddDate(0, 0, 1).Equal(nextCursor)", !cursor.AddDate(0, 0, 1).Equal(nextCursor))
		fmt.Println("i == len(times)-1", i == len(times)-1)
		if i == len(times)-1 {
			if startCursor.Equal(nextCursor) {
				result = append(result, nextCursor.Format("2006-01-02"))
			} else if !cursor.AddDate(0, 0, 1).Equal(nextCursor) {
				result = append(result, fmt.Sprintf("between %s and %s", startCursor.Format("2006-01-02"), cursor.Format("2006-01-02")))
				result = append(result, nextCursor.Format("2006-01-02"))
			} else {
				result = append(result, fmt.Sprintf("between %s and %s", startCursor.Format("2006-01-02"), nextCursor.Format("2006-01-02")))
			}
		} else if !cursor.AddDate(0, 0, 1).Equal(nextCursor) {
			result = append(result, fmt.Sprintf("between %s and %s", startCursor.Format("2006-01-02"), cursor.Format("2006-01-02")))
			startCursor = nextCursor
		}

		cursor = nextCursor
	}

	return result
}

// function to run the pull data, process and report data
func (a Application) Run() error {
	inputArray, err := a.reader.GetInputs()
	if err != nil {
		return err
	}

	for _, input := range inputArray {
		min, max, mean := a.findMinMaxMean(input.Content)
		median, firstQuartile, IQR := a.findMedianFirstQuartileIQR(input.Content)
		underPerformancePeriod := a.findUnderPerformance(input.Content, firstQuartile, IQR)
		minDate, maxDate := a.findMinMaxDate(input.Content)

		minValue := math.Min(min, math.Min(max, math.Min(median, mean)))
		unit, time := a.findOptimalUnit(minValue)

		fileName := strings.Split(input.Name, ".")
		var output string
		if len(underPerformancePeriod) > 0 {
			output = fmt.Sprintf(`SamKnows Metric Analyser v1.0.0
===============================

Period checked:

    From: %s
    To:   %s

Statistics:

    Unit: %s

    Average: %.2f
    Min: %.2f
    Max: %.2f
    Median: %.2f

Under-performing periods:

    * The period %s
      was under-performing.
`, minDate.Format("2006-01-02"), maxDate.Format("2006-01-02"), unit, mean*8/math.Pow(1000, float64(time)), min*8/math.Pow(1000, float64(time)), max*8/math.Pow(1000, float64(time)), median*8/math.Pow(1000, float64(time)), strings.Join(a.DateArrayConcatString(underPerformancePeriod), ", "))
		} else {
			output = fmt.Sprintf(`SamKnows Metric Analyser v1.0.0
===============================

Period checked:

    From: %s
    To:   %s

Statistics:

    Unit: %s

    Average: %.2f
    Min: %.2f
    Max: %.2f
    Median: %.2f
`, minDate.Format("2006-01-02"), maxDate.Format("2006-01-02"), unit, mean*8/math.Pow(1000, float64(time)), min*8/math.Pow(1000, float64(time)), max*8/math.Pow(1000, float64(time)), median*8/math.Pow(1000, float64(time)))
		}

		a.writer.WriteOutput(fileName[0]+".output", []byte(output))
	}

	return nil
}
