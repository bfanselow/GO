
/***************************************************************************
 git_api.go
 bfanselow 2017-10-31

 This app will pull a list of projects from githup.com API for a particular
 user.

 GitHup API: "https://api.github.com/users/bfanselow/repos"

 Usage: go run git_api.go -u <username> [options]

 Options:
  -d <debug-level>  Default=0

 Example:
 > go run git_api.go -u bfanselow

 Output:
Found 2 Project(s) for user (bfanselow)
 (1) Project-Id: 108909372
 (1) Project-name: GO-tests
 (1) Full-name: bfanselow/GO-tests
==================
 (2) Project-Id: 108910653
 (2) Project-name: Python
 (2) Full-name: bfanselow/Python
==================

***************************************************************************/
package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag" // alias pflag to "flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// constants
const (
	apiURLbase   = "https://api.github.com"
	userEndpoint = "/users/"
)

// CLI data
var (
	user string
)

// behavior defaults
var debug = 0

//=================================================================================

// GitProject struct represents the JSON user data from GitHub API
type GitProject struct {
	Project_id   int    `json:"id"`
	Project_name string `json:"name"`
	Full_name    string `json:"full_name"`
}

// API response is a Slice of these structs
// We create to store the "Unmarshal"-ed (aka parsed JSON) data
var GitProjectList = []GitProject{}

//=================================================================================
func init() {
	// We pass the user variable we declared at the package level (above).
	// The "&" character means we are passing the variable "by reference" (as opposed to "by value"),
	// meaning: we don't want to pass a copy of the user variable. We want to pass the original variable.
	flag.StringVarP(&user, "user", "u", "", "Username for whith to search API")
	flag.IntVarP(&debug, "debug", "d", 0, "Change DEBUG level")
}

//=================================================================================
// printUsage is a custom function we created to print app usage
func printUsage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}

//=================================================================================
// Main()
func main() {
	flag.Parse() // parse CLI flags

	// If no input, print usage
	if flag.NFlag() == 0 {
		printUsage()
	}

	// If multiple users are passed separated by commas, store them in a "users" array
	users := strings.Split(user, ",")

	// If no user specified, print usage
	if len(users) == 0 {
		printUsage()
	}

	for _, uname := range users {
		if debug > 1 {
			fmt.Printf("Searching API for user: %s\n", uname)
		}

		// getGitProjectsForUser() returns slice of bytes
		// API returns an array of json objects, each containg info about a project.
		// ioutil.ReadAll(resp.Body) converts this to slice of bytes
		bs_user := getGitProjectsForUser(uname)

		if debug > 1 {
			// quick-n-dirty transformation of bytes to string will show us the raw JSON object
			json_user := string(bs_user)
			fmt.Printf(" %s\n", json_user)
		}

		// create a "resp" variable of type "GitProjectList" to store the "Unmarshal"-ed json data
		resp := GitProjectList
		json.Unmarshal(bs_user, &resp)

		var N_projects = len(resp)
		fmt.Printf("Found %d Project(s) for user (%s)\n", N_projects, uname)

		for i, GitProject := range resp {
			var n = i + 1
			fmt.Printf(" (%d) Project-Id: %d\n", n, GitProject.Project_id)
			fmt.Printf(" (%d) Project-name: %s\n", n, GitProject.Project_name)
			fmt.Printf(" (%d) Full-name: %s\n", n, GitProject.Full_name)
			fmt.Println("==================")
		}
	}
}

//=================================================================================
// getGitProjectsForUser() queries API for a given user. Return byte-array.
func getGitProjectsForUser(name string) []byte {

	var apiURL = apiURLbase + userEndpoint + name + "/repos"

	if debug > 0 {
		fmt.Printf("Getting user (%s) from API: [%s]\n", name, apiURL)
	}

	// send GET request to API with the requested user
	resp, err := http.Get(apiURL)

	// Always good practice to defer closing the response body.
	defer resp.Body.Close()

	// if err occurs during GET request, then throw error and quit application
	if err != nil {
		log.Fatalf("Error retrieving data: %s\n", err)
	}

	// read the response body and handle any errors during reading. body is byte-slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading data: %s\n", err)
	}
	return body
}
