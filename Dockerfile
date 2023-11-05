FROM golang:1.20.4-alpine

RUN set -xe && \
	apk add git make

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV GO111MODULE=on

ADD ./src/ /g/
WORKDIR /g

RUN set -xe && \
    rm -rf /g/release && \
    make release