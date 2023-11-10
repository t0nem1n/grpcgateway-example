FROM golang:1.21 as build

RUN apt-get update && apt install -y unzip

ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
RUN curl -LO $PB_REL/download/v24.4/protoc-24.4-linux-x86_64.zip && unzip protoc-24.4-linux-x86_64.zip -d /local

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="$PATH:/local/bin:$(go env GOPATH)/bin"

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

