package timewheel

import "time"

var tw *timeWheel = NewTimeWheel(time.Second, 3600)

func init() {
	tw.Start()
}

func At(at time.Time, key string, job func()) {
	tw.AddJob(time.Until(at), key, job)
}
