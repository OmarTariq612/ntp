package ntp

import (
	"time"
)

//go:generate go run golang.org/x/tools/cmd/stringer@latest -output ntp_string.go -type=Version,LeapIndicator,Mode

type Version byte

// Defaults
const (
	V4 Version = 4

	DefaultNTPServer = "time.windows.com:123"
	DefaultPort      = 123
	DefaultTimeout   = 1 * time.Second
)

type LeapIndicator byte

const (
	NoWarning            LeapIndicator = 0
	LastMinOfDayHas61Sec LeapIndicator = 1
	LastMinOfDayHas59Sec LeapIndicator = 2
	ClockUnsynchronized  LeapIndicator = 3
)

type Mode byte

const (
	Reserved          Mode = 0
	SymmetricActive   Mode = 1
	SymmetricPassive  Mode = 2
	Client            Mode = 3
	Server            Mode = 4
	Broadcast         Mode = 5
	NTPControlMessage Mode = 6
	ReservedPrivate   Mode = 7
)

func LIVNModeToByte(li LeapIndicator, v Version, m Mode) uint8 {
	return (uint8(li) << 6) | (uint8(v) << 3) | uint8(m)
}

func LIVNModeFromByte(b uint8) (li LeapIndicator, v Version, m Mode) {
	li = LeapIndicator((b >> 6) & 0x03)
	v = Version((b >> 3) & 0x07)
	m = Mode(b & 0x07)
	return li, v, m
}

// The NTP packet header has 12 words followed by optional
// extension fields and finally an optional message authentication code (MAC)
// consisting of the Key Identifier field and Message Digest field.
type Header struct {
	// leap indicator (2 bits) - Version (3 bits) - Mode (3 bits)
	LIVNMode uint8

	// Stratum (stratum): 8-bit integer representing the stratum, with values:
	//
	// +--------+-----------------------------------------------------+
	// | Value  | Meaning                                             |
	// +--------+-----------------------------------------------------+
	// | 0      | unspecified or invalid                              |
	// | 1      | primary server (e.g., equipped with a GPS receiver) |
	// | 2-15   | secondary server (via NTP)                          |
	// | 16     | unsynchronized                                      |
	// | 17-255 | reserved                                            |
	// +--------+-----------------------------------------------------+
	Stratum uint8

	// Poll: 8-bit signed integer representing the maximum interval between
	// successive messages, in log2 seconds.
	//
	// Suggested default limits for minimum and maximum poll intervals are 6 and 10, respectively.
	Poll int8

	// Precision: 8-bit signed integer representing the precision of the
	// system clock, in log2 seconds.
	//
	// For instance, a value of -18 corresponds to a precision of about one microsecond.
	// The precision can be determined when the service first starts up as the minimum time
	// of several iterations to read the system clock.
	Precision int8

	// Root Delay (rootdelay): Total round-trip delay to the reference
	// clock, in NTP short format.
	RootDelay ShortFormat

	// Root Dispersion (rootdisp): Total dispersion to the reference clock,
	// in NTP short format.
	RootDispersion ShortFormat

	// Reference ID (refid): 32-bit code identifying the particular server or reference clock.
	//
	// The interpretation depends on the value in the stratum field.
	ReferenceID uint32

	// Reference Timestamp: Time when the system clock was last set or corrected,
	// in NTP timestamp format.
	ReferenceTimestamp TimestampFormat

	// Origin Timestamp (org): Time at the client when the request departed
	// for the server, in NTP timestamp format.
	OriginTimestamp TimestampFormat

	// Receive Timestamp (rec): Time at the server when the request arrived
	// from the client, in NTP timestamp format.
	ReceiveTimestamp TimestampFormat

	// Transmit Timestamp (xmt): Time at the server when the response left
	// for the client, in NTP timestamp format.
	TransmitTimestamp TimestampFormat
}

func (h *Header) LeapIndicator() LeapIndicator {
	return LeapIndicator((h.LIVNMode >> 6) & 0x03)
}

func (h *Header) Version() Version {
	return Version((h.LIVNMode >> 3) & 0x07)
}

func (h *Header) Mode() Mode {
	return Mode(h.LIVNMode & 0x07)
}

func Offset(org, rec, xmt, dst TimestampFormat) time.Duration {
	// ((rec - org) + (xmt - dst)) / 2
	a := rec.Time().Sub(org.Time())
	b := xmt.Time().Sub(dst.Time())
	return (a + b) / time.Duration(2)
}

func Delta(org, rec, xmt, dst TimestampFormat) time.Duration {
	// (dst - org) - (xmt - rec)
	a := dst.Time().Sub(org.Time())
	b := xmt.Time().Sub(rec.Time())
	return a - b
}
