FROM golang:1.16-alpine

# install git
RUN apk update && apk add git

ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/go-module

COPY . .

# Build
RUN go mod init go-module
RUN go mod tidy
RUN go mod vendor
RUN go build -o app

ENTRYPOINT ["./app"]

EXPOSE 8080