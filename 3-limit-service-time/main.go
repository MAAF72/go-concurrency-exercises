//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync/atomic"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// Freetime define free time in second
const Freetime int64 = 10

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	var start time.Time
	var end time.Time
	var status bool

	chanProcess := make(chan bool, 1)
	go func() {
		start = time.Now()
		process()
		chanProcess <- true
	}()

	select {
	case status = <-chanProcess:
		status = true
	case <-time.After(time.Duration(Freetime-atomic.LoadInt64(&u.TimeUsed)) * time.Second):
		if u.IsPremium {
			// continue the process if the user are a paid premium user.
			status = <-chanProcess
		} else {
			status = false
		}
	}

	end = time.Now()

	atomic.AddInt64(&u.TimeUsed, (end.Unix() - start.Unix()))

	return status
}

func main() {
	RunMockServer()
}
