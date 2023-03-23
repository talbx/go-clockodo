package cache

import (
	"fmt"
	"github.com/talbx/go-clockodo/pkg/model"
	"github.com/talbx/go-clockodo/pkg/util"
	"sync"
)

type Cache interface {
	Get(id int) string
}

type CustomerNameCache struct {
	customerDb *sync.Map
	mutex      sync.Mutex
}

func CreateCache() CustomerNameCache {
	return CustomerNameCache{
		customerDb: &sync.Map{},
		mutex:      sync.Mutex{},
	}
}

func (cache *CustomerNameCache) Get(id int) string {

	defer cache.mutex.Unlock()
	cache.mutex.Lock()
	util.SugaredLogger.Infof("[Cache] looking for customerId %v", id)
	if name, ok := cache.customerDb.Load(id); ok {
		util.SugaredLogger.Infof("[Cache] found customerId %v in cache: %v", id, name)
		return name.(string)
	}
	util.SugaredLogger.Infof("[Cache] never saw customerId %v before; will ask API", id)
	name := cache.getCustomerName(id)
	cache.customerDb.Store(id, name)
	return name
}

func (cache *CustomerNameCache) getCustomerName(id int) string {
	var customer model.Customer
	cal := fmt.Sprintf("v2/customers/%v", id)
	util.CallApi(cal, &customer)
	return customer.Name
}
