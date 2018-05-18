package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	sdk "github.com/Noteof/sdk-go"
	"github.com/google/subcommands"
)

type DeleteCmd struct {
	api *sdk.AuthenticatedAPI
	yes bool
}

func (*DeleteCmd) Name() string     { return "delete" }
func (*DeleteCmd) Synopsis() string { return "delete a note." }
func (*DeleteCmd) Usage() string {
	return `list:
	list your notes.
  `
}

func (p *DeleteCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&p.yes, "y", false, "force")
}
func (p *DeleteCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	hasError := false

	if fs.NArg() == 0 {
		fmt.Println("no noteIDs given")
		return subcommands.ExitUsageError
	}

	for _, v := range fs.Args() {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Println(err)
			return subcommands.ExitUsageError
		}

		deleted, err := p.api.DeleteNote(i)
		if deleted && err == nil {
			fmt.Println(i, "DELETED")
		} else if err != nil {
			fmt.Println(i, err)
			hasError = true
		} else {
			fmt.Println(i, "NOT FOUND")
			hasError = true
		}
	}

	if hasError {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
