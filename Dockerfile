FROM golang:1.21.0 AS builder

LABEL stage=builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/src/build

ADD go.mod ./
ADD go.sum ./

RUN go mod download

COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./cmd ./cmd

RUN go build -o task ./cmd/app/main.go

FROM alpine

LABEL authors="ekvo"

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=task-manager
ENV DB_PASSWORD=qwert12345
ENV DB_NAME=task-store
ENV DB_SSLMODE=disable
ENV SRV_ADDR=3000

EXPOSE ${SRV_ADDR}

WORKDIR /usr/src/app

RUN apk update && \
    apk add postgresql-client

COPY --from=builder /usr/src/build/task /usr/src/app/task

COPY script/start.sh /start.sh
RUN chmod +x /start.sh

#ENTRYPOINT ["sh","/start.sh"]
