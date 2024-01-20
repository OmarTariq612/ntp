package ntp

import "time"

const (
	t10_3 = 1_000
	t10_6 = 1_000_000
	t10_9 = 1_000_000_000
)

var (
	ntpEpoch  = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	unixEpoch = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	// time.Duration from the NTP epoch (1 / 1 / 1900) to the Unix epoch (1 / 1 / 1970)
	fromNTPtoUnixEpoch = unixEpoch.Sub(ntpEpoch)
)

type NTPTimeFormat interface {
	// time.Duration since the NTP epoch
	Duration() time.Duration

	Time() time.Time
}

var (
	_ NTPTimeFormat = TimestampFormat(0)
	_ NTPTimeFormat = ShortFormat(0)
)

type TimestampFormat uint64

func (ts TimestampFormat) Seconds() uint32 {
	return uint32(ts >> 32)
}

func (ts TimestampFormat) Fraction() uint32 {
	return uint32(ts & 0xffffffff)
}

func (ts TimestampFormat) Duration() time.Duration {
	secs := uint64(ts>>32) * t10_9
	fraction := uint64(ts&0xffffffff) * t10_9
	nsecs := fraction >> 32
	if fraction >= 0x80000000 {
		nsecs++
	}
	return time.Duration(secs + nsecs)
}

func (ts TimestampFormat) Time() time.Time {
	return ntpEpoch.Add(ts.Duration())
}

func (ts TimestampFormat) UnixNano() int64 {
	fromNTPDuration := int64(ts.Duration())
	return fromNTPDuration - int64(fromNTPtoUnixEpoch)
}

func (ts TimestampFormat) UnixMicro() int64 {
	return ts.UnixNano() / t10_3
}

func (ts TimestampFormat) UnixMilli() int64 {
	return ts.UnixNano() / t10_6
}

func (ts TimestampFormat) Unix() int64 {
	return ts.UnixNano() / t10_9
}

func TimestampFormatFromUnixNano(nano int64) TimestampFormat {
	nano += int64(fromNTPtoUnixEpoch)
	secs := uint64(nano / t10_9)
	nsecs := uint64(uint64(nano)-secs*t10_9) << 32
	fraction := uint64(nsecs / t10_9)
	if nsecs%t10_9 >= t10_9/2 {
		fraction++
	}
	return TimestampFormat(secs<<32 | fraction)
}

type ShortFormat uint32

func (s ShortFormat) Seconds() uint16 {
	return uint16(s >> 16)
}

func (s ShortFormat) Fraction() uint16 {
	return uint16(s & 0xffff)
}

func (s ShortFormat) Duration() time.Duration {
	secs := uint64(s>>16) * t10_9
	fraction := uint64(s&0xffff) * t10_9
	nsecs := fraction >> 16
	if fraction >= 0x8000 {
		nsecs++
	}
	return time.Duration(secs + nsecs)
}

func (s ShortFormat) Time() time.Time {
	return ntpEpoch.Add(s.Duration())
}

func (s ShortFormat) UnixNano() int64 {
	fromNTPDuration := int64(s.Duration())
	return fromNTPDuration - int64(fromNTPtoUnixEpoch)
}

func (s ShortFormat) UnixMicro() int64 {
	return s.UnixNano() / t10_3
}

func (s ShortFormat) UnixMilli() int64 {
	return s.UnixNano() / t10_6
}

func (s ShortFormat) Unix() int64 {
	return s.UnixNano() / t10_9
}

func ShortFormatFromUnixNano(nano int64) ShortFormat {
	nano += int64(fromNTPtoUnixEpoch)
	secs := uint64(nano / t10_9)
	nsecs := uint64(uint64(nano)-secs*t10_9) << 16
	fraction := uint64(nsecs / t10_9)
	if nsecs%t10_9 >= t10_9/2 {
		fraction++
	}
	return ShortFormat(secs<<16 | fraction)
}
