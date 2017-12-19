package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"

	noteofitcli "github.com/NoteOfIt/noteofit-cli"
	sdk "github.com/NoteOfIt/sdk-go"
	"github.com/google/subcommands"
)

type PostCmd struct {
	editor string

	config *noteofitcli.Config
	api    *sdk.AuthenticatedAPI
}

func (*PostCmd) Name() string     { return "post" }
func (*PostCmd) Synopsis() string { return "post a new note." }
func (*PostCmd) Usage() string {
	return `post:
	post a new note.
  `
}

var wsre = regexp.MustCompile("\\s")

func (p *PostCmd) SetFlags(f *flag.FlagSet) {}
func (p *PostCmd) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if p.editor != "" {
		body, err := execEditor(p.editor)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

		n := &sdk.Note{
			CurrentText: &sdk.NoteText{
				NoteTextValue: string(body),
			},
		}

		n2, err := p.api.PostNewNote(n)
		if err != nil {
			log.Println(err)
			return subcommands.ExitFailure
		}

		fmt.Println(n.NoteID, getTitleLine(n2.CurrentText.NoteTextValue))
	}

	return subcommands.ExitSuccess
}

func execEditor(editor string) ([]byte, error) {
	parts := wsre.Split(editor, -1)

	tmpfile, err := ioutil.TempFile("", "post")
	tmpfile.WriteString("<new note text>")
	tmpPath := tmpfile.Name()
	tmpfile.Close()
	if err != nil {
		return []byte{}, err
	}

	args := []string{}
	if len(parts) > 1 {
		args = append(args, parts[1:]...)
	}

	args = append(args, tmpPath)

	cmd := exec.Command(parts[0], args...)
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadFile(tmpPath)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}
