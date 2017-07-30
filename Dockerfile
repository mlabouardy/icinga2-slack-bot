FROM golang:1.8.0
MAINTAINER mlabouardy <mohamed@labouardy.com>

RUN mkdir -p /go/src/github

COPY . /go/src/github/

WORKDIR /go/src/github/

RUN go get -v

CMD go run $(ls -1 *.go | grep -v _test.go)
