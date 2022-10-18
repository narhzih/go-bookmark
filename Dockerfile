FROM golang:alpine AS builder
MAINTAINER Mypipe Developers <dev@mypipeapp.com>

# Copy base files and install dependency
COPY go.mod go.sum /mypipeapp/
WORKDIR /mypipeapp/
RUN go mod download

# Copy project files and build
COPY . /mypipeapp
COPY .env.staging /mypipeapp/.env.docker

# Remove test env files
RUN rm /mypipeapp/.env && rm /mypipeapp/.env.staging && rm /mypipeapp/.env.example

# Create build
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/api ./cmd/api

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /mypipeapp/bin/api /usr/bin/api

# start project files
ENTRYPOINT ["/usr/bin/api"]