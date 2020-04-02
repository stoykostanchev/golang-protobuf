FROM golang:alpine as build-env

RUN mkdir ./backend

RUN apk update && apk add bash ca-certificates gcc g++ libc-dev

WORKDIR /backend

COPY ./main.go /backend

RUN go build -o backend .

EXPOSE 8080
CMD ./backend