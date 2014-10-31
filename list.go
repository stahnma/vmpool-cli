package main

import "net/http"
import "io/ioutil"
import "strings"
import "os"
import "log"

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
	resp, err := http.Get(vmpool_url)
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
	pattern := args[0]
	if pattern == "" {
		printStrings(list)
	} else {
		pattern = strings.ToLower(pattern)
		list = filterStrings(list, pattern)
		printStrings(list)
	}
}
