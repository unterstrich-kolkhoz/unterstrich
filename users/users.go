package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

  "github.com/hellerve/artifex/model"
)

type User struct {
  model.Base
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
	Social    *Social  `json:"social,omitempty"`
}
type Address struct {
  model.Base
	Line1 string `json:"line1,omitempty"`
	Line2 string `json:"line2,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}
type Social struct {
  model.Base
	Github  string `json:"github,omitempty"`
	Ello    string `json:"ello,omitempty"`
	Website string `json:"website,omitempty"`
}


func Initialize(db* gorm.DB, router* mux.Router) {
	router.HandleFunc("/users", endpoint(db, GetUsers)).Methods("GET")
	router.HandleFunc("/users", endpoint(db, CreateUser)).Methods("POST")
	router.HandleFunc("/users/{id}", endpoint(db, GetUser)).Methods("GET")
	router.HandleFunc("/users/{id}", endpoint(db, UpdateUser)).Methods("PUT")
	router.HandleFunc("/users/{id}", endpoint(db, DeleteUser)).Methods("DELETE")

  db.AutoMigrate(&User{}, &Address{}, &Social{})
}

func endpoint(db* gorm.DB, wrapped func(http.ResponseWriter,
*http.Request, *gorm.DB)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,
		r *http.Request) {
    wrapped(w, r, db)
  }
}

func GetUsers(w http.ResponseWriter, r *http.Request, db* gorm.DB) {
  var users []User
  db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request, db* gorm.DB) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
    return
	}

  var user *User
  db.First(user, id)

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid ID: not found"))
    return
	}

	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request, db* gorm.DB) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body: "))
		w.Write([]byte(err.Error()))
    return
	}

  if !db.NewRecord(user) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already present: "))
		w.Write([]byte(string(user.ID)))
    return
  }

  db.Create(&user)

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request, db* gorm.DB) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
    return
	}

  var user *User
  db.First(user, id)

  if user == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
    return
  }

	json.NewEncoder(w).Encode("")
}

func UpdateUser(w http.ResponseWriter, r *http.Request, db* gorm.DB) {
	params := mux.Vars(r)

	_, err := strconv.Atoi(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID: must be numerical"))
	}

	var user *User
	err = json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid body: "))
		w.Write([]byte(err.Error()))
	}

  if db.NewRecord(user) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
    return
  }

  db.Save(&user)

	json.NewEncoder(w).Encode(user)
}
