package s3

import "dominicbreuker/tinyaws/pkg/client"

type Service struct {
	Client *client.Client
}

func NewService(accessKey string, secretKey string, sessionToken string) *Service {
	return &Service{
		Client: client.NewClient(accessKey, secretKey, sessionToken),
	}
}
