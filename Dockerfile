FROM golang:1.12-alpine

RUN apk add --no-cache --update alpine-sdk

COPY . /go/src/github.com/VirtusLab/cloud-file-server
RUN cd /go/src/github.com/VirtusLab/cloud-file-server && go build

FROM alpine:3.9

EXPOSE 8080

RUN apk --no-cache upgrade && \
    apk --no-cache add --update ca-certificates bash

COPY --from=0 /go/src/github.com/VirtusLab/cloud-file-server/cloud-file-server /usr/local/bin/cloud-file-server

ENTRYPOINT ["cloud-file-server"]