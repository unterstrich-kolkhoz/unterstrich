package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID        int      `json:"id,omitempty"`
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
	Social    *Social  `json:"social,omitempty"`
}
type Address struct {
	Line1 string `json:"line1,omitempty"`
	Line2 string `json:"line2,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}
type Social struct {
	Github  string `json:"github,omitempty"`
	Ello    string `json:"ello,omitempty"`
	Website string `json:"website,omitempty"`
}

var users []User

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
	}

	found := false
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			found = true
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid ID: not found"))
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body: "))
		w.Write([]byte(err.Error()))
	}

	user.ID = len(users) + 1
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
	}

	users = append(users[:id-1], users[id:]...)

	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body: "))
		w.Write([]byte(err.Error()))
	}

	user.ID = id

	found := false
	for _, u := range users {
		if u.ID == id {
			found = true
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid ID: not found"))
	}

	users[id-1] = user
	json.NewEncoder(w).Encode(user)
}
