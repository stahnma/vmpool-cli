package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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
	resp, err := Request("grab", "POST", strings.Join(args, "+"), "{}")
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(2)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	perror(err)
	var j map[string]interface{}
	err = json.Unmarshal(contents, &j)
	perror(err)
	if j["ok"] == false {
		if len(args) > 1 {
			// Check which of the pools they gave aren't valid
			invalidMap := make(map[string]bool)
			for _, arg := range args {
				invalidMap[arg] = true
			}
			pools := Vmpools()
			for _, pool := range pools {
				if invalidMap[pool] {
					invalidMap[pool] = false
				}
			}
			invalidPools := make([]string, 0)
			for pool := range invalidMap {
				if invalidMap[pool] {
					invalidPools = append(invalidPools, pool)
				}
			}
			log.Printf("Invalid pool name(s): [ %s ]\n", strings.Join(invalidPools, ", "))
		} else {
			log.Printf("Invalid pool name: %s\n", args[0])
		}
		log.Printf("Something went wrong with grab. Perhaps an invalid token?")
		fmt.Println("Something went wrong with grab. Perhaps an invalid token?")
		os.Exit(2)
	}
	for _, arg := range unique(args) {
		var vm string
		a := j[arg].(map[string]interface{})
		hosts := a["hostname"]
		switch hosts.(type) {
		case string:
			host := hosts.(string)
			vm = host
			fmt.Printf("%v: %v.%v\n", arg, host, j["domain"])
			logmsg(fmt.Sprintf("Retrieved %v from pool %v\n", host, args[0]))
		case []interface{}:
			fmt.Printf("%v:\n", arg)
			hosts := hosts.([]interface{})
			for _, host := range hosts {
				host := host.(string)
				vm = shortname(host)
				fmt.Printf("    %v.%v\n", host, j["domain"])
				logmsg(fmt.Sprintf("Retrieved %v from pool %v\n", host, args[0]))
			}
		}
		appendTags(shortname(vm))
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
