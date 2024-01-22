package adjtime

import (
	"time"

	"github.com/OmarTariq612/ntp/proto"
	"golang.org/x/sys/unix"
)

func adjtime(delta time.Duration) (err error) {
	timeVal := unix.Timeval{Sec: int64(delta) / proto.T10_9, Usec: int64(delta) % proto.T10_9}
	for timeVal.Usec < 0 {
		timeVal.Sec -= 1
		timeVal.Usec += proto.T10_9
	}
	_, err = unix.Adjtimex(&unix.Timex{
		Modes: unix.ADJ_SETOFFSET | unix.ADJ_NANO,
		Time:  timeVal,
	})
	return err
}
