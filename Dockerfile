# start a golang base image, version 1.8
FROM golang:1.14

WORKDIR /go/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /server ./cmd/server/*.go
COPY ./openapi.yaml /openapi.yaml
CMD ["/server"]