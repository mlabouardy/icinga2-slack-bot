FROM golang:1.8.0
MAINTAINER mlabouardy <mohamed@labouardy.com>

COPY . /app

WORKDIR /app

RUN go get github.com/nlopes/slack
RUN go get github.com/BurntSushi/toml

CMD go run *.go
