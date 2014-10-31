package main

import (
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

func runList(cmd *Command, args []string) {
	resp, err := Request("GET", "", "{}")
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	s := string(body[:])
	for _, char := range []string{`[`, `]`, `"`, `,`} {
		s = strings.Replace(s, char, "", -1)
	}
	list := strings.Fields(s)
	if len(args) < 1 {
		printStrings(list)
	} else {
		if len(args) != 1 {
			cmd.Usage()
		}
		pattern := args[0]
		pattern = strings.ToLower(pattern)
		list = filterStrings(list, pattern)
		printStrings(list)
	}
}
