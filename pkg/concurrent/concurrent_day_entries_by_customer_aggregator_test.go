package concurrent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/talbx/go-clockodo/pkg/model"
	"testing"
)

func TestConcurrentDayEntriesByCustomerAggregator_Aggregate(t *testing.T) {
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

	karte := make(map[string][]model.TimeEntry)
	karte["02.01.2023"] = append(karte["02.01.2023"], timeEntry1)
	karte["01.01.2023"] = append(karte["01.01.2023"], timeEntry2)
	karte["01.01.2023"] = append(karte["01.01.2023"], timeEntry3)

	aggregator := ConcurrentDayEntriesByCustomerAggregator{&karte, custIds}
	result := *aggregator.Aggregate()

	fmt.Println(result)
	assert.Len(t, result, 2)
	assert.Len(t, result["02.01.2023"], 1)
	assert.Len(t, result["01.01.2023"], 2)
	fmt.Println(result)
}
