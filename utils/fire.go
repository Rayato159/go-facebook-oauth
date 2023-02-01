package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func httpMethodCheck(method string) error {
	methods := map[string]string{
		"GET":    "GET",
		"POST":   "POST",
		"PUT":    "PUT",
		"PATCH":  "PATCH",
		"DELETE": "DELETE",
	}

	if methods[method] == "" || methods[method] != method {
		return fmt.Errorf("unknown http method: %s", method)
	}
	return nil
}

func FireHttpRequest(method, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*120))
	defer cancel()

	// Http method validation
	if err := httpMethodCheck(method); err != nil {
		return nil, err
	}

	// Config before request
	config := http.Header{
		"Content-Type": {"application/json"},
	}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request config error: %v", err)
	}
	req.Header = config
	client := new(http.Client)

	// Request fire!!!
	res, err := client.Do(req)
	if err != nil {
		defer func() {
			req.Close = true
			res.Body.Close()
		}()
		return nil, fmt.Errorf("http request error: %v", err)
	}
	defer func() {
		req.Close = true
		res.Body.Close()
	}()

	// Body response
	resJson, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("response error: %v", err)
	}
	return resJson, nil
}
