#A Dockerfile to test usergo
FROM ubuntu:14.04

MAINTAINER davejfranco <davefranco1987@gmail.com>

RUN apt-get update
RUN apt-get install wget git vim -y

#Golang Installation
RUN wget https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
RUN tar -xvf go1.8.3.linux-amd64.tar.gz
ENV GOROOT /go
ENV PATH "$PATH:/go/bin"

RUN mkdir /code
WORKDIR /code

COPY . /code
