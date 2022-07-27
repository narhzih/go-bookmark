FROM golang:alpine AS builder 

COPY go.mod go.sum /go/src/gitlab.com/trencetech/mypipe-api/
WORKDIR /go/src/gitlab.com/trencetech/mypipe-api/
RUN go mod download
COPY . /go/src/gitlab.com/trencetech/mypipe-api/
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o bin/mypipe ./main.go

FROM alpine 
RUN apk add --no-cache ca-certificates && update-ca-certificates
# Install make
RUN apt update && apt install -y make
COPY --from=builder /go/src/gitlab.com/trencetech/mypipe-api/bin/mypipe /usr/bin/mypipe


# COPY .env /usr/bin/mypipe
RUN  chmod +x /usr/bin/mypipe
EXPOSE 5555
ENTRYPOINT ["/usr/bin/mypipe"]