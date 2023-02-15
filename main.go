package main

import (
	"dominicbreuker/tinyaws/pkg/services/s3"
	"dominicbreuker/tinyaws/pkg/services/sts"
	"log"
)

func main() {
	accessKey := "YOUR_KEY_ID"
	secretKey := "YOUR_SECRET_KEY"

	TestSTS(accessKey, secretKey)
	TestS3(accessKey, secretKey)
}

func TestSTS(accessKey string, secretKey string) {
	svc := sts.NewService(accessKey, secretKey, "")

	arn, err := svc.GetCallerIdentity()
	if err != nil {
		log.Fatalf("Error: GetCallerIdentity: %s", err)
	}
	log.Printf("Your ARN is %s\n", arn)
}

func TestS3(accessKey string, secretKey string) {
	region := "us-east-2"
	bucket := "fancy-bucket-name"
	key := "somefile.txt"
	data := []byte("your ad could be here")

	svc := s3.NewService(accessKey, secretKey, "")

	log.Printf("Uploading file: s3://%s/%s", bucket, key)
	if err := svc.PutObject(region, bucket, key, data); err != nil {
		log.Fatalf("Error: PutObject: %s", err)
	}

	log.Printf("Downloading file: s3://%s/%s", bucket, key)
	data, err := svc.GetObject(region, bucket, key)
	if err != nil {
		log.Fatalf("Error: GetObject: %s", err)
	}
	log.Printf("  Content of file: %+v\n", string(data))

	log.Printf("Listing files: s3://%s/*", bucket)
	keys, err := svc.ListObjectsV2(region, bucket, "")
	if err != nil {
		log.Fatalf("Error: ListObjectsV2: %s", err)
	}
	for _, key := range keys {
		log.Printf("  - %s\n", key)
	}

	log.Printf("Deleting file: s3://%s/%s", bucket, key)
	if err := svc.DeleteObject(region, bucket, key); err != nil {
		log.Fatalf("Error: DeleteObject: %s", err)
	}
}
