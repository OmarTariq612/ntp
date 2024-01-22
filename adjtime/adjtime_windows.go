package adjtime

import (
	"time"

	"golang.org/x/sys/windows"
)

const (
	FileTime1970 = 0x019db1ded53e8000
)

func adjtime(delta time.Duration) (err error) {
	// TODO: add support for windows
	// Ref: https://github.com/ntp-project/ntp/blob/9c75327c3796ff59ac648478cd4da8b205bceb77/ports/winnt/libntp/SetSystemTime.c#L29
	var currTime windows.FileTime
	windows.GetSystemTimeAsFileTime(&currTime)
	return nil
}
