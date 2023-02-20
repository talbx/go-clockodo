package util

import (
	"fmt"
	"log"
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

var customerDb map[int]string = make(map[int]string)

type CustomerNameCache struct {}

func (cache CustomerNameCache) Get(id int) string {
	log.Printf("Looking for customer with id %v\n", id)
	for k, v := range customerDb {
		if k == id {
			log.Printf("Found customer %v with id %v in cache!\n",v, id)
			return v
		}
	}
	log.Printf("No customer found for id %v; will contact api!\n", id)
	customerName := getCustomerName(id)
	customerDb[id] = customerName
	log.Printf("customer retrieved from api and stored in cache: %v - %v\n",id, customerName)
	return customerName
}