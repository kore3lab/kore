package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := NewRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type KoreCtlOptions struct {
	Options
}

func NewRootCommand() *cobra.Command {

	o := KoreCtlOptions{}

	// ctl
	cmds := &cobra.Command{
		Use:                   "korectl",
		Short:                 "Kore3. command-line-interface manager",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Persistent Flags
	cmds.PersistentFlags().StringVarP(&o.ProfileName, "profile", "p", "", "Profile")
	cmds.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "kore-system", "Name of namespace")
	cmds.PersistentFlags().StringVarP(&o.Filename, "file", "f", "", "Filename")
	cmds.PersistentFlags().StringArrayVar(&o.Values, "set", []string{}, `Override an setting value 
(can specify multiple or separate values with commas: key1=val1,key2=val2), 
e.g. enable or disable components (--set components.dashboard.enabled=true), 
or override values settings (--set metrics-server.enabled=true).`)

	// add commands
	cmds.AddCommand(&cobra.Command{
		Use:                   "version",
		Short:                 "Print the version number of korectl",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Println("Version=%s, buildTime=%s", BuildVersion, BuildTime)
		},
	})

	cmds.AddCommand(NewCommandProfile(&o.Options))   // profile command
	cmds.AddCommand(NewCommandInstall(&o.Options))   // install command
	cmds.AddCommand(NewCommandUninstall(&o.Options)) // uninstall command
	cmds.AddCommand(NewCommandManifest(&o.Options))  // manifest command
	cmds.AddCommand(NewCommandOperator(&o.Options))  // operator command

	return cmds

}

func ValidateError(c *cobra.Command, err error) {

	if err != nil {
		msg := err.Error()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

}
func bindCommandArgs(values ...*string) func(c *cobra.Command, args []string) error {

	return func(c *cobra.Command, args []string) error {

		for i, v := range args {
			if len(values) > i {
				*values[i] = v
			}
		}
		return nil
	}

}
