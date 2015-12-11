package main

import (
	"fmt"
	"log"
	"os"
)

var cmdToken = &Command{
	Run:       runToken,
	UsageLine: "token",
	Short:     "work with tokens from vmpooler",
	Long: `

    token list (default):
      - see all of your tokens.
    token request:
       - get a token for LDAP_USERNAME using LDAP_PASSWORD
    token delete <token>:
      - delete a token
    token purge
      - purge all tokens

    `,
}

func isValidToken(token string) bool {
	debug("Function: isValidToken")
	debug("  Token: " + token)
	//params := ldapsetup()
	params := token
	contents, output_json := RequestWrapper("token", params, "GET", "{}")
	pjson(contents)
	if _, ok := output_json[token]; ok {
		debug("  Returning true from isValidToken")
		return true
	}
	debug("  Returning false from isValidToken")
	return false
}

func listToken() {
	debug("Function listToken")
	params := ldapsetup()
	contents, output_json := RequestWrapper("token", params, "GET", "{}")
	if output_json["ok"] == false {
		fmt.Println("No tokens setup.")
		os.Exit(255)
	}

	pjson(contents)
}

func purgeTokens() {
	debug("Function purgeTokens")
	params := ldapsetup()
	_, output_json := RequestWrapper("token", params, "GET", "{}")
	if output_json["ok"] == false {
		fmt.Println("No tokens setup.")
		os.Exit(255)
	}
	for key, _ := range output_json {
		if key == "ok" {
			continue
		}
		debug("  Key: " + key)
		deleteToken(key)
	}

}

func deleteToken(token string) {
	debug("Function: deleteToken")
	params := ldapsetup()
	params = params + "|" + token
	_, output_json := RequestWrapper("token", params, "DELETE", "{}")

	if output_json["ok"] == true {
		fmt.Println("Deleted token: " + token)
		logmsg("Deleted token: " + token)
	} else {
		fmt.Println("Token " + token + " not found.")
		os.Exit(255)
	}
}

func processEnvForToken() string {
	debug("Function: processEnvForToken")
	// VMPOOL_TOKEN is deprecated, but still should work
	if os.Getenv("VMPOOL_TOKEN") != "" {
		os.Setenv("VMPOOLER_TOKEN", os.Getenv("VMPOOL_TOKEN"))
	}
	if os.Getenv("VMPOOLER_TOKEN") != "" {
		//		if isValidToken(os.Getenv("VMPOOLER_TOKEN")) {
		logmsg("Using VMPOOLER_TOKEN...assuming it is valid: " + os.Getenv("VMPOOLER_TOKEN"))
		debug("VMPOOLER_TOKEN...assuming it is valid: " + os.Getenv("VMPOOLER_TOKEN"))
		debug("  Returning " + os.Getenv("VMPOOLER_TOKEN"))
		return os.Getenv("VMPOOLER_TOKEN")
		//}

	}
	debug("  Returning: \"\"" + " from processEnvForToken")
	return ""
}

func grantToken() string {
	debug("Function grantToken")
	newtoken := processEnvForToken()
	if newtoken == "" {
		debug("Did not find token in ENV, requesting a new one.")
		params := ldapsetup()
		_, output_json := RequestWrapper("token", params, "POST", "{}")
		newtoken = fmt.Sprintf("%v", output_json["token"])
		os.Setenv("VMPOOLER_TOKEN", newtoken)
		logmsg(fmt.Sprintf("Aquired token: %v", newtoken))
		fmt.Println(fmt.Sprintf("Aquired token: %v", newtoken))
	}
	debug("Token is: " + newtoken)
	return newtoken
}

// What happens if I don't have a token?
func retrieveToken() string {
	debug("Function: retrieveToken")
	token := ""
	token = processEnvForToken()
	if token != "" {
		return token
	}

	params := ldapsetup()
	_, output_json := RequestWrapper("token", params, "GET", "{}")
	// Take the first token returned
	for k, _ := range output_json {
		if k == "ok" {
			continue
		}
		token = k
		return token
	}
	if token == "" {
		token = grantToken()
	}
	return token
}

func runToken(cmd *Command, args []string) {
	debug("Function: runToken")
	var subcmd string
	if len(args) > 0 {
		subcmd = args[0]
	} else {
		subcmd = ""
	}

	switch {
	case (subcmd == "request"):
		grantToken()
	case (subcmd == "list" || subcmd == ""):
		listToken()
	case subcmd == "delete":
		if len(args) < 2 {
			log.Fatal("delete requires token argument. ")
		}
		token := args[1]
		deleteToken(token)
	case subcmd == "purge":
		purgeTokens()
	default:
		log.Fatal(args[0] + " is not a recognized action for token.")
	}
}
