package types

import "time"

// type for decoding date from JSON to time.Time
type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}

	t.Time = date
	return nil
}

type Mesurement struct {
	MetricValue float64  `json:"metricValue"`
	Dtime       JSONTime `json:"dtime"`
}

type InputFormat struct {
	Name    string
	Content []Mesurement
}

type OutputFormat struct {
	Name    string
	Content []byte
}
