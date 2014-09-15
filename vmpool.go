package main

import "net/http"
import "fmt"
import "io/ioutil"
import "strings"
import "os"
import "regexp"
import "bytes"
import "encoding/json"

var version string

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		matched, err := regexp.MatchString(a, b)
		perror(err)
		if matched {
			return true
		}
	}
	return false
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func list(url string, pattern string) {
	resp, err := http.Get(url)
	perror(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	perror(err)
	s := string(body[:])
	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	s = strings.Replace(s, `"`, "", -1)
	s = strings.Replace(s, `,`, "", -1)
	list := strings.Fields(s)
	// No pattern specified
	if pattern == "" {
		for i := 0; i < len(list); i++ {
			fmt.Println(list[i])
		}
	} else {
		pattern = strings.ToLower(pattern)
		matched := stringInSlice(pattern, list)
		if matched == false {
			// Invalid Pattern
			for i := 0; i < len(list); i++ {
				fmt.Println(list[i])
			}
		} else {
			var output []string
			for i := 0; i < len(list); i++ {
				matched, err := regexp.MatchString(pattern, list[i])
				perror(err)
				if matched == true {
					output = append(output, list[i])
				}
			}
			for i := 0; i < len(output); i++ {
				fmt.Println(output[i])
			}
		}
	}
	os.Exit(0)
}

func delete(url string, args string) {
	// Truncate the domain if required
	matched, err := regexp.MatchString("delivery.puppetlabs.net", args)
	if matched == true {
		foo := strings.Split(args, ".")
		args = foo[0]
	}
	uri := url + "/" + args
	f := "{}"
	input_json := []byte(f)
	body := bytes.NewBuffer(input_json)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", uri, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli")
	resp, err := client.Do(req)
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	var j interface{}
	err = json.Unmarshal(contents, &j)
	perror(err)
	m := j.(map[string]interface{})
	status := m["ok"]

	if status == true {
		fmt.Println(args + " deleted")
		os.Exit(0)
	} else {
		fmt.Println("Host " + args + " not found.")
		os.Exit(1)
	}
}

func grab(url string, args string) {
	uri := url + "/" + args
	f := "{}"
	input_json := []byte(f)
	body := bytes.NewBuffer(input_json)
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli")
	resp, err := client.Do(req)
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	var j interface{}
	err = json.Unmarshal(contents, &j)
	m := j.(map[string]interface{})
	ok := m["ok"]
	if ok == false {
		fmt.Println("Invalid pool name")
		os.Exit(1)
	}
	throwaway := (m[args])
	a := throwaway.(map[string]interface{})
	host := a["hostname"].(string)
	domainname := m["domain"].(string)
	perror(err)
	ret := args + ": " + host + "." + domainname
	fmt.Println(ret)
	os.Exit(0)
}

func usage() {
	fmt.Println("vmpool version: " + version + " - 2014")
	fmt.Println("vmpool list <pattern>")
	fmt.Println("vmpool grab <poolname>")
	fmt.Println("vmpool delete <hostname>")
	os.Exit(1)
}

func parseArgs(args []string) (c string, a string) {
	argument := ""
	valid_subcommands := []string{"list", "grab", "delete", "--help", "help", "version", "--version"}
	if len(args) < 2 {
		fmt.Println("vmpool requires a command")
		usage()
	}
	if !stringInSlice(args[1], valid_subcommands) {

		fmt.Println("invalid subcommand")
		usage()
	}

	if args[1] == "--help" || args[1] == "help" {
		usage()
	}

	if args[1] == "version" || args[1] == "--version" {
		usage()
	}

	if len(args) >= 3 {
		argument = args[2]
	}

	command := args[1]
	return command, argument
}

func main() {
	//  vlcoud base configuration can be overidden via ENV variable
	url := os.Getenv("VMPOOL_URL")
	if url == "" {
		url = "http://vcloud.delivery.puppetlabs.net/vm"
	} else {
		url = os.Getenv("VMPOOL_URL")
	}

	command, arguments := parseArgs(os.Args)
	if command == "list" {
		list(url, arguments)
	}

	if command == "grab" {
		grab(url, arguments)
	}

	if command == "delete" {
		delete(url, arguments)
	}

}
