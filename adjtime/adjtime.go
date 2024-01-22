package adjtime

import "time"

func AdjTime(delta time.Duration) error {
	return adjtime(delta)
}
