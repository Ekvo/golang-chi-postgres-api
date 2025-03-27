FROM golang:alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy