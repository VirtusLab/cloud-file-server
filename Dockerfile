FROM alpine:3.12

RUN apk --no-cache upgrade && \
    apk --no-cache add --update ca-certificates

ADD cross/cloud-file-server-linux-amd64 /usr/local/bin/cloud-file-server

EXPOSE 8080

ENTRYPOINT ["cloud-file-server"]
