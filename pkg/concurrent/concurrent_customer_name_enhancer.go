package concurrent

import (
	. "github.com/talbx/go-clockodo/pkg/cache"
	. "github.com/talbx/go-clockodo/pkg/model"
)

type ConcurrentCustomerNameEnhancer struct {
	TimeEntries map[string][]DayByCustomer
}

func (c ConcurrentCustomerNameEnhancer) Aggregate() *map[string][]DayByCustomer {
	kache := CreateCache()
	veryNewMap := make(map[string][]DayByCustomer)

	// for every weekday in that map
	for k := range c.TimeEntries {
		channel := make(chan DayByCustomer)

		// go over every task entry of a weekday
		for _, customerDay := range c.TimeEntries[k] {
			// and fetch its customerName, publish it to a weekday-group channel
			go getAndAddCustomerNamesToEntries(customerDay, &kache, channel)
		}

		// collect all updated task entries for each week day
		resultList := make([]DayByCustomer, 0)
		for i := 0; i < len(c.TimeEntries[k]); i++ {
			resultList = append(resultList, <-channel)
		}
		veryNewMap[k] = resultList
	}
	return &veryNewMap
}

func getAndAddCustomerNamesToEntries(v DayByCustomer, cache *CustomerNameCache, c1 chan<- DayByCustomer) {
	v.Customer = cache.Get(v.CustomerId)
	c1 <- v
}
