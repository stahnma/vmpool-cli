package main

import (
	"io/ioutil"
)

var cmdStatus = &Command{
	Run:       runStatus,
	UsageLine: "status",
	Short:     "Get basic health of the vmpooler",
	Long: `
Display vmpooler health information via the status endpoint.
    `,
}

func runStatus(cmd *Command, args []string) {
	debug("Function: runStatus")
	resp, err := Request("status", "GET", "", "{}")
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	perror(err)
	pjson(contents)
}
