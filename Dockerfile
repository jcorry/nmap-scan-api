FROM golang:1.12.0-alpine

RUN apk --update upgrade
RUN apk add bash
RUN apk add sqlite
RUN apk add gcc
RUN apk add g++
RUN apk add git mercurial

RUN go get github.com/mattn/go-sqlite3
RUN go install github.com/mattn/go-sqlite3

# removing apk cache
RUN rm -rf /var/cache/apk/*

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN touch db/data.db

ENV CGO_ENABLED=1

RUN go build -v -mod=vendor ./cmd/api

EXPOSE 8080

CMD ["./api"]