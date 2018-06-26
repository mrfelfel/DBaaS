package main

import (
	"encoding/json"
	"net/http"
	"log"
	"context"
	"strings"
	"time"
	"database/sql"
	"fmt"
	"errors"

	"github.com/gorilla/mux"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	_ "github.com/go-sql-driver/mysql" //mysql driver
)

var dock_cli *client.Client 	// Golang Docker Client 

// Exposed Docker Container information
type Container struct {
	ID string
	Name string
	Image string
	Cmd string
}

func __getDatabase(w http.ResponseWriter, r *http.Request) (*sql.DB, error){
	// get URL parameters
	params := mux.Vars(r)

	// Get container data
	dckr_res, dckr_err := __getContainer(params["cid"])
	if dckr_err != nil{
		return nil, dckr_err
	}

	// Not found
	if dckr_res == nil{
		return nil, errors.New("docker error: " + params["cid"] + " not found.")
	}

	// set up connection
	user := "docker"
	pass := "docker123"
	ip := "10.0.75.1"
	port := dckr_res.Ports[0].PublicPort
	host := fmt.Sprintf("%s:%s@tcp(%s:%d)/",user, pass, ip, port)
	//log.Print(host)
	db, db_err := sql.Open("mysql", host)
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return nil, db_err
	}

	return db, nil
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
	return nil, errors.New("Docker error: " + id + " not found.")
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

	// get URL parameters
	params := mux.Vars(r)

	// extract parameters & build the config request
	log.Print(r.FormValue("Image"))
	log.Print(r.FormValue("Cmd"))

	config := &container.Config{
		Image : r.FormValue("Image"),
		Cmd : []string {r.FormValue("Cmd")},
		ExposedPorts: nat.PortSet{
			nat.Port("3306/tcp"):{},
		},
	}

    	hostConfig := &container.HostConfig{
    		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
		},
    		PortBindings: nat.PortMap{
    			"3306/tcp": []nat.PortBinding{{HostIP:"0.0.0.0"}},
    		},
	}

	// create container
	dckr_res, dckr_err := dock_cli.ContainerCreate(context.Background(), config ,hostConfig, nil, params["cid"])
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

	db, db_err := __getDatabase(w,r)
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// ping db
	db_err = db.Ping()
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}

	// query away!
	rows, db_err := db.Query("show databases")
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dbname string
	var output []string
	for rows.Next() {
		err := rows.Scan(&dbname)
		if err != nil {
			log.Fatal(err)
		}
		output = append(output, dbname)
	}
	json.NewEncoder(w).Encode(output)
}

func CreateDatabase(w http.ResponseWriter, r *http.Request) {
	db, db_err := __getDatabase(w,r)
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// ping db
	db_err = db.Ping()
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}

	// get URL parameters
	params := mux.Vars(r)
	rows, db_err := db.Query("create database " + params["dbid"])
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	/*var res string
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			log.Fatal(err)
		}
	}*/

	json.NewEncoder(w).Encode("Created database " + params["dbid"] + ".")
}

func GetDatabase(w http.ResponseWriter, r *http.Request) {
}

func RemoveDatabase(w http.ResponseWriter, r *http.Request) {
	db, db_err := __getDatabase(w,r)
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// ping db
	db_err = db.Ping()
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}

	// get URL parameters
	params := mux.Vars(r)

	rows, db_err := db.Query("drop database " + params["dbid"])
	if db_err != nil {
		http.Error(w, db_err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	/*var res string
	for rows.Next() {
		err := rows.Scan(&res)
		if err != nil {
			log.Fatal(err)
		}
	}*/

	json.NewEncoder(w).Encode("Dropped database " + params["dbid"] + ".")
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
	router.HandleFunc("/container/{cid}", GetContainer).Methods("GET")
	router.HandleFunc("/container/{cid}", CreateContainer).Methods("POST")
	router.HandleFunc("/container/{cid}", DeleteContainer).Methods("DELETE")

	/* service operations */
	router.HandleFunc("/container/{cid}/start", StartContainer).Methods("POST")
	router.HandleFunc("/container/{cid}/stop", StopContainer).Methods("POST")

	/* database operations */
	router.HandleFunc("/database/{cid}", ListDatabases).Methods("GET")
	router.HandleFunc("/database/{cid}/{dbid}", GetDatabase).Methods("GET")
	router.HandleFunc("/database/{cid}/{dbid}", CreateDatabase).Methods("POST")
	router.HandleFunc("/database/{cid}/{dbid}", RemoveDatabase).Methods("DELETE")


	//start listening
	log.Fatal(http.ListenAndServe(":8000", router))
}

var containers []Container