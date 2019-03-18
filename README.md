# cloud-file-server

`cloud-file-server` is an application that serves files over HTTP using configured connectors (i.e. S3, file, dir)

[![Build Status](https://travis-ci.org/VirtusLab/cloud-file-server.svg?branch=master)](https://travis-ci.org/VirtusLab/cloud-file-server)
[![Docker Repository on Quay](https://quay.io/repository/VirtusLab/cloud-file-server/status "Docker Repository on Quay")](https://quay.io/repository/VirtusLab/cloud-file-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/cloud-file-server)](https://goreportcard.com/report/github.com/VirtusLab/cloud-file-server)
[![GoDoc](https://godoc.org/github.com/VirtusLab/cloud-file-server?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/cloud-file-server)

## Connectors
- [AWS S3 bucket](https://aws.amazon.com/s3/)
- local directory
- local file

## Run

    ./cloud-file-server --config example-config.yaml
    
## Example config

    listen: :8080
    logRequests: true
    connectors:
    - type: s3
      uri: s3://aws-s3-bucket-name/example/path
      region: eu-west-1
      pathPrefix: /s3
    - type: file
      uri: file:///example/path/file.yaml
      pathPrefix: /file
    - type: directory
      uri: file:///example/path/directory
      pathPrefix: /dir
