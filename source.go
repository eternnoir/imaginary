package main

import (
	"net/http"
	"net/url"
)

type ImageSourceType string
type ImageSourceFactoryFunction func(*SourceConfig) ImageSource

type SourceConfig struct {
	AuthForwarding  bool
	Authorization   string
	MountPath       string
	Type            ImageSourceType
	AllowedOrigings []*url.URL
	MaxAllowedSize  int

	EnableS3    bool
	S3AccessKey string
	S3Secret    string
	S3Bucket    string
	S3Region    string
	S3Path      string
}

var imageSourceMap = make(map[ImageSourceType]ImageSource)
var imageSourceFactoryMap = make(map[ImageSourceType]ImageSourceFactoryFunction)

type ImageSource interface {
	Matches(*http.Request) bool
	GetImage(*http.Request) ([]byte, error)
}

func RegisterSource(sourceType ImageSourceType, factory ImageSourceFactoryFunction) {
	imageSourceFactoryMap[sourceType] = factory
}

func LoadSources(o ServerOptions) {
	for name, factory := range imageSourceFactoryMap {
		imageSourceMap[name] = factory(&SourceConfig{
			Type:            name,
			MountPath:       o.Mount,
			AuthForwarding:  o.AuthForwarding,
			Authorization:   o.Authorization,
			AllowedOrigings: o.AlloweOrigins,
			MaxAllowedSize:  o.MaxAllowedSize,
			EnableS3:        o.EnableS3,
			S3AccessKey:     o.S3AccessKey,
			S3Secret:        o.S3Secret,
			S3Bucket:        o.S3Bucket,
			S3Region:        o.S3Region,
			S3Path:          o.S3Path,
		})
	}
}

func MatchSource(req *http.Request) ImageSource {
	for _, source := range imageSourceMap {
		if source.Matches(req) {
			return source
		}
	}
	return nil
}
