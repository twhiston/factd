FROM golang:1.10

LABEL maintainer="Tom Whiston <tom.whiston@gmail.com>"

RUN apt-get update -y && apt-get install -y rpm --no-install-recommends && \
    rm -rf /var/lib/apt/lists/* && \
    go get -u github.com/golang/dep/cmd/dep && \
    go get github.com/schrej/godacov && \
    go get github.com/goreleaser/goreleaser && \
    go get -u github.com/haya14busa/goverage && \
    go get -u github.com/alecthomas/gometalinter && \
    gometalinter --install