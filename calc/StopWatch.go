package calc

import "time"

type StopWatch struct {
	records []time.Duration

	startTime time.Time
}

const maxRecords = 10

//Start Begins the Stopwatch timer
func (s *StopWatch) Start() {
	s.startTime = time.Now()
}

//Stop Stops the Stopwatch Timer
func (s *StopWatch) Stop() {

	now := time.Now()
	diff := now.Sub(s.startTime)

	if len(s.records) == maxRecords {
		s.records = s.records[1:len(s.records)]
	}

	s.records = append(s.records, diff)
}

func (s *StopWatch) GetCurrentElaspedTime() time.Duration {
	now := time.Now()
	diff := now.Sub(s.startTime)
	return diff
}

func (s *StopWatch) IsStarted() bool {
	return s.startTime != time.Time{}
}

//GetAverage Returns the Average time of all records within the timer
func (s *StopWatch) GetAverage() time.Duration {

	var total time.Duration

	for i := 0; i < len(s.records); i++ {
		total += s.records[i]
	}

	var div = time.Duration(len(s.records))

	if div == 0 {
		div = 1
	}

	return total / div

}
