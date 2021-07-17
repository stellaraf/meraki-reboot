package util

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"stellar.af/meraki-reboot/types"
)

const requestTimeout int = 30

var httpClient *http.Client
var headers map[string]string

func CreateHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{Jar: jar, Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(requestTimeout) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}, Timeout: time.Duration(requestTimeout) * time.Second}

	return client
}

func CreateInsecureHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{Jar: jar, Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}, Timeout: time.Duration(requestTimeout) * time.Second}

	return client
}

// Perform an HTTP request and return the body as a string.
func Request(m string, u string, p ...string) (b string, e error) {
	// Ensure the base URL has a trailing slash.
	if string(u[len(u)-1:]) != "/" {
		u = u + "/"
	}
	reqURL, err := url.Parse(u + strings.Join(p, "/"))

	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest(m, reqURL.String(), nil)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("Request to %s failed:\n%s", reqURL, err.Error())
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error requesting data from %s - %s", reqURL, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", fmt.Errorf("Unable to parse response from %s\n%s", reqURL, err.Error())
	}

	return string(body), nil
}

// BuildUrl creates a "guaranteed correct" URL from a base URL and path.
func BuildUrl(base string, path string, query types.QueryParams) (u *url.URL) {
	// Create a URL object from the base URL.
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	// Create an array of all path elements.
	pathParts := SplitRemoveEmpty(path, "/")
	// Ensure the final URL has a trailing slash.
	// Override the URL object's path property with the cleaned path.
	u.Path = strings.Join(pathParts, "/")
	// Add query params.
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%s", v))
	}
	u.RawQuery = q.Encode()
	return u
}
