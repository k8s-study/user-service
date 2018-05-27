FROM golang:1.10-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR $GOPATH/src/github.com/k8s-study/user-service
ADD . $GOPATH/src/github.com/k8s-study/user-service

RUN dep ensure

ENV PORT 8080
ENV DB_HOST localhost
ENV DB_PORT 5432
ENV DB_NAME users
ENV DB_USERNAME postgres
ENV DB_PASSWORD postgres
ENV KONG_HOST http://kong-ingress-controller.kong:8001

EXPOSE 8080

CMD go run main.go
