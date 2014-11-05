package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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
	for _, arg := range args {
		// Truncate the domain if required
		if matchPattern("delivery.puppetlabs.net")(arg) {
			arg = strings.Split(arg, ".")[0]
		}
		resp, err := Request("DELETE", arg, "{}")
		if err != nil {
			log.Printf("%v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		perror(err)
		var output_json map[string]bool
		err = json.Unmarshal(contents, &output_json)
		perror(err)
		if !output_json["ok"] {
			log.Printf("Host %s not found.\n", arg)
			os.Exit(1)
		} else {
			fmt.Printf("%s deleted\n", arg)
		}
	}
}
