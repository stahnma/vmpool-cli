package main

import (
	"bytes"
	"net/http"
)

func Request(method, arg, input_json string) (resp *http.Response, err error) {
	url := vmpool_url + "/" + arg
	body := bytes.NewBuffer([]byte(input_json))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli")
	return client.Do(req)
}
