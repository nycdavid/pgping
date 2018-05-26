FROM golang:1.10.2-alpine3.7

RUN mkdir -p /go/src/github.com/velvetreactor/pgping

WORKDIR /go/src/github.com/velvetreactor/pgping

COPY ./ ./

RUN apk update && \
  apk add curl git && \
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
  dep ensure
