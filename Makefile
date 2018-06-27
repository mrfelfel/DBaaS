all:
	go get -t -v ./...
	echo -e "Done getting dependencies..."

	cd ${GOPATH}/src/github.com/docker
	rm -rf go-connections


	cd "${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/"
	mv go-connections ${GOPATH}/src/github.com/docker

	cd ${HOME}/gopath/src/github.com/sixtop/DBaaS
	go build -o bin/api api/api.go

linux:
	
	echo -e "Done getting dependencies..."

	rm -rf ${GOPATH}/src/github.com/docker/go-connections
	mv ${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/go-connections ${GOPATH}/src/github.com/docker

	cd ${HOME}/gopath/src/github.com/sixtop/DBaaS
	go build -o bin/api api/api.go

windows:
	cd "${GOPATH}\\src\\github.com\\docker\\docker\\vendor\\github.com\\docker\\"
	mv go-connections ${GOPATH}\\src\\github.com\\docker

	cd ${HOME}\\gopath\\src\\github.com\\sixtop\\DBaaS

	echo "what"

	go build -o bin/api.exe api/api.go

