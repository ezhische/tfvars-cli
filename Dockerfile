FROM golang:1.23.4

COPY . /
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /
RUN go mod download
RUN go build -ldflags="-s -w" -o ./cmd/main.go
ENTRYPOINT [ "sleep","600" ]