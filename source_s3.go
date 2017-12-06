package main

import (
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const ImageSourceTypeS3 ImageSourceType = "S3"

type S3ImageSource struct {
	Config *SourceConfig
	creds  *credentials.Credentials
}

func NewS3ImageSource(config *SourceConfig) ImageSource {
	creds := credentials.NewStaticCredentials(config.S3AccessKey, config.S3Secret, "")
	return &S3ImageSource{Config: config, creds: creds}
}

func (s *S3ImageSource) Matches(r *http.Request) bool {
	return r.Method == "GET" && r.URL.Query().Get("s3") != ""
}

func (s *S3ImageSource) GetImage(r *http.Request) ([]byte, error) {
	filename := r.URL.Query().Get("s3")
	key := path.Join(s.Config.S3Path, filename)
	cfg := aws.NewConfig().WithRegion(s.Config.S3Region).WithCredentials(s.creds)
	buff := &aws.WriteAtBuffer{}
	s3dl := s3manager.NewDownloader(session.New(cfg))
	_, err := s3dl.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(s.Config.S3Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func init() {
	RegisterSource(ImageSourceTypeS3, NewS3ImageSource)
}
