package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/example/cmd/action"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:     "action",
	Short:   "action example",
	Aliases: []string{"alias"},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	actionCmd.Flags().StringP("files", "f", "", "files flag")
	actionCmd.Flags().String("directories", "", "files flag")
	actionCmd.Flags().StringP("groups", "g", "", "groups flag")
	actionCmd.Flags().String("hosts", "", "hosts flag")
	actionCmd.Flags().StringP("message", "m", "", "message flag")
	actionCmd.Flags().StringP("net_interfaces", "n", "", "net_interfaces flag")
	actionCmd.Flags().String("usergroup", "", "user:group flag")
	actionCmd.Flags().StringP("users", "u", "", "users flag")
	actionCmd.Flags().StringP("values", "v", "", "values flag")
	actionCmd.Flags().StringP("values_described", "d", "", "values with description flag")
	//actionCmd.Flags().StringS("shorthandonly", "s", "", "shorthandonly flag")
	actionCmd.Flags().StringP("kill", "k", "", "kill signals")
	actionCmd.Flags().StringP("optarg", "o", "", "optional arg with default value blue")
	actionCmd.Flag("optarg").NoOptDefVal = "blue"

	carapace.Gen(actionCmd).FlagCompletion(carapace.ActionMap{
		"files":            carapace.ActionFiles(".go"),
		"directories":      carapace.ActionDirectories(),
		"groups":           action.ActionGroups(),
		"hosts":            action.ActionHosts(),
		"message":          carapace.ActionMessage("message example"),
		"net_interfaces":   action.ActionNetInterfaces(),
		"usergroup":        action.ActionUserGroup(),
		"users":            action.ActionUsers(),
		"values":           carapace.ActionValues("values", "example"),
		"values_described": carapace.ActionValuesDescribed("values", "valueDescription", "example", "exampleDescription"),
		"kill":             action.ActionKillSignals(),
		"optarg":           carapace.ActionValues("blue", "red", "green", "yellow"),
	})

	carapace.Gen(actionCmd).PositionalCompletion(
		carapace.ActionValues("positional1", "p1"),
		carapace.ActionValues("positional2", "p2"),
	)
}
