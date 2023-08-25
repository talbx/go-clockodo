package concurrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/talbx/go-clockodo/pkg/model"
	"testing"
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

	fmt.Println(result)
	assert.Len(t, result, 2)
	var onlySundayAndMondayKeysExist assert.Comparison = func() bool {

		if nil != result["01.01.2023"] && nil != result["02.01.2023"] && nil == result["03.01.2023"] {
			return true
		}
		return false
	}
	assert.Condition(t, onlySundayAndMondayKeysExist)
	assert.Len(t, result["01.01.2023"], 2)
	assert.Len(t, result["02.01.2023"], 1)
}
