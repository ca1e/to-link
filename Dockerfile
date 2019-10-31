FROM golang:alpine as builder

RUN apk --no-cache add git ca-certificates

WORKDIR /go/src/github.com/go/tolnkcn/

COPY to-link/src/ .

RUN go get -d -v github.com/gin-gonic/gin \
  && go get -d -v github.com/gomodule/redigo/redis \
  &&  go get -d -v github.com/spaolacci/murmur3

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app . 

FROM alpine:latest as prod

WORKDIR /root/
COPY to-link/ .

COPY --from=0 /go/src/github.com/go/tolnkcn/app ./bin/
# ENV GIN_MODE=release
WORKDIR /root/bin
CMD ["./app"]
