all:
	go get github.com/go-sql-driver/mysql
	go get github.com/gorilla/mux
	go get github.com/pkg/errors
	go get golang.org/x/net/proxy
	go get github.com/docker/docker/api/types
	go get github.com/docker/docker/api/types/container
	go get github.com/docker/docker/client
	go get github.com/docker/go-connections/nat


	echo -e "Done getting dependencies... test"

	if [ -d "${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/go-connections" ]; then \
		rm -rf ${GOPATH}/src/github.com/docker/go-connections; \
		mv ${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/go-connections ${GOPATH}/src/github.com/docker; \
	fi

	go build -o bin/api api/api.go

linux:	
	go get github.com/go-sql-driver/mysql
	go get github.com/gorilla/mux
	go get github.com/pkg/errors
	go get golang.org/x/net/proxy
	go get github.com/docker/docker/api/types
	go get github.com/docker/docker/api/types/container
	go get github.com/docker/docker/client
	go get github.com/docker/go-connections/nat


	echo -e "Done getting dependencies..."

	if [ -d "${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/go-connections" ]; then \
		rm -rf ${GOPATH}/src/github.com/docker/go-connections; \
		mv ${GOPATH}/src/github.com/docker/docker/vendor/github.com/docker/go-connections ${GOPATH}/src/github.com/docker; \
	fi

	go build -o bin/api api/api.go

windows:
	cd "${GOPATH}\\src\\github.com\\docker\\docker\\vendor\\github.com\\docker\\"
	mv go-connections ${GOPATH}\\src\\github.com\\docker

	cd ${HOME}\\gopath\\src\\github.com\\sixtop\\DBaaS

	echo "what"

	go build -o bin/api.exe api/api.go

