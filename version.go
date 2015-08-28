package main

import (
	"fmt"
	"strconv"
	"time"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print Vmpool version",
	Long:      `Version prints the Vmpool version.`,
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
	}

	t := time.Now()
	year := strconv.Itoa(t.Year())
	fmt.Printf("vmpool version: %s - "+year+"\n", version)
}
