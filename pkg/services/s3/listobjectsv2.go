package s3

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
func (s *Service) ListObjectsV2(region, bucket, prefix string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucket, region), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %s", err)
	}

	q := req.URL.Query()
	q.Add("list-type", "2")
	q.Add("max-keys", "10") // change if you like, max 1000
	if prefix != "" {
		q.Add("prefix", prefix)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.Client.Do(req, region, "s3", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading HTTP response body: %s", err)
	}

	var result S3ListBucketResult
	xml.Unmarshal(data, &result)

	var out []string
	for _, key := range result.Contents.Keys {
		out = append(out, key)
	}

	return out, nil
}

type S3ListBucketResult struct {
	KeyCount    int              `xml:"KeyCount"`
	MaxKeys     int              `xml:"MaxKeys"`
	IsTruncated bool             `xml:"IsTruncated"`
	Contents    S3BucketContents `xml:"Contents"`
}

type S3BucketContents struct {
	Keys []string `xml:"Key"`
}
