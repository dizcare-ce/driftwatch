// Package throttle implements a sliding-window token bucket used to cap
// the rate at which drift-check runs may be triggered.
//
// # Usage
//
//	th := throttle.New(time.Minute, 10) // at most 10 runs per minute
//
//	if !th.Allow() {
//		log.Println("throttled: too many runs in the current window")
//		return
//	}
//	// proceed with the drift check run
//
// Remaining returns the number of calls still permitted in the current
// window. Reset clears all recorded timestamps, useful in tests or after
// a configuration reload.
package throttle
