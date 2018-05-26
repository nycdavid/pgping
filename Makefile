test:
	docker run \
	-v $(shell pwd):/go/src/github.com/velvetreactor/pgping/ \
	nycdavid/pgping:latest \
	go test -v ./...

image:
	docker build -t nycdavid/pgping:latest .
