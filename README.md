# DBaaS
DBaaS challenge 

[![Build Status](https://travis-ci.org/sixtop/DBaaS.png?branch=master)](https://travis-ci.org/sixtop/DBaaS)

# Endpoints

### Container operations
	router.HandleFunc("/container", GetContainers).Methods("GET")
	router.HandleFunc("/container/{cid}", GetContainer).Methods("GET")
	router.HandleFunc("/container/{cid}", CreateContainer).Methods("POST")
	router.HandleFunc("/container/{cid}", DeleteContainer).Methods("DELETE")

### service operations
	router.HandleFunc("/container/{cid}/start", StartContainer).Methods("POST")
	router.HandleFunc("/container/{cid}/stop", StopContainer).Methods("POST")

### database operations
	router.HandleFunc("/database/{cid}", ListDatabases).Methods("GET")
	router.HandleFunc("/database/{cid}/{dbid}", GetDatabase).Methods("GET")
	router.HandleFunc("/database/{cid}/{dbid}", CreateDatabase).Methods("POST")
	router.HandleFunc("/database/{cid}/{dbid}", RemoveDatabase).Methods("DELETE")
