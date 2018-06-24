package main

import (
	"encoding/json"
	"net/http"
	"log"
	"context"
	"strings"
	//"fmt"

	"github.com/gorilla/mux"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var dock_cli *client.Client 	// Golang Docker Client 

// Exposed Docker Container information
type Container struct {
	ID string
	Name string
    Image string
    Cmd string
}

/* This function gets a list of IDs of the currently running containers */
func GetContainers(w http.ResponseWriter, r *http.Request) {
	
	// Get container data
	crs, err := dock_cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil{
		panic(err)
		return
	}

	// Extract container IDs
	id_list := make([]string, len(crs))
	for i := range crs{
		id_list[i] = crs[i].ID
	}

	// pretty print
	out, jsonerr := json.MarshalIndent(id_list,"", "\t")
	if jsonerr != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}


/* Gets information of a selected container*/
func GetContainer(w http.ResponseWriter, r *http.Request) {

	// Get container data
	crs, err := dock_cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil{
		panic(err)
		return
	}

	// get URL parameters
	params := mux.Vars(r)

	// loop through containers
	for _, item := range crs{

		// Match parameter with container ID
		if strings.HasPrefix(item.ID, params["id"]){

			// pretty print
			out, jsonerr := json.MarshalIndent(item,"", "\t")
			if jsonerr != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
		    	return
			}

			// write response
			w.Header().Set("Content-Type", "application/json")
			w.Write(out)
			return
		}		
	}

	//write empty list
	json.NewEncoder(w).Encode(&Container{})
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	// extract parametes
	params := mux.Vars(r)
	var c Container
	_ = json.NewDecoder(r.Body).Decode(&c)
	c.ID = params["id"]
	c.Image = r.FormValue("image")
	c.Cmd = r.FormValue("cmd")
	c.Image = r.FormValue("test")

	log.Print(r.FormValue("image"))
	log.Print(r.FormValue("cmd"))
	log.Print(r.FormValue("test"))

	// create container
	var config container.Config
	config.Image = c.Image
	config.Cmd = []string {c.Cmd}

	resp, err := dock_cli.ContainerCreate(context.Background(), &config ,nil, nil, c.ID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// pretty print
	out, jsonerr := json.MarshalIndent(config,"", "\t")
	if jsonerr != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}

	//GetContainer(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
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

	// Connect to Docker server
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
	if err!=nil {
		panic(err)
	}
	dock_cli = cli

	// Create RESTful API endpoints
	router := mux.NewRouter()
	router.HandleFunc("/container", GetContainers).Methods("GET")
	router.HandleFunc("/container/create", CreateContainer).Methods("POST")
	router.HandleFunc("/container/{id}", GetContainer).Methods("GET")
	router.HandleFunc("/container/{id}/start", StartContainer).Methods("POST")
	router.HandleFunc("/container/{id}/stop", StopContainer).Methods("POST")
	router.HandleFunc("/container/{id}", DeleteContainer).Methods("DELETE")

	//start listening
	log.Fatal(http.ListenAndServe(":8000", router))
}

var containers []Container