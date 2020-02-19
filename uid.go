package zsh

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func uidCommand(cmd *cobra.Command) string {
	names := make([]string, 0)
	current := cmd
	for {
		names = append(names, current.Name())
		current = current.Parent()
		if current == nil {
			break
		}
	}

	reverse := make([]string, len(names))
	for i, entry := range names {
		reverse[len(names)-i-1] = entry
	}

	return "_" + strings.Join(reverse, "__")
}

func uidFlag(cmd *cobra.Command, flag *pflag.Flag) string {
	// TODO ensure flag acually belongs to command (force error)
	// TODO handel unknown flag nil error
	return fmt.Sprintf("%v##%v", uidCommand(cmd), flag.Name)
}

func uidPositional(cmd *cobra.Command, position int) string {
	// TODO complete function
	return fmt.Sprintf("%v#%v", uidCommand(cmd), position)
}

func parse(uid string) []string {
	var splitted []string
	if splitted = strings.Split(uid[1:], "#"); len(splitted) == 0 { // TODO check for empty uid string
		return nil
	}
	return strings.Split(splitted[0], "__")
}

func find(cmd *cobra.Command, uid string) *cobra.Command {
	var splitted []string
	if splitted = strings.Split(uid[1:], "#"); len(splitted) == 0 { // TODO check for empty uid string
		return nil
	}
	c, _, err := cmd.Root().Find(strings.Split(splitted[0], "__")[1:]) // TODO root if jut one arg
	if err != nil {
		log.Fatal(err)
	}
	return c
}
