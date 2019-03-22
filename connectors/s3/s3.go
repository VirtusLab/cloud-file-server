package s3

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/VirtusLab/cloud-file-server/config"
	"github.com/VirtusLab/cloud-file-server/connectors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bugsnag/bugsnag-go/errors"
)

const (
	prefixURI = "s3://"
)

type s3Handler struct {
	pathPrefix    string
	bucketName    string
	bucketFolders string
	service       *s3.S3
}

// New creates a connector that provides files from AWS s3 bucket
func New(config config.ConnectorConfig) (handler http.Handler, err error) {
	if len(config.URI) == 0 {
		return handler, errors.Errorf("URI parameter missing in connector: %#v", config)
	}
	if len(config.Region) == 0 {
		return handler, errors.Errorf("Region parameter missing in connector: %#v", config)
	}
	if !strings.HasPrefix(config.URI, prefixURI) {
		return handler, errors.Errorf("Invalid URI parameter in connector, expected '%s': %#v", prefixURI, config)
	}

	uriWithOutS3Prefix := strings.Replace(config.URI, prefixURI, "", 1)
	if uriWithOutS3Prefix == "" {
		return handler, errors.Errorf("Bucket name missing in URI parameter in connector: %#v", config)
	}

	bucketFolders := ""
	var bucketName string
	if index := strings.Index(uriWithOutS3Prefix, "/"); index == -1 {
		bucketName = uriWithOutS3Prefix
	} else {
		bucketName = uriWithOutS3Prefix[:index]
		bucketFolders = uriWithOutS3Prefix[index:]
	}

	awsConfig := aws.NewConfig().
		WithRegion(config.Region).
		WithCredentialsChainVerboseErrors(true)

	awsOptions := session.Options{
			Config: *awsConfig,
			Profile: config.Profile,
	}

	awsSession, err := session.NewSessionWithOptions(awsOptions)
	if err != nil {
		return nil, err
	}

	service := s3.New(awsSession)

	handler = &s3Handler{
		pathPrefix:    config.PathPrefix,
		bucketName:    bucketName,
		bucketFolders: bucketFolders,
		service:       service,
	}

	return handler, nil
}

func (h *s3Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	key := h.bucketFolders + strings.Replace(req.URL.Path, h.pathPrefix, "", 1)
	log.Printf("S3 request path %q, key %q", req.URL.Path, key)

	if strings.HasSuffix(key, "/") {
		http.Error(rw, connectors.PageNotFoundMessage, http.StatusNotFound)
		return
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(key),
	}
	if v := req.Header.Get("If-None-Match"); v != "" {
		input.IfNoneMatch = aws.String(v)
	}

	var is304 bool
	resp, err := h.service.GetObject(input)
	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case s3.ErrCodeNoSuchKey:
			http.Error(rw, connectors.PageNotFoundMessage, http.StatusNotFound)
			return
		case "NotModified":
			is304 = true
			// continue so other headers get set appropriately
		case "NoCredentialProviders":
			log.Printf("AWS Error: %v: %v", awsErr.Code(), awsErr.Message())
			http.Error(rw, connectors.UnauthorizedMessage, http.StatusUnauthorized)
			return
		case "AccessDenied":
			log.Printf("AWS Error: %v: %v", awsErr.Code(), awsErr.Message())
			http.Error(rw, connectors.ForbiddenMessage, http.StatusForbidden)
			return
		default:
			log.Printf("AWS Error: %v: %v", awsErr.Code(), awsErr.Message())
			http.Error(rw, connectors.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		log.Printf("non-AWS error %v %s", err, err)
		http.Error(rw, connectors.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	var contentType string
	ext := path.Ext(req.URL.Path)
	contentType = mime.TypeByExtension(ext)

	if contentType == "" {
		if resp.ContentType != nil {
			contentType = *resp.ContentType
		} else {
			log.Printf("Could not set conentType for %q", key)
			http.Error(rw, connectors.InternalServerErrorMessage, http.StatusInternalServerError)
		}
	}

	if resp.ETag != nil && *resp.ETag != "" {
		rw.Header().Set("Etag", *resp.ETag)
	}

	if contentType != "" {
		rw.Header().Set("Content-Type", contentType)
	}
	if resp.ContentLength != nil && *resp.ContentLength > 0 {
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))
	}

	if is304 {
		rw.WriteHeader(304)
	} else {
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Error occured during copy file %q, %q", key, err)
			http.Error(rw, connectors.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		_ = resp.Body.Close()
	}
}
