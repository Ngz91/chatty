package services

import (
	"fmt"
	"net"
	"sync"
)

type Room struct {
	Id    string
	Users []net.Conn
}

type DataStore struct {
	sync.Mutex
	datastore map[string]Room
}

func newDataStore() DataStore {
	return DataStore{
		sync.Mutex{},
		make(map[string]Room),
	}
}

/*
Puts the key and value in the DataStore.
If the key (k) already exists will replace with the provided value (v).
*/
func (ds *DataStore) put(k string, conn net.Conn) {
	ds.Lock()
	room, ok := ds.datastore[k]
	if !ok {
		room = Room{Id: k}
		ds.datastore[k] = room
	}
	room.Users = append(room.Users, conn)

	ds.Unlock()
}

/*
Returns true, if the DataStore has the key (k)
*/
func (ds *DataStore) containsKey(k string) bool {
	ds.Lock()
	if _, ok := ds.datastore[k]; ok {
		return ok
	}
	ds.Unlock()
	return false
}

/*
Puts the key and value in the DataStore, ONLY if the key (k) is not present.
Returns true, if the put operation is successful,
false if the key (k) ia already present in the DataStore
*/
func (ds *DataStore) putIfAbsent(k string, conn net.Conn) bool {
	ds.Lock()
	if _, ok := ds.datastore[k]; ok {
		fmt.Println("datastore contains key: ", k)
		ds.Unlock()
		return false
	}

	room := Room{Id: k, Users: []net.Conn{conn}}
	ds.datastore[k] = room
	ds.Unlock()
	return true
}

/*
Returns the Data value associated with the key (k).
*/
func (ds *DataStore) get(k string) Room {
	return ds.datastore[k]
}

/*
Removes the entry for the given key(k)
*/
func (ds *DataStore) removeKey(k string) {
	delete(ds.datastore, k)
}

/*
Removes the entry for the given key(k)
*/
func (ds *DataStore) removeKeys(k ...string) {
	for _, d := range k {
		delete(ds.datastore, d)
	}
}

/*
Prints the keys and values
*/
func (ds *DataStore) print() {
	ds.Lock()
	for k, v := range ds.datastore {
		fmt.Println(k, " ", v)
	}
	ds.Unlock()
}
