all:
	go get -t -v ./...
	echo -e "Done getting dependencies..."
	cd "${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker"
	mv go-connections ${GOPATH}/src/github.com/docker
	cd ${HOME}/gopath/src/github.com/sixtop/DBaaS

	go build -o bin/api.exe api/api.go

windows:
	cd "${GOPATH}\\src\\github.com\\docker\\docker\\vendor\\github.com\\docker\\"
	mv go-connections ${GOPATH}\\src\\github.com\\docker

	cd ${HOME}\\gopath\\src\\github.com\\sixtop\\DBaaS

	echo "what"

	go build -o bin/api.exe api/api.go

