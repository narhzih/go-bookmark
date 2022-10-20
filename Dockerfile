FROM golang:alpine AS builder
MAINTAINER Mypipe Developers <dev@mypipeapp.com>

# Copy base files and install dependency
COPY go.mod go.sum /mypipeapp/
WORKDIR /mypipeapp/
RUN go mod download

# Copy project files and build
COPY . /mypipeapp

# Create build
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/api ./cmd/api

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /mypipeapp/ /usr/mypipeapp

# start project files
ENV APP_ENV=prod
ENTRYPOINT ["/usr/mypipeapp/bin/api"]