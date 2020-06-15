package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/k15r/codeowners/pkg/codeowners"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "implement me"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	approvedOwners := arrayFlags{}
	codeownersFile := flag.String("project-dir", ".", "")
	flag.Var(&approvedOwners, "a", ".")
	flag.Parse()

	files := []string{}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	} else {
		fmt.Println("stdin is from a terminal")
	}

	if len(files) == 0 {
		fmt.Print("No input")
		os.Exit(0)
	}

	owners, err := codeowners.NewCodeowners(*codeownersFile)
	if err != nil {
		fmt.Print(err)
	}

	groups := owners.Groups(files, approvedOwners)

	for _, group := range groups {
		fmt.Println(group)
	}

}

type solution struct {
	owners        []string
	remaining_req map[string]map[string]bool
}

//TODO find minimal solution
func filter(sol solution) []solution {
	solutions := []solution{}
	if len(sol.remaining_req) == 0 {
		return append(solutions, sol)
	}
	keys := make([]string, 0, len(sol.remaining_req))
	for k := range sol.remaining_req {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	file := keys[0]
	for owner := range sol.remaining_req[file] {
		newsol := sol
		newsol.owners = append(newsol.owners, owner)
		newsol.remaining_req = filterFilesByOwner(owner, sol.remaining_req)
		addtionalSols := filter(newsol)
		if len(addtionalSols) > 0 {
			solutions = append(solutions, addtionalSols...)
		}
	}

	return solutions
}

func filterFilesByOwner(owner string, req map[string]map[string]bool) map[string]map[string]bool {
	newmap := map[string]map[string]bool{}
	for file, owners := range req {
		if !owners[owner] {
			newmap[file] = owners
		}
	}
	return newmap
}
