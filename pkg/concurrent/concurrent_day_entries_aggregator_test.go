package concurrent

import (
	"github.com/stretchr/testify/assert"
	"github.com/talbx/go-clockodo/pkg/model"
	"testing"
	"time"
)

func TestGroupEntriesByDay(t *testing.T) {

	timeEntry1 := model.TimeEntry{
		// monday
		StartTime: "2023-01-02T15:04:05Z",
	}
	timeEntry2 := model.TimeEntry{
		// sunday
		StartTime: "2023-01-01T15:04:05Z",
	}
	timeEntry3 := model.TimeEntry{
		// sunday
		StartTime: "2023-01-01T16:04:05Z",
	}
	input := model.TimeEntriesResponse{Entries: []model.TimeEntry{timeEntry1, timeEntry2, timeEntry3}}

	aggregator := ConcurrentDayEntriesAggregator{
		input,
	}
	result := *aggregator.Aggregate()

	assert.Len(t, result, 2)
	var onlySundayAndMondayKeysExist assert.Comparison = func() bool {
		if nil != result[time.Sunday] && nil != result[time.Monday] && nil == result[time.Tuesday] {
			return true
		}
		return false
	}
	assert.Condition(t, onlySundayAndMondayKeysExist)
	assert.Len(t, result[time.Sunday], 2)
	assert.Len(t, result[time.Monday], 1)
}
