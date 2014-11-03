package main

import (
	"fmt"
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

	fmt.Printf("vmpool version: %s - 2014\n", version)
}
