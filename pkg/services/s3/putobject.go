package s3

import (
	"bytes"
	"fmt"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (s *Service) PutObject(region string, bucket string, key string, data []byte) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := s.Client.Do(req, region, "s3", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
