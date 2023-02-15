package s3

import (
	"fmt"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (s *Service) DeleteObject(region string, bucket string, key string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key), nil)
	if err != nil {
		return fmt.Errorf("creating request: %s", err)
	}

	resp, err := s.Client.Do(req, region, "s3", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
