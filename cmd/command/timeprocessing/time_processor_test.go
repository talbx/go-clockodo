package timeprocessing

import (
	"github.com/magiconair/properties/assert"
	"github.com/talbx/go-clockodo/pkg/model"
	"strconv"
	"testing"
	"time"
)

func TestAddLeadingZero(t *testing.T) {
	// given
	num := 2
	bignum := 12

	// when
	numWithLeadingZero := AddLeadingZero(num)
	bignumWithoutLeadingZero := AddLeadingZero(bignum)

	// then
	assert.Equal(t, numWithLeadingZero, "02")
	assert.Equal(t, bignumWithoutLeadingZero, strconv.Itoa(bignum))
}

func TestDurationToHM(t *testing.T) {
	// given
	duration := 30000

	// when
	hours, minutes := DurationToHM(duration)

	assert.Equal(t, hours, 8)
	assert.Equal(t, minutes, 20)
}

func TestEOB(t *testing.T) {
	// given
	now := time.Now()

	// when
	eob := EOB(now)

	// then
	assert.Equal(t, eob.Hour(), 23)
	assert.Equal(t, eob.Minute(), 59)
	assert.Equal(t, eob.Day(), now.Day())
}

func TestProcessClock(t *testing.T) {
	type args struct {
		clock *model.ClockResponse
	}
	var tests []struct {
		name        string
		args        args
		wantHours   int
		wantMinutes int
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHours, gotMinutes := ProcessClock(tt.args.clock)
			if gotHours != tt.wantHours {
				t.Errorf("ProcessClock() gotHours = %v, want %v", gotHours, tt.wantHours)
			}
			if gotMinutes != tt.wantMinutes {
				t.Errorf("ProcessClock() gotMinutes = %v, want %v", gotMinutes, tt.wantMinutes)
			}
		})
	}
}

func TestRound(t *testing.T) {
	// given
	testingTable := [][]int{
		{23, 30, 23, 30},
		{23, 37, 23, 30},
		{23, 38, 23, 45},
		{23, 45, 23, 45},
		{23, 52, 23, 45},
		{23, 53, 24, 0},
		{8, 29, 8, 30},
	}

	for _, tt := range testingTable {
		h, m := Round(tt[0], tt[1])
		assert.Equal(t, h, tt[2], " - Hours do not match")
		assert.Equal(t, m, tt[3], " - Minutes do not match")
	}
}

func TestSOB(t *testing.T) {
	// given
	now := time.Now()

	// when
	eob := SOB(now)

	// then
	assert.Equal(t, eob.Hour(), 0)
	assert.Equal(t, eob.Minute(), 0)
	assert.Equal(t, eob.Day(), now.Day())
}

func Test_d(t *testing.T) {
	// given
	now := time.Now()

	// when
	day := d(now, 15, 15)

	// then
	assert.Equal(t, day.Hour(), 15)
	assert.Equal(t, day.Minute(), 15)
}
