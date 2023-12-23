FROM golang:1.21-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY pkg ./pkg

WORKDIR /app/cmd/driver
RUN go build -o main
