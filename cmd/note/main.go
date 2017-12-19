package main

import (
	"context"
	"flag"
	"log"
	"os"

	noteofitcli "github.com/NoteOfIt/noteofit-cli"
	sdk "github.com/NoteOfIt/sdk-go"
	"github.com/google/subcommands"
)

func main() {
	config, err := noteofitcli.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	uapi := &sdk.UnauthenticatedAPI{}
	aapi := sdk.NewAuthenticatedApi(config.GetToken())

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&InitCmd{config, uapi}, "")
	subcommands.Register(&ListCmd{aapi}, "")
	subcommands.Register(&GetCmd{"vim", aapi}, "")
	subcommands.Register(&DeleteCmd{api: aapi}, "")
	subcommands.Register(&NewCmd{"vim", config, aapi}, "")
	subcommands.Register(&EditCmd{"vim", aapi}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))

	// aapi.
}
