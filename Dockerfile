FROM golang:1.16-alpine

# install git
RUN apk update && apk add git

ENV GO111MODULE=on

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/course-golang

COPY . .

# Build
RUN go mod init course-golang
RUN go build -o main

ENTRYPOINT ["./app/main"]

EXPOSE 8080