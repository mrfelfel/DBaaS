package main

import (
	"encoding/json"
	"net/http"
	"log"
	"context"
	"strings"
	"time"
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

/* Gets a list of all containers */
func __getContainers() ([]types.Container, error) {

	// Get container data
	dckr_res, dckr_err := dock_cli.ContainerList(context.Background(), types.ContainerListOptions{All:true})

	return dckr_res, dckr_err
}

/* Gets a specific container reference */
func __getContainer(id string) (*types.Container, error) {

	// Get container data
	dckr_res, dckr_err := __getContainers()
	if dckr_err != nil {
		return nil, dckr_err
	}

	// loop through containers
	for _, item := range dckr_res{

		// Match parameter with container ID
		if strings.HasPrefix(item.ID, id){
			return &item, nil
		}		
	}
	return nil, nil
}

/* This function gets a list of IDs of the currently running containers */
func GetContainers(w http.ResponseWriter, r *http.Request) {
	
	// Get container data
	dckr_res, dckr_err := dock_cli.ContainerList(context.Background(), types.ContainerListOptions{All:true})
	if dckr_err != nil {
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract container IDs
	id_list := make([]string, len(dckr_res))
	for i := range dckr_res{
		id_list[i] = dckr_res[i].ID
	}

	// pretty print
	out, json_err := json.MarshalIndent(id_list,"", "\t")
	if json_err != nil{
		http.Error(w, json_err.Error(), http.StatusInternalServerError)
    		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}


/* Gets information of a selected container*/
func GetContainer(w http.ResponseWriter, r *http.Request) {

	// get URL parameters
	params := mux.Vars(r)

	// Get container data
	dckr_res, dckr_err := __getContainer(params["cid"])
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// Not found
	if dckr_res == nil{
		json.NewEncoder(w).Encode(params["cid"] + " not found.")
		return
	}

	// pretty print
	json_out, json_err := json.MarshalIndent(dckr_res,"", "\t")
	if json_err != nil{
		http.Error(w, json_err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_out)
}

/* Creates a container */
func CreateContainer(w http.ResponseWriter, r *http.Request) {

	// extract parameters & build the config request
	log.Print(r.FormValue("Image"))
	log.Print(r.FormValue("Cmd"))
	var config container.Config
	config.Image = r.FormValue("Image")
	config.Cmd = []string {r.FormValue("Cmd")}

	// create container
	dckr_res, dckr_err := dock_cli.ContainerCreate(context.Background(), &config ,nil, nil, r.FormValue("id"))
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("start") == "true" {
		log.Print("Starting container ")
		dckr_err := dock_cli.ContainerStart(context.Background(), dckr_res.ID,  types.ContainerStartOptions{})
		if dckr_err != nil{
			http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// pretty print
	json_out, json_err := json.MarshalIndent(dckr_res,"", "\t")
	if json_err != nil{
		http.Error(w, json_err.Error(), http.StatusInternalServerError)
    		return
	}

	//GetContainer(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json_out)
}

/* Deletes a container */
func DeleteContainer(w http.ResponseWriter, r *http.Request) {

	// get URL parameters
	params := mux.Vars(r)

	// Get container data
	dckr_res, dckr_err := __getContainer(params["cid"])
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// Not found
	if dckr_res == nil{
		json.NewEncoder(w).Encode(params["cid"] + " not found.")
		return
	}

	//remove it
	dckr_err = dock_cli.ContainerRemove(context.Background(), dckr_res.ID, types.ContainerRemoveOptions{false,false,true});
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	json.NewEncoder(w).Encode("removed: " + dckr_res.ID)
}

func StartContainer(w http.ResponseWriter, r *http.Request) {

	// get URL parameters
	params := mux.Vars(r)

	// Get container data
	dckr_res, dckr_err := __getContainer(params["cid"])
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// Not found
	if dckr_res == nil{
		json.NewEncoder(w).Encode(params["cid"] + " not found.")
		return
	}

	dckr_err = dock_cli.ContainerStart(context.Background(), dckr_res.ID,  types.ContainerStartOptions{})
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}
	// write response
	json.NewEncoder(w).Encode("Container started: " + dckr_res.ID)
}

func StopContainer(w http.ResponseWriter, r *http.Request) {

	// get URL parameters
	params := mux.Vars(r)

	// Get container data
	dckr_res, dckr_err := __getContainer(params["cid"])
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}

	// Not found
	if dckr_res == nil{
		json.NewEncoder(w).Encode(params["cid"] + " not found.")
		return
	}

	timeout := time.Second
	dckr_err = dock_cli.ContainerStop(context.Background(), dckr_res.ID, &timeout)
	if dckr_err != nil{
		http.Error(w, dckr_err.Error(), http.StatusInternalServerError)
		return
	}
	// write response
	json.NewEncoder(w).Encode("Container stopped: " + dckr_res.ID)
}

func ListDatabases(w http.ResponseWriter, r *http.Request) {	
}

func CreateDatabase(w http.ResponseWriter, r *http.Request) {	
}

func GetDatabase(w http.ResponseWriter, r *http.Request) {	
}

func RemoveDatabase(w http.ResponseWriter, r *http.Request) {	
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

	/* Container operations */
	router.HandleFunc("/container", GetContainers).Methods("GET")
	router.HandleFunc("/container/create", CreateContainer).Methods("POST")
	router.HandleFunc("/container/{cid}", DeleteContainer).Methods("DELETE")
	router.HandleFunc("/container/{cid}", GetContainer).Methods("GET")

	/* service operations */
	router.HandleFunc("/container/{cid}/start", StartContainer).Methods("POST")
	router.HandleFunc("/container/{cid}/stop", StopContainer).Methods("POST")

	/* database operations */
	router.HandleFunc("/container/{cid}/list", ListDatabases).Methods("GET")
	router.HandleFunc("/container/{cid}/createDB", CreateDatabase).Methods("POST")
	router.HandleFunc("/container/{cid}/{dbid}", GetDatabase).Methods("GET")
	router.HandleFunc("/container/{cid}/{dbid}", RemoveDatabase).Methods("DELETE")


	//start listening
	log.Fatal(http.ListenAndServe(":8000", router))
}

var containers []Container