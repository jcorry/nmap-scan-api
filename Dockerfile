FROM golang:1.12.0-alpine
# Get OS/build dependencies for cgo
RUN apk --update upgrade
RUN apk add bash
RUN apk add sqlite
RUN apk add gcc
RUN apk add g++
RUN apk add git mercurial

# Get and install cgo driver
RUN go get github.com/mattn/go-sqlite3
RUN go install -i github.com/mattn/go-sqlite3

# Set workdir
WORKDIR /app

# Remove apk cache
RUN rm -rf /var/cache/apk/*
# env vars
ENV GO111MODULE=on
ENV CGO_ENABLED=1

# Move source files
COPY . .

# Since this wasn't in the repo, create it here...will not persist between image builds
RUN touch db/data.db

# Build the app
RUN go build -v -mod=vendor ./cmd/api

# New container in which this will run
FROM alpine
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy necessary files
COPY --from=0 /app/api .
COPY --from=0 /app/db ./db
COPY --from=0 /app/sql ./sql
CMD ["./api"]

# Port
EXPOSE 8080

# Run command to start service
CMD ["./api"]