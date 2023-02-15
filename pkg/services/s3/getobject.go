package s3

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (s *Service) GetObject(region string, bucket string, key string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %s", err)
	}

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
		return nil, fmt.Errorf("reading body: %s", err)
	}

	return data, nil
}
