#A Dockerfile to test usergo
FROM ubuntu:14.04

MAINTAINER davejfranco <davefranco1987@gmail.com>

RUN apt-get update
RUN apt-get install wget git -y

#Golang Installation
RUN wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
RUN tar -xvf go1.8.3.linux-amd64.tar.gz
ENV GOROOT /go
ENV PATH "$PATH:/go/bin"
ENV GOPATH /code

#ADD code
RUN mkdir -p /code/src/github.com/davejfranco/usergo
ADD . /code/src/github.com/davejfranco/usergo

WORKDIR /code/src/github.com/davejfranco/usergo
