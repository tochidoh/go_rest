package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/tochidoh/go_rest/data"
	"github.com/tochidoh/go_rest/handlers"
)

// end points

// get /users
// get /users/{id}
// post /users

func main() {
	fmt.Println("hello world")

	mux := http.NewServeMux()

	// should make a constructor for this
	userHandler := &handlers.UserHandler{
		Store: &data.Datastore{
			M: map[string]data.User{
				"1": data.User{Id: 1, Name: "bob"},
			},
			RWMutex: &sync.RWMutex{},
		},
	}

	mux.Handle("/user/", userHandler)
	mux.Handle("/user", userHandler)

	http.ListenAndServe("localhost:8080", mux)
}
