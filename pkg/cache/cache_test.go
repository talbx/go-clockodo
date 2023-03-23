package cache

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

type MockCache struct {
	interceptFn func(id int) string
}

func (c MockCache) Get(id int) string {
	return c.interceptFn(id)
}

func TestCustomerNameCache_Get(t *testing.T) {

	// given
	someId := 1234
	custName := "peter"
	var cache Cache = MockCache{
		interceptFn: func(id int) string {
			return custName
		},
	}

	// when
	customerName := cache.Get(someId)

	// then
	assert.Equal(t, customerName, custName)

}
