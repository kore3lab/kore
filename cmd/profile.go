package cmd

import (
	"github.com/spf13/cobra"

	"kore3lab.io/kore/manifests"
)

// a struct to support command
type ProfileOptions struct {
	*Options
	ProfileName string
}

// returns a cobra command
func NewCommandProfile(options *Options) *cobra.Command {
	o := &ProfileOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:   "profile",
		Short: "The profile subcommand",
		Long:  "",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// list
	cmdList := &cobra.Command{
		Use:                   "list",
		Short:                 "Lists available Kore configuration profiles",
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			ValidateError(c, func() error {
				if files, err := manifests.ListProfiles(); err != nil {
					o.Println(err.Error())
				} else {
					for _, file := range files {
						o.Println(file)
					}
				}
				return nil
			}())
		},
	}
	cmds.AddCommand(cmdList)

	// get
	cmdGet := &cobra.Command{
		Use:                   "dump",
		Short:                 "Get a Kore configuration profile",
		DisableFlagsInUseLine: true,
		Args:                  bindCommandArgs(&o.ProfileName),
		Run: func(c *cobra.Command, args []string) {
			ValidateError(c, func() error {
				if o.ProfileName == "" {
					o.ProfileName = "default"
				}
				if by, err := manifests.ReadFile(o.ProfileName); err != nil {
					o.Println(err.Error())
				} else {
					o.Println(string(by))
				}
				return nil
			}())
		},
	}
	cmdGet.Flags().StringVarP(&o.ProfileName, "name", "", "", "Name of profile")

	cmds.AddCommand(cmdGet)

	return cmds
}
