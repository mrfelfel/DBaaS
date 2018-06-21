package main

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

func GetContainers(w http.ResponseWriter, r *http.Request) {
	
	out, _ := json.Marshal(containers)
	log.Print(string(out))
	json.NewEncoder(w).Encode(containers)
}

func GetContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range containers{
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}		
	}
	json.NewEncoder(w).Encode(&Container{})
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var container Container
	_ = json.NewDecoder(r.Body).Decode(&container)
	container.ID = params["id"]
	containers = append(containers, container)
	json.NewEncoder(w).Encode(containers)
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range containers {
		if item.ID == params["id"]{
			containers = append(containers[:index], containers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(containers)
}

// our main function
func main() {

	containers = append(containers, Container{"1.1", "base"})
	containers = append(containers, Container{ID:"2.1", Name:"testing"})
	containers = append(containers, Container{ID:"3.1", Name:"prod"})

	router := mux.NewRouter()
	router.HandleFunc("/container", GetContainers).Methods("GET")
	router.HandleFunc("/container/{id}", GetContainer).Methods("GET")
	router.HandleFunc("/container/{id}", CreateContainer).Methods("POST")
	router.HandleFunc("/container/{id}", DeleteContainer).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}


type User struct{
	ID 	string	//`json:"id,omitempty"`
	name	string	//`json:"name,omitempty"`
	containers	*Container //`json:"containers,omitempty"`
}

type Container struct{
	ID 	string	`json:"id_nums,omitempty"`
	Name	string	`json:"name,omitempty"`
}

var user []User
var containers []Container


