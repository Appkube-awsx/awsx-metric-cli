package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	contentTypeJSON           = "application/json"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

const (
	headerKeyAccept        = "Accept"
	headerKeyContentType   = "Content-Type"
	headerKeyAuthorization = "Authorization"
	headerKeyIdToken       = "X-ID-Token"
)

func GetRequest(url string) (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodGet, url, nil)
	return req, err
}

func ProcessApiCall(url string, allowedHosts []string) (obj any, statusCode int, duration time.Duration, err error) {
	req, err := GetRequest(url)
	if err != nil {
		fmt.Println("creating http get request failed. error: ", err.Error())
		return
	}
	startTime := time.Now()
	if !CanAllowURL(req.URL.String(), allowedHosts) {
		fmt.Println("url is not in the allowed list. make sure to match the base URL with the allowed list", "url", req.URL.String())
		return nil, http.StatusUnauthorized, 0, fmt.Errorf("requested URL is not allowed. To allow this URL, add this URL to Allowed Hosts section")
	}
	client := &http.Client{}
	res, err := client.Do(req)
	duration = time.Since(startTime)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil && res != nil {
		fmt.Println("error getting response from server", "url", url, "method", req.Method, "error", err.Error(), "status code", res.StatusCode)
		return nil, res.StatusCode, duration, fmt.Errorf("error getting response from %s", url)
	}
	if err != nil && res == nil {
		fmt.Println("error getting response from server. no response received", "url", url, "error", err.Error())
		return nil, http.StatusInternalServerError, duration, fmt.Errorf("error getting response from url %s. no response received. Error: %s", url, err.Error())
	}
	if err == nil && res == nil {
		fmt.Println("invalid response from server and also no error", "url", url, "method", req.Method)
		return nil, http.StatusInternalServerError, duration, fmt.Errorf("invalid response received for the URL %s", url)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, res.StatusCode, duration, fmt.Errorf(res.Status)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body", "url", url, "error", err.Error())
		return nil, res.StatusCode, duration, err
	}
	bodyBytes = removeBOMContent(bodyBytes)
	if CanParseAsJSON(res.Header) {
		var out any
		err := json.Unmarshal(bodyBytes, &out)
		if err != nil {
			fmt.Println("error un-marshaling JSON response", "url", url, "error", err.Error())
		}
		return out, res.StatusCode, duration, err
	}
	return string(bodyBytes), res.StatusCode, duration, err
}

func CanParseAsJSON(responseHeaders http.Header) bool {
	contentType := responseHeaders.Get(headerKeyContentType)
	if strings.Contains(strings.ToLower(contentType), contentTypeJSON) {
		return true
	}
	return false
}

func removeBOMContent(input []byte) []byte {
	return bytes.TrimPrefix(input, []byte("\xef\xbb\xbf"))
}
func CanAllowURL(url string, allowedHosts []string) bool {
	allow := false
	if len(allowedHosts) == 0 {
		return true
	}
	for _, host := range allowedHosts {
		if strings.HasPrefix(url, host) {
			return true
		}
	}
	return allow
}
