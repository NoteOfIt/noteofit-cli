package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	noteofitcli "github.com/NoteOfIt/noteofit-cli"
	sdk "github.com/NoteOfIt/sdk-go"
	"github.com/google/subcommands"
	"golang.org/x/crypto/ssh/terminal"
)

type InitCmd struct {
	config *noteofitcli.Config
	api    *sdk.UnauthenticatedAPI
}

func (*InitCmd) Name() string     { return "init" }
func (*InitCmd) Synopsis() string { return "Initialize the noteofit cli." }
func (*InitCmd) Usage() string {
	return `init:
	Initialize the noteofit cli.
  `
}

func (p *InitCmd) SetFlags(f *flag.FlagSet) {}
func (p *InitCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	for {
		username, password, err := credentials()
		if err != nil {
			os.Stdout.WriteString(err.Error() + "\n")
			return subcommands.ExitFailure
		} else if username == "" && password == "" {
			os.Stdout.WriteString("exiting unsuccessfully\n")
			return subcommands.ExitFailure
		}

		tr, err := p.api.DoAuth(username, password, "noteofit-cli")
		if err != nil {
			continue
		}

		token := tr.APIToken
		p.config.SetToken(token)
		fmt.Println("token saved successfully")
		return subcommands.ExitSuccess
	}
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println()

	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}