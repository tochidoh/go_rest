package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/tochidoh/go_rest/data"
)

var (
	getUsersRe   = regexp.MustCompile(`\/user[\/]*$`)
	getUserRe    = regexp.MustCompile(`\/user\/(\d+)$`)
	createUserRe = regexp.MustCompile(`\/user[\/]*$`)
)

type UserHandler struct {
	Store *data.Datastore
	Num   int
}

func (uh *UserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("serving this user handler")
	log.Println(r.Method, r.URL.Path)
	log.Println(getUsersRe.MatchString(r.URL.Path))
	log.Println(getUserRe.MatchString(r.URL.Path))
	log.Println(createUserRe.MatchString(r.URL.Path))

	rw.Header().Set("content-type", "application/json")

	switch {
	case r.Method == http.MethodGet && getUsersRe.MatchString(r.URL.Path):
		log.Println("get users")
		uh.GetUsers(rw, r)
		return
	case r.Method == http.MethodGet && getUserRe.MatchString(r.URL.Path):
		log.Println("get single user")
		uh.GetUser(rw, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		log.Println("post single user")
		uh.CreateUser(rw, r)
		return
	default:
		log.Println("goes down default path")
		fmt.Println("something hereeeeee")
		// could not find a valid handler path
		uh.NotFound(rw, r)
		return
	}
}

func (uh *UserHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	users := make([]data.User, 0, len(uh.Store.M))
	uh.Store.RLock()
	for _, user := range uh.Store.M {
		users = append(users, user)
	}
	uh.Store.RUnlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		internalServerError(rw, r)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}

func (uh *UserHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	// need to extract id from the path
	matches := getUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		uh.NotFound(rw, r)
		return
	}
	uh.Store.RLock()
	user, ok := uh.Store.M[matches[1]]
	uh.Store.RUnlock()
	if !ok {
		uh.NotFound(rw, r)
		return
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		internalServerError(rw, r)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)

}

func (uh *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	user := data.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		badRequest(rw, r)
		return
	}

	uh.Store.Lock()
	id := strconv.Itoa(user.Id)
	uh.Store.M[id] = user
	uh.Store.Unlock()

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		internalServerError(rw, r)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}

func (uh *UserHandler) NotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("not found"))
	rw.Write([]byte(`{"error": "not found"}`))
}

func internalServerError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(`{"error": "internal server error"`))
}

func badRequest(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(`{"error": "bad request"}`))
}
