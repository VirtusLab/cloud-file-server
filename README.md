# cloud-file-server

`cloud-file-server` is an application that serves files over HTTP using configured connectors (i.e. S3, file, dir)

[![Version](https://img.shields.io/badge/version-v0.0.5-brightgreen.svg)](https://github.com/VirtusLab/cloud-file-server/releases/tag/v0.0.5)
[![Build Status](https://travis-ci.org/VirtusLab/cloud-file-server.svg?branch=master)](https://travis-ci.org/VirtusLab/cloud-file-server)
[![Docker Repository on Quay](https://quay.io/repository/VirtusLab/cloud-file-server/status "Docker Repository on Quay")](https://quay.io/repository/VirtusLab/cloud-file-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/cloud-file-server)](https://goreportcard.com/report/github.com/VirtusLab/cloud-file-server)
[![GoDoc](https://godoc.org/github.com/VirtusLab/cloud-file-server?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/cloud-file-server)

## Connectors
- [AWS S3 bucket](https://aws.amazon.com/s3/)
- local directory
- local file

## Usage

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

## Operating system support

We provide cross-compiled binaries for most platforms, but is currently used mainly with `linux/amd64`.

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/cloud-file-server/issues) 
or [pull requests](https://github.com/VirtusLab/cloud-file-server/pulls).

## Development

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    mkdir -p $GOPATH/src/github.com/VirtusLab
    cd $GOPATH/src/github.com/VirtusLab
    git clone git@github.com:VirtusLab/cloud-file-server.git
    cd cloud-file-server
    
    go get -u github.com/golang/dep/cmd/dep
    make all

## The name

We believe in obvious names. It serves files from cloud. It's `cloud-file-server`.
