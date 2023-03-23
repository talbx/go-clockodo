package concurrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/talbx/go-clockodo/pkg/model"
	"github.com/talbx/go-clockodo/pkg/util"
	"testing"
	"time"
)

func TestConcurrentDayEntriesByCustomerAggregator_Aggregate(t *testing.T) {
	util.CreateSugaredLogger()
	timeEntry1 := model.TimeEntry{
		// monday
		StartTime:  "2023-01-02T15:04:05Z",
		CustomerId: 1234,
	}
	timeEntry2 := model.TimeEntry{
		// sunday
		StartTime:  "2023-01-01T15:04:05Z",
		CustomerId: 1234,
	}
	timeEntry3 := model.TimeEntry{
		// sunday
		StartTime:  "2023-01-01T16:04:05Z",
		CustomerId: 2345,
	}

	custIds := []int{1234, 2345}

	karte := make(map[time.Weekday][]model.TimeEntry)
	karte[time.Monday] = append(karte[time.Monday], timeEntry1)
	karte[time.Sunday] = append(karte[time.Sunday], timeEntry2)
	karte[time.Sunday] = append(karte[time.Sunday], timeEntry3)

	aggregator := ConcurrentDayEntriesByCustomerAggregator{&karte, custIds}
	result := *aggregator.Aggregate()

	fmt.Println(result)
	assert.Len(t, result, 2)
	assert.Len(t, result[time.Monday], 1)
	assert.Len(t, result[time.Sunday], 2)
	fmt.Println(result)
}
