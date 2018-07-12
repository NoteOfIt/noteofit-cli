package main

import (
	"context"
	"flag"
	"log"
	"os"

	noteofcli "github.com/Noteof/noteof-cli"
	sdk "github.com/Noteof/sdk-go"
	"github.com/google/subcommands"
)

func main() {
	config, err := noteofcli.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	editor := config.GetEditor()
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	uapi := &sdk.UnauthenticatedAPI{}
	aapi := sdk.NewAuthenticatedApi(config.GetToken())

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&InitCmd{config, uapi}, "")

	list := &ListCmd{aapi}
	subcommands.Register(list, "")
	subcommands.Register(subcommands.Alias("ls", list), "")
	subcommands.Register(&GetCmd{aapi}, "")

	delete := &DeleteCmd{api: aapi}
	subcommands.Register(delete, "")
	subcommands.Register(subcommands.Alias("rm", delete), "")
	subcommands.Register(&NewCmd{editor, config, aapi}, "")
	subcommands.Register(&EditCmd{editor, aapi}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
