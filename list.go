package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list [pattern]",
	Short:     "list the available platforms on vmpooler",
	Long: `
List queries vmpooler for the available platforms
and returns a list matching the given pattern
or the whole list if no pattern is specified.
    `,
}

func Vmpools() []string {
	resp, err := Request("list", "GET", "", "{}")
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	perror(err)
	var pools []string
	err = json.Unmarshal(contents, &pools)
	perror(err)
	return pools
}

func runList(cmd *Command, args []string) {
	pools := Vmpools()
	if len(args) != 1 && len(args) != 0 {
		cmd.Usage()
	}
	if len(args) == 1 {
		pattern := args[0]
		pattern = strings.ToLower(pattern)
		pools = filterStrings(pools, pattern)
	}
	printStrings(pools)
	if len(pools) == 0 {
		os.Exit(1)
	}
}
