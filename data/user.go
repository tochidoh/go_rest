package data

import "sync"

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Datastore struct {
	M map[string]User
	*sync.RWMutex
}
