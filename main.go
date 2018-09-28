package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/errskipower/find-issues/issues"
	flags "github.com/jessevdk/go-flags"
)

type opts struct {
	Creator string `long:"creator" description:"Filter issues based on their creator."`
	Label   string `long:"label" description:"Filter issues based on their labels."`
}

func main() {
	var options opts
	parser := flags.NewParser(&options, flags.HelpFlag|flags.PrintErrors)
	remainingArgs, err := parser.ParseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var repo string
	if len(remainingArgs) == 0 {
		log.Fatal("Missing repo. Pass in something like `ghc-tdd/find-issues`!")
	}

	repo = remainingArgs[1]

	httpClient := &http.Client{}

	service := issues.NewService(repo, httpClient)

	issues, err := service.Get(options.Label, options.Creator)
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		fmt.Printf("#%d: %s\n", issue.Number, issue.Title)
	}
}
