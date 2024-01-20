package ntp

import "testing"

func TestTimestampFormatSecondsAndFraction(t *testing.T) {
	testcases := []struct {
		timestampFormat TimestampFormat
		secs            uint32
		fraction        uint32
	}{
		{timestampFormat: TimestampFormat(0x1111111105050505), secs: 0x11111111, fraction: 0x05050505},
		{timestampFormat: TimestampFormat(0xa0a0a0a0ffffffff), secs: 0xa0a0a0a0, fraction: 0xffffffff},
		{timestampFormat: TimestampFormat(0x0123456789abcdef), secs: 0x01234567, fraction: 0x89abcdef},
	}

	for _, test := range testcases {
		// testing the fraction first to make sure that seconds aren't affecting the results
		if fraction := test.timestampFormat.Fraction(); fraction != test.fraction {
			t.Errorf("0x%X != 0x%X", fraction, test.fraction)
		}
		if secs := test.timestampFormat.Seconds(); secs != test.secs {
			t.Errorf("0x%X != 0x%X", secs, test.secs)
		}
	}
}

func TestShortFormatSecondsAndFraction(t *testing.T) {
	testCases := []struct {
		shortFormat ShortFormat
		secs        uint16
		fraction    uint16
	}{
		{shortFormat: ShortFormat(0x11110505), secs: 0x1111, fraction: 0x0505},
		{shortFormat: ShortFormat(0xa0a0ffff), secs: 0xa0a0, fraction: 0xffff},
		{shortFormat: ShortFormat(0x01234567), secs: 0x0123, fraction: 0x4567},
	}

	for _, test := range testCases {
		// testing the fraction first to make sure that seconds aren't affecting the results
		if fraction := test.shortFormat.Fraction(); fraction != test.fraction {
			t.Errorf("0x%X != 0x%X", fraction, test.fraction)
		}
		if secs := test.shortFormat.Seconds(); secs != test.secs {
			t.Errorf("0x%X != 0x%X", secs, test.secs)
		}
	}
}
