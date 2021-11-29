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
# Download all the dependencies that are required in your source files and update go.mod file with that dependency.
# Remove all dependencies from the go.mod file which are not required in the source files.
RUN go mod tidy

# enviroment
WORKDIR cmd/pro
RUN go build -o app

ENTRYPOINT ["./app"]

EXPOSE 8080