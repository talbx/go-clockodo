package util

import (
	"fmt"
	"sync"
)

type Customer struct {
	Cus `json:"customer"`
}

type Cus struct {
	Name string `json:"name"`
}

func getCustomerName(id int) string {
	var customer Customer
	cal := fmt.Sprintf("v2/customers/%v", id)
	CallApi(cal, &customer)
	return customer.Name
}

type Cache interface {
	Get(id int) string
}

var customerDb sync.Map

type CustomerNameCache struct{}

func (cache CustomerNameCache) Get(id int) string {
	customerName, _ := customerDb.LoadOrStore(id, getCustomerName(id))
	return customerName.(string)
}
