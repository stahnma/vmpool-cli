package main

import "net/http"
import "fmt"
import "io/ioutil"
import "strings"
import "os"
import "log"
import "bytes"
import "encoding/json"

var cmdDelete = &Command{
	Run:       runDelete,
	UsageLine: "delete [hostname]",
	Short:     "delete a vm from vmpooler",
	Long: `
Delete removes a vm with the given hostname
from vmpooler.
    `,
}

func runDelete(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
	}
	arg := args[0]
	// Truncate the domain if required
	if matchPattern("delivery.puppetlabs.net")(arg) {
		arg = strings.Split(arg, ".")[0]
	}
	url := vmpool_url + "/" + arg
	input_json := []byte("{}")
	body := bytes.NewBuffer(input_json)
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, body)
	perror(err)
	req.Header.Add("User-Agent", "vmpool-cli")
	resp, err := client.Do(req)
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	var j map[string]interface{}
	err = json.Unmarshal(contents, &j)
	perror(err)
	status := j["ok"]

	if status == false {
		log.Printf("Host %s not found.\n", arg)
		os.Exit(1)
	}

	fmt.Printf("%s deleted\n", arg)
}
