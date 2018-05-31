test:
	docker run \
	-v $(shell pwd)/main.go:/go/src/github.com/velvetreactor/pgping/main.go \
	-v $(shell pwd)/main_test.go:/go/src/github.com/velvetreactor/pgping/main_test.go \
	nycdavid/pgping:latest \
	go test -v ./...

image:
	docker build -t nycdavid/pgping:latest .
