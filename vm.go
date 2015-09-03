package main

import (
	"fmt"
	"log"
	"os"
)

var cmdVm = &Command{
	Run:       runVm,
	UsageLine: "vm",
	Short:     "work with vms from vmpooler",
	Long: `

    vm grab <pool name>:
      - request a vm from pool
    vm delete <vm>:
      - delete a vm
    vm info <vm>:
       - retreive info about a specific vm
    vm lifetime <vm> <TTL in hours>
    `,
}

func vmTag(hostname string, tags string) {
	debug("Funciton vmTag")
	debug("  Hostname: " + hostname)
	debug("  Tags: " + tags)
	token := retrieveToken()
	params := token + "|" + hostname
	_, output_json := RequestWrapper("vm", params, "PUT", tags)
	if output_json["ok"] == true {
		debug(fmt.Sprintf("Tags changed to " + tags + " on " + hostname))
	} else {
		fmt.Println("Something went wrong. Unable to adjust Tags.")
		os.Exit(255)
	}
}

func vmLifetime(hostname string, ttl string) {
	debug("Function: vmLifetime")
	debug("  Hostname: " + hostname)
	debug("  TTL: " + ttl)
	token := retrieveToken()
	params := token + "|" + hostname
	lifetime := "{\"lifetime\":" + ttl + "}"
	input_json := lifetime
	_, output_json := RequestWrapper("vm", params, "PUT", input_json)
	if output_json["ok"] == true {
		fmt.Println("Lifetime changed to " + ttl + " on " + hostname)
	} else {
		fmt.Println("Something went wrong. Unable to adjust TTL.")
		os.Exit(255)
	}
}

func appendTags(vm string) {
	debug("Function: appendTags")
	user := ""
	if os.Getenv("LDAP_USERNAME") != "" {
		user = os.Getenv("LDAP_USERNAME")
	} else {
		user = os.Getenv("USER")
	}
	tags := `{"tags": {"user":"` + user + `", "client":"vmpool-cli-` + version + `"}}`
	debug("Tags:" + tags)
	vmTag(vm, tags)
}

func runVm(cmd *Command, args []string) {
	debug("Function: runVM")
	if len(args) < 1 {
		log.Fatal("You need arguments.")
	}
	debug("  ARGS for runVM are" + fmt.Sprintf("%v", args))
	var subcmd string
	if len(args) > 0 {
		subcmd = args[0]
	} else {
		subcmd = ""
	}

	var params string
	var http_action string
	var vm string
	if len(args) > 1 {
		params = shortname(args[1])
	}
	if len(args) < 2 {
		log.Fatal("info requires a vm name argument.")
	} else {
		vm = shortname(args[1])
	}
	switch {
	case (subcmd == "info"):
		http_action = "GET"
		// how do I know args is correct?
		contents, output_json := RequestWrapper("vm", params, http_action, "{}")
		if output_json["ok"] == false {
			fmt.Println("VM not found.")
			os.Exit(255)
		}
		pjson(contents)
		os.Exit(0)
	case (subcmd == "delete"):
		host := []string{vm}
		runDelete(cmdDelete, host)
		os.Exit(0)
	case (subcmd == "grab"):
		host := []string{vm}
		runGrab(cmdGrab, host)
		os.Exit(0)

	case (subcmd == "lifetime"):
		ttl := ""
		if len(args) < 3 {
			log.Fatal("lifetime action requires a TTL.")
		} else {
			ttl = args[2]
		}
		vmLifetime(vm, ttl)
	default:
		log.Fatal("Unkonwn argument to vm.")
	}
}
