package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func matchPattern(pattern string) func(name string) bool {
	reg := regexp.MustCompile(pattern)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}

func filterStrings(list []string, pattern string) []string {
	cleaned := []string{}
	match := matchPattern(pattern)
	for _, elm := range list {
		if match(elm) {
			cleaned = append(cleaned, elm)
		}
	}
	return cleaned
}

func printStrings(list []string) {
	for _, elm := range list {
		fmt.Println(elm)
	}
}

func perror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func shortname(hostname string) string {
	if matchPattern(".")(hostname) {
		hostname = strings.Split(hostname, ".")[0]
	}
	return hostname
}

func logmsg(message string) {
	logfile := os.Getenv("HOME") + "/.vmpool.log"
	if os.Getenv("VMPOOL_LOGFILE") != "" {
		logfile = os.Getenv("VMPOOL_LOGFILE")
	}
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	message = strings.TrimSpace(message)
	log.Println(message)
}

func ldapsetup() string {
	errmsg := ""
	bombout := 0
	if os.Getenv("LDAP_USERNAME") == "" {
		errmsg = errmsg + "The environment variable LDAP_USERNAME must be set to request a token.\n"
		bombout = 1
	}
	if os.Getenv("LDAP_PASSWORD") == "" {
		errmsg = errmsg + "The environment variable LDAP_PASSWORD must be set to request a token.\n"
		bombout = 1
	}
	if bombout == 1 {
		log.Fatal(errmsg)
	}
	ldap_user := os.Getenv("LDAP_USERNAME")
	ldap_pass := os.Getenv("LDAP_PASSWORD")
	debug("LDAP_USERNAME " + ldap_user)
	debug("LDAP_PASSWORD " + ldap_pass)
	return ldap_user + "|" + ldap_pass
}

func pjson(contents []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, contents, "", "  ")
	perror(error)
	fmt.Println(string(prettyJSON.Bytes()))
}

func debug(message string) {
	if os.Getenv("DEBUG") == "1" || os.Getenv("DEBUG") == "true" {
		fmt.Printf("%v\n", message)
	}
}
