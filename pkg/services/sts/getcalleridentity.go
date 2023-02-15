package sts

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://docs.aws.amazon.com/STS/latest/APIReference/API_GetCallerIdentity.html
func (s *Service) GetCallerIdentity() (string, error) {
	body := url.Values{}
	body.Set("Action", "GetCallerIdentity")
	body.Set("Version", "2011-06-15")
	reqData := []byte(body.Encode())

	req, err := http.NewRequest("POST", "https://sts.amazonaws.com/", bytes.NewReader(reqData))
	if err != nil {
		return "", fmt.Errorf("creating request: %s", err)
	}

	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req, "us-east-1", "sts", reqData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body: %s", err)
	}

	var result GetCallerIdentityResponse
	xml.Unmarshal(data, &result)
	if result.Result.Arn == "" {
		return "", fmt.Errorf("no ARN in result")
	}

	return result.Result.Arn, nil
}

type GetCallerIdentityResponse struct {
	Result GetCallerIdentityResult `xml:"GetCallerIdentityResult"`
}

type GetCallerIdentityResult struct {
	Arn string `xml:"Arn"`
}
