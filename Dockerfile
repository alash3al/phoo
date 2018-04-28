FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/alash3al/http2fcgi

ENTRYPOINT ["http2fcgi"]

WORKDIR /root/