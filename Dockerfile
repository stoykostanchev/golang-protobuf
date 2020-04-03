FROM golang:alpine as build-env

RUN apk update && apk add bash ca-certificates gcc g++ libc-dev

WORKDIR /backend

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o backend .

EXPOSE 8080
CMD ./backend