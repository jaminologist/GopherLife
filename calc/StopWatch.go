package calc

import "time"

type StopWatch struct {
	records []int64

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

	s.records = append(s.records, diff.Nanoseconds())
}

//GetAverage Returns the Average time of all records within the timer
func (s *StopWatch) GetAverage() int64 {

	var total int64

	for i := 0; i < len(s.records); i++ {
		total += s.records[i]
	}

	return total / int64(len(s.records))

}
