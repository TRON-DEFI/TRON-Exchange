package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	httpClient *http.Client
)

// init HTTPClient
func InitHttpClient() {
	if httpClient == nil {
		httpClient = createHTTPClient()
		return
	}
}

const (
	MaxIdleConns        = 300
	MaxIdleConnsPerHost = 300
	IdleConnTimeout     = 300
)

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   300 * time.Second,
				KeepAlive: 300 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},
		Timeout: 300 * time.Second,
	}
	return client
}

func Get(url, token string, tokenFlag bool) (string, error) {
	InitHttpClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("http.NewRequest error")
	}
	req.Header.Set("Content-Type", "application/json")
	if tokenFlag {
		req.Header.Set("Authorization", "Bearer " + token)
	}
	// use httpClient to send request
	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		fmt.Println("httpClient.Do error")
		return "", err
	} else {
		// Close the connection to reuse it
		defer response.Body.Close()
		// Let's check if the work actually is done
		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Couldn't parse response body")
			return "", err
		}
		return string(body), nil
	}
}

func Post(url string, data interface{}) (string, error) {
	InitHttpClient()
	value, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal error")
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(value)))
	if err != nil {
		fmt.Println("http.NewRequest error")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// use httpClient to send request
	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		fmt.Println("httpClient.Do error")
		return "", err
	} else {
		// Close the connection to reuse it
		defer response.Body.Close()
		// Let's check if the work actually is done
		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Couldn't parse response body")
			return "", err
		}

		return string(body), nil
	}
}
