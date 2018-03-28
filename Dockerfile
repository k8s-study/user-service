FROM golang:1.10-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR $GOPATH/src/github.com/k8s-study/user-service
ADD . $GOPATH/src/github.com/k8s-study/user-service

RUN dep ensure

EXPOSE 8080

CMD go run main.go
