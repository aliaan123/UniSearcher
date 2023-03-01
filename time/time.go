package time

import "time"

var TimeStart = time.Now()

// function that returns the time since TimeStart was initialized
func TimeSinceStart() float64 {
	// returns the duration of the time since the server started
	return time.Since(TimeStart).Seconds()
}
