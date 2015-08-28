package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func RequestWrapper(method string, params string, http_action string, input_json string) ([]byte, map[string]interface{}) {
	debug("Function RequestWrapper")
	debug("  Method: " + method)
	debug("  Params: " + params)
	debug("  http_action: " + http_action)
	debug("  input_json: " + input_json)

	resp, err := Request(method, http_action, params, input_json)
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	perror(err)
	var output_json map[string]interface{}
	err = json.Unmarshal(contents, &output_json)
	perror(err)
	return contents, output_json
}

func Request(task, http_action, arg, input_json string) (resp *http.Response, err error) {
	debug("Function Request")
	debug("  Task:" + task)
	debug("  Method:" + http_action)
	debug("  Arg: " + fmt.Sprintf("%v", arg))
	debug("  Input_JSON: " + input_json)
	var url string
	var vm string
	var token string
	if task != "token" {
		token = retrieveToken()
	}
	switch {
	case task == "token":
		url = vmpool_url + "/api/v1/token"
		if http_action == "DELETE" {
			url = vmpool_url + "/api/v1/token/" + strings.Split(arg, "|")[2]
		}
	case task == "status":
		url = vmpool_url + "/api/v1/status"
	case task == "summary":
		url = vmpool_url + "/api/v1/summary"
	case task == "vm":

		params := strings.Split(arg, "|")
		if http_action == "PUT" {
			token = params[0]
			vm = params[1]
		}
		if http_action == "GET" {
			vm = strings.Split(arg, "|")[0]
		}
		url = vmpool_url + "/api/v1/vm/" + vm
	default:
		url = vmpool_url + "/api/v1/vm/" + arg
	}
	debug("  Using " + http_action + " at " + url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	body := bytes.NewBuffer([]byte(input_json))
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http_action, url, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli-"+version)
	if task == "token" {
		s := strings.Split(arg, "|")
		req.SetBasicAuth(s[0], s[1])
	}
	if token != "" {
		debug("  Adding X-AUTH-TOKEN header " + token)
		req.Header.Add("X-AUTH-TOKEN", token)
	}

	return client.Do(req)
}
