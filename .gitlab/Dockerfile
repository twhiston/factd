FROM golang:1.10

RUN apt-get update -y && apt-get install -y rpm --no-install-recommends && \
    go get -u github.com/golang/dep/cmd/dep && \
    go get github.com/schrej/godacov && \
    go get github.com/goreleaser/goreleaser && \
    go get -u github.com/haya14busa/goverage && \
    go get -u github.com/alecthomas/gometalinter && \
    gometalinter --install