package main

import (
	"io/ioutil"
)

var cmdSummary = &Command{
	Run:       runSummary,
	UsageLine: "summary",
	Short:     "Get detailed summary of the vmpooler",
	Long: `
Display summary information for the vmpooler.
  Warning: This can be long and verbose.
    `,
}

func runSummary(cmd *Command, args []string) {
	debug("Fucntion: runSummary")
	resp, err := Request("summary", "GET", "", "{}")
	perror(err)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	perror(err)
	pjson(contents)
}
