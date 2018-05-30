test:
	docker run \
	-v $(shell pwd):/go/src/github.com/velvetreactor/pgping/ \
	nycdavid/pgping:latest \
	go test -v ./...

cmd:
	docker run \
	-v $(shell pwd):/go/src/github.com/velvetreactor/pgping/ \
	nycdavid/pgping:latest \
	$(CMD)

image:
	docker build -t nycdavid/pgping:latest .
