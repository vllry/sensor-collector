FROM golang:1.13 as builder

RUN apt-get update && apt-get install -y zip

RUN adduser --disabled-password --gecos '' app
RUN mkdir -p /go/src/github.com/vllry/sensor-collector

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip && \
    unzip protoc-3.10.0-linux-x86_64.zip && mv ./bin/protoc /usr/bin/ && mv ./include/ /usr/local/include/

COPY ./pkg /go/src/github.com/vllry/sensor-collector/pkg
COPY main.go go.mod go.sum /go/src/github.com/vllry/sensor-collector/

WORKDIR /go/src/github.com/vllry/sensor-collector
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
RUN chmod +x sensor-collector



FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/vllry/sensor-collector/sensor-collector /app

USER app
CMD ["/app"]