FROM golang:1.8.3

RUN mkdir -p /Golang_api
WORKDIR /Golang_api

ADD . /Golang_api

RUN go get -v

