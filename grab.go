package main

import "net/http"
import "fmt"
import "log"
import "io/ioutil"
import "os"
import "bytes"
import "strings"
import "encoding/json"

var cmdGrab = &Command{
	Run:       runGrab,
	UsageLine: "grab [pool...]",
	Short:     "check-out a vm or vms",
	Long: `
Grab posts to vmpooler to fetch a vm from the given pool
and returns the domain name of the vm.
    `,
}

func runGrab(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
	}
	url := vmpool_url + "/" + strings.Join(args, "+")
	input_json := []byte("{}")
	body := bytes.NewBuffer(input_json)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli")
	resp, err := client.Do(req)
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	var j map[string]interface{}
	err = json.Unmarshal(contents, &j)
	perror(err)
	ok := j["ok"]
	if ok == false {
		log.Printf("Invalid pool name in: [ %s ]\n", strings.Join(args, ", "))
		os.Exit(1)
	}
	for _, arg := range unique(args) {
		a := j[arg].(map[string]interface{})
		b := a["hostname"]
		switch b.(type) {
		case string:
			host := b.(string)
			fmt.Printf("%v: %v.%v\n", arg, host, j["domain"])
		case []interface{}:
			b := b.([]interface{})
			hosts := interfacesToStrings(b)
			suffix := fmt.Sprintf(".%s\n    ", j["domain"])
			hostsString := fmt.Sprintf("\n    %s.%s\n", strings.Join(hosts, suffix), j["domain"])
			fmt.Printf("%v:%s", arg, hostsString)
		}
	}
}

func unique(args []string) []string {
	m := make(map[string]bool)
	for _, arg := range args {
		m[arg] = true
	}
	uniques := make([]string, 0)
	for key := range m {
		uniques = append(uniques, key)
	}
	return uniques
}

func interfacesToStrings(list []interface{}) []string {
	output := make([]string, len(list))
	for i, elm := range list {
		output[i] = elm.(string)
	}
	return output
}
