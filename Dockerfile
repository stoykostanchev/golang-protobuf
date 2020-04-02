FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN mkdir ./backend
RUN mkdir -p ./backend/generated

RUN apk update && apk add bash ca-certificates gcc g++ libc-dev

WORKDIR /backend

COPY ./generated/proto/joke.pb.go /backend/generated/proto/
COPY ./main.go /backend

COPY ./go.mod /backend
COPY ./go.sum /backend

RUN go mod download

RUN go build -o backend .

EXPOSE 8080
CMD ./backend