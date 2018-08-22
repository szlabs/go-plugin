package plugin

import (
	"sync"

	"github.com/steven-zou/go-plugin/pkg/spec"
)

//Store defines how to maintain the loaded plugin
type Store interface {
	//The total count of current items in the store
	Size() uint

	//Append the plugin item to the store
	//If forced is set to be true, new plugin item will overwrite the existing one
	//Try best to append, ignore any errors
	Put(item *spec.PluginItem, forced bool)

	//Get the plugin item by name
	//If existing, return the item and set the bool flag to true
	Get(name string) (*spec.PluginItem, bool)

	//Remove the plugin out of the store and return the removed plugin item
	//If successfully removed, set the bool flag to true
	Remove(name string) (*spec.PluginItem, bool)
}

//BaseStore is the default implementation of Store interface
type BaseStore struct {
	//internal lock
	lock *sync.RWMutex

	//internal list
	hash map[string]*spec.PluginItem
}

//NewBaseStore is constructor of BaseStore
func NewBaseStore() *BaseStore {
	return &BaseStore{
		lock: new(sync.RWMutex),
		hash: make(map[string]*spec.PluginItem),
	}
}

//Size is the implementation of same method in Store interface
func (bs *BaseStore) Size() uint {
	bs.lock.Lock()
	defer bs.lock.Unlock()

	return (uint)(len(bs.hash))
}

//Put is the implementation of same method in Store interface
func (bs *BaseStore) Put(item *spec.PluginItem, forced bool) {
	if item == nil || item.Spec == nil || item.Executor == nil {
		return
	}

	name := item.Spec.Name
	if len(name) == 0 {
		return
	}

	_, ok := bs.hash[name]
	if !ok || (ok && forced) {
		bs.lock.Lock()
		defer bs.lock.Unlock()

		bs.hash[name] = item
	}
}

//Get is the implementation of same method in Store interface
func (bs *BaseStore) Get(name string) (*spec.PluginItem, bool) {
	bs.lock.RLock()
	defer bs.lock.RUnlock()

	item, ok := bs.hash[name]

	return item, ok
}

//Remove is the implementation of same method in Store interface
func (bs *BaseStore) Remove(name string) (*spec.PluginItem, bool) {
	bs.lock.Lock()
	defer bs.lock.Unlock()

	item, ok := bs.hash[name]
	if ok {
		delete(bs.hash, name)
		return item, ok
	}

	return nil, false
}
