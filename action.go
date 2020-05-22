package carapace

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/rsteube/carapace/bash"
	"github.com/rsteube/carapace/elvish"
	"github.com/rsteube/carapace/fish"
	"github.com/rsteube/carapace/powershell"
	"github.com/rsteube/carapace/zsh"
	"github.com/spf13/cobra"
)

type Action struct {
	Bash       string
	Elvish     string
	Fish       string
	Zsh        string
	Powershell string
	Callback   CompletionCallback
}
type ActionMap map[string]Action
type CompletionCallback func(args []string) Action

// finalize replaces value if a callback function is set
func (a Action) finalize(cmd *cobra.Command, uid string) Action {
	if a.Callback != nil {
		if a.Bash == "" {
			a.Bash = bash.Callback(cmd.Root().Name(), uid)
		}
		if a.Elvish == "" {
			a.Elvish = elvish.Callback(cmd.Root().Name(), uid)
		}
		if a.Fish == "" {
			a.Fish = fish.Callback(cmd.Root().Name(), uid)
		}
		if a.Powershell == "" {
			a.Powershell = powershell.Callback(cmd.Root().Name(), uid)
		}
		if a.Zsh == "" {
			a.Zsh = zsh.Callback(uid)
		}
	}
	return a
}

func (a Action) Value(shell string) string {
	switch shell {
	case "bash":
		return a.Bash
	case "fish":
		return a.Fish
	case "elvish":
		return a.Elvish
	case "powershell":
		return a.Powershell
	case "zsh":
		return a.Zsh
	default:
		return ""
	}
}

func (a Action) NestedValue(args []string, shell string, maxDepth int) string {
	if value := a.Value(shell); value == "" && a.Callback != nil && maxDepth > 0 {
		return a.Callback(args).NestedValue(args, shell, maxDepth-1)
	} else {
		return value
	}
}

func (m *ActionMap) Shell(shell string) map[string]string {
	actions := make(map[string]string, len(completions.actions))
	for key, value := range completions.actions {
		actions[key] = value.Value(shell)
	}
	return actions
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{Callback: callback}
}

// ActionExecute uses command substitution to invoke a command and evalues it's result as Action
func ActionExecute(command string) Action {
	return Action{
		Bash:       bash.ActionExecute(command),
		Elvish:     elvish.ActionExecute(command),
		Fish:       fish.ActionExecute(command),
		Powershell: powershell.ActionExecute(command),
		Zsh:        zsh.ActionExecute(command),
	}
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

func ActionDirectories() Action {
	return Action{
		Bash:       bash.ActionDirectories(),
		Elvish:     elvish.ActionDirectories(),
		Fish:       fish.ActionDirectories(),
		Powershell: powershell.ActionDirectories(),
		Zsh:        zsh.ActionDirectories(),
	}
}

func ActionFiles(suffix string) Action {
	return Action{
		Bash:       bash.ActionFiles(suffix),
		Elvish:     elvish.ActionFiles(suffix),
		Fish:       fish.ActionFiles(suffix),
		Powershell: powershell.ActionFiles(suffix),
		Zsh:        zsh.ActionFiles("*" + suffix),
	}
}

// ActionNetInterfaces completes network interface names
func ActionNetInterfaces() Action {
	return Action{
		Bash:       bash.ActionNetInterfaces(),
		Elvish:     elvish.ActionNetInterfaces(),
		Fish:       fish.ActionNetInterfaces(),
		Powershell: powershell.ActionNetInterfaces(),
		Zsh:        zsh.ActionNetInterfaces(),
	}
}

// ActionUsers completes user names
func ActionUsers() Action {
	return Action{
		Bash: bash.ActionUsers(),
		Fish: fish.ActionUsers(),
		Zsh:  zsh.ActionUsers(),
		Callback: func(args []string) Action {
			users := []string{}
			if content, err := ioutil.ReadFile("/etc/passwd"); err == nil {
				for _, entry := range strings.Split(string(content), "\n") {
					user := strings.Split(entry, ":")[0]
					if len(strings.TrimSpace(user)) > 0 {
						users = append(users, user)
					}
				}
			}
			return ActionValues(users...)
		},
	}
}

// ActionGroups completes group names
func ActionGroups() Action {
	return Action{
		Bash: bash.ActionGroups(),
		Fish: fish.ActionGroups(),
		Zsh:  zsh.ActionGroups(),
		Callback: func(args []string) Action {
			// TODO windows
			groups := []string{}
			if content, err := ioutil.ReadFile("/etc/group"); err == nil {
				for _, entry := range strings.Split(string(content), "\n") {
					group := strings.Split(entry, ":")[0]
					if len(strings.TrimSpace(group)) > 0 {
						groups = append(groups, group)
					}
				}
			}
			return ActionValues(groups...)
		},
	}
}

// ActionHosts completes host names
func ActionHosts() Action {
	return Action{
		Bash: bash.ActionHosts(),
		Fish: fish.ActionHosts(),
		Zsh:  zsh.ActionHosts(),
		Callback: func(args []string) Action {
			hosts := []string{}
			if file, err := homedir.Expand("~/.ssh/known_hosts"); err == nil {
				if content, err := ioutil.ReadFile(file); err == nil {
					r := regexp.MustCompile(`^(?P<host>[^ ,#]+)`)
					for _, entry := range strings.Split(string(content), "\n") {
						if r.MatchString(entry) {
							hosts = append(hosts, r.FindStringSubmatch(entry)[0])
						}
					}
				} else {
					return ActionValues(err.Error())
				}
			}
			return ActionValues(hosts...)
		},
	}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	return Action{
		Bash:       bash.ActionValues(values...),
		Elvish:     elvish.ActionValues(values...),
		Fish:       fish.ActionValues(values...),
		Powershell: powershell.ActionValues(values...),
		Zsh:        zsh.ActionValues(values...),
	}
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	return Action{
		Bash:       bash.ActionValuesDescribed(values...),
		Elvish:     elvish.ActionValuesDescribed(values...),
		Fish:       fish.ActionValuesDescribed(values...),
		Powershell: powershell.ActionValuesDescribed(values...),
		Zsh:        zsh.ActionValuesDescribed(values...),
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return Action{
		Bash:       bash.ActionMessage(msg),
		Elvish:     elvish.ActionMessage(msg),
		Fish:       fish.ActionMessage(msg),
		Powershell: powershell.ActionMessage(msg),
		Zsh:        zsh.ActionMessage(msg),
	}
}

func ActionPrefixValues(prefix string, values ...string) Action {
	return Action(Action{
		Bash:       bash.ActionPrefixValues(prefix, values...),
		Elvish:     elvish.ActionPrefixValues(prefix, values...),
		Fish:       fish.ActionPrefixValues(prefix, values...),
		Powershell: powershell.ActionPrefixValues(prefix, values...),
		Zsh:        zsh.ActionPrefixValues(prefix, values...),
	})
}

// TODO find a better solution for this
var CallbackValue string

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char
func ActionMultiParts(divider string, callback func(args []string, parts []string) []string) Action {
	return ActionCallback(func(args []string) Action {
		// TODO multiple dividers by splitting on each char
		index := strings.LastIndex(CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = CallbackValue
		} else if index != -1 {
			prefix = CallbackValue[0 : index+1]
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 {
			parts = parts[0 : len(parts)-1]
		}

		return ActionPrefixValues(prefix, callback(args, parts)...)
	})
}
