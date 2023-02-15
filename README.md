# tinyaws

A PoC for an AWS API client based only on the Go standard library, without the bloated official AWS SDK.
If you need only a few API calls and don't want to blow up the size of your binary, this may be for you.

This code will be little more than a blueprint.
It implements only a handful of API actions for STS and S3.
Use them as examples to implement the actions you actually need.
The authentication logic (AWS signature) though should be reusable.

## Getting Started

Check out the file [main.go](main.go), which shows all implemented API actions in action.
They are:
- STS: GetCallerIdentity
- S3: PutObject
- S3: GetObject
- S3: ListObjectsV2
- S3: DeleteObject

To make it work you have to add IAM credentials into the variables first:

```go
func main() {
	accessKey := "YOUR_KEY_ID"
	secretKey := "YOUR_SECRET_KEY"

	TestSTS(accessKey, secretKey)
	TestS3(accessKey, secretKey)
}
```

To test the S3 actions you need a bucket and write permissions for it.
Configure the name on the top of the `TestS3` function:

```go
func TestS3(accessKey string, secretKey string) {
	region := "us-east-2"
	bucket := "fancy-bucket-name"
	key := "somefile.txt"
	data := []byte("your ad could be here")
    ...
}
```

Then run `go run main.go` and you should see some log output related to the API actions.

## Implementing API Actions

AWS provides detailed documentation for each API action.
All you have to do is build the request using Go's [net/http Request](https://pkg.go.dev/net/http#Request).
The tricky part is authentication.
AWS wants you to create a signature for each request, with your credentials used as a secret key.
Details are [in the docs](https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html).
This repo provides an implementation of that in [pkg/client/signature.go](pkg/client/signature.go).
The function is called `signRequest`.

For example, consider the implementation of S3 GetObject.
The code for it is in [pkg/services/s3/getobject.go](pkg/services/s3/getobject.go).
AWS documentation is [here](https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html).
This is the code, defined of the S3 `Service`:

```go
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
```

Nothing special to see.
You could add headers or query parameters if needed, or leave it like this if you are fine with it.
All authentication happens under the hood.
The S3 service has a `Client` which wraps a [net/http Client](https://pkg.go.dev/net/http#Client)
but adds the signature before doing the request.

You should be able to add more API actions in a similar fashion, without having to deal with the signatures.
However, understand that all request data (headers, URI, parameters, body) is part of the signature and AWS is very fussy about these.
Your request must look exactly the same way the AWS docs describe it.
Expect that my implementation of the signature scheme won't resolve all ambiguities for you.

