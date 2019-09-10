package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func newRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {

	requestURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Parse endpoint ERROR:")
	}

	return http.NewRequest(method, requestURL.String(), body)
}

func newClient() (*http.Client, error) {

	client := &http.Client{}

	if config().proxyEnabled() {

		proxyStr := config().proxyURL()
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, errors.Wrap(err, "Parse proxy url ERROR:")
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		//adding the Transport object to the http Client
		client.Transport = transport
	}

	return client, nil
}
