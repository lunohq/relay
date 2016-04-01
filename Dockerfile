FROM golang:1.5.3
MAINTAINER Michael Hahn <michael@lunohq.com>

RUN go get github.com/Masterminds/glide

ENV DOCKER_HOST unix:///var/run/docker.sock

ADD . /go/src/github.com/lunohq/relay
WORKDIR /go/src/github.com/lunohq/relay
RUN glide install
RUN GO15VENDOREXPERIMENT=1 go install ./cmd/relay

ENTRYPOINT ["/go/bin/relay"]
