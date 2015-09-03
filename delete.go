package main

import (
	"fmt"
	"log"
	"os"
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
		shortname := shortname(arg)
		_, output_json := RequestWrapper("delete", shortname, "DELETE", "{}")

		if output_json["ok"] == false {
			fmt.Println("Host not found or token invalid.")
			log.Printf("Host %s not found or token invalid.\n", shortname)
			os.Exit(1)
		} else {
			fmt.Printf("%s deleted\n", shortname)
			logmsg(fmt.Sprintf("%s deleted\n", shortname))
		}
	}
}
