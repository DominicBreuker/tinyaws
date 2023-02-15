package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const longTimeFormat = "20060102T150405Z"
const shortTimeFormat = "20060102"
const algorithm = "AWS4-HMAC-SHA256"

var newLine = []byte{'\n'}

func signRequest(req *http.Request, accessKey, secretKey, sessionToken, region, service string, data []byte) {
	t := time.Now().UTC()

	setXAmzDateHeader(req, t)
	setXAmzContentSHA256Header(req, data)
	if sessionToken != "" {
		setXAmzSecurityTokenHeader(req, sessionToken)
	}

	scope := getScope(t, region, service)

	signingKey := getSigningKey(t, secretKey, service, region)
	stringToSign := getStringToSign(t, scope, req, data)
	signature := fmt.Sprintf("%x", HMACSHA256(signingKey, stringToSign.Bytes()))

	setAuthorizationHeader(req, accessKey, scope, signature)
}

func setXAmzDateHeader(req *http.Request, t time.Time) {
	req.Header.Set("X-Amz-Date", t.Format(longTimeFormat))
}

func setXAmzContentSHA256Header(req *http.Request, data []byte) {
	value := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // default for no content

	if len(data) > 0 {
		value = fmt.Sprintf("%x", SHA256(data))
	}

	req.Header.Set("x-amz-content-sha256", value)
}

func setXAmzSecurityTokenHeader(req *http.Request, sessionToken string) {
	req.Header.Set("X-Amz-Security-Token", sessionToken)
}

func getScope(t time.Time, region, service string) string {
	return t.Format(shortTimeFormat) + "/" + region + "/" + service + "/aws4_request"
}

func getSigningKey(t time.Time, secretKey, service, region string) []byte {
	h := HMACSHA256([]byte("AWS4"+secretKey), []byte(t.Format("20060102")))
	h = HMACSHA256(h, []byte(region))
	h = HMACSHA256(h, []byte(service))
	h = HMACSHA256(h, []byte("aws4_request"))
	return h
}

func getStringToSign(t time.Time, scope string, req *http.Request, data []byte) *bytes.Buffer {
	w := &bytes.Buffer{}

	w.Write([]byte(algorithm))
	w.Write(newLine)

	w.Write([]byte(t.Format(longTimeFormat)))
	w.Write(newLine)

	w.Write([]byte(scope))
	w.Write(newLine)

	fmt.Fprintf(w, "%x", SHA256(getCanonicalRequest(req, data).Bytes()))

	return w
}

func getCanonicalRequest(req *http.Request, data []byte) *bytes.Buffer {
	w := &bytes.Buffer{}

	writeHTTPMethod(w, req)
	w.Write(newLine)

	writeCanonicalURI(w, req)
	w.Write(newLine)

	writeCanonicalQueryString(w, req)
	w.Write(newLine)

	writeCanonicalHeaders(w, req)
	w.Write(newLine)
	w.Write(newLine)

	writeSignedHeaders(w, req)
	w.Write(newLine)

	writeHashedPayload(w, data)

	return w
}

func setAuthorizationHeader(req *http.Request, accessKey, scope, signature string) {
	authorization := &bytes.Buffer{}

	authorization.Write([]byte(algorithm + " Credential=" + accessKey + "/" + scope))
	authorization.Write([]byte{',', ' '})

	authorization.Write([]byte("SignedHeaders="))
	writeSignedHeaders(authorization, req)
	authorization.Write([]byte{',', ' '})

	authorization.Write([]byte("Signature=" + signature))

	req.Header.Set("Authorization", authorization.String())
}

func writeHTTPMethod(w io.Writer, req *http.Request) {
	w.Write([]byte(req.Method))
}

func writeCanonicalURI(w io.Writer, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}

func writeCanonicalQueryString(w io.Writer, r *http.Request) {
	var params []string

	for name, values := range r.URL.Query() {
		for _, value := range values {
			if value == "" {
				params = append(params, url.QueryEscape(name))
			} else {
				params = append(params, url.QueryEscape(name)+"="+url.QueryEscape(value))
			}
		}
	}

	sort.Strings(params)

	w.Write([]byte(strings.Join(params, "&")))
}

func writeCanonicalHeaders(w io.Writer, r *http.Request) {
	var headers []string

	for name, value := range r.Header {
		sort.Strings(value)
		headers = append(headers, strings.ToLower(name)+":"+strings.Join(value, ","))
	}

	sort.Strings(headers)

	w.Write([]byte(strings.Join(headers, "\n")))
}

func writeSignedHeaders(w io.Writer, r *http.Request) {
	var headers []string

	for name := range r.Header {
		headers = append(headers, strings.ToLower(name))
	}

	sort.Strings(headers)

	w.Write([]byte(strings.Join(headers, ";")))
}

func writeHashedPayload(w io.Writer, data []byte) {
	fmt.Fprintf(w, "%x", SHA256(data))
}

func SHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
