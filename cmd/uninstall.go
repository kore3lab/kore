package cmd

import (
	"github.com/spf13/cobra"

	"kore3lab.io/kore/manifests"
	"kore3lab.io/kore/pkg/helmreconciler"
)

// a struct to support command
type UninstallOptions struct {
	*Options
	ComponentName string
}

// returns a cobra command
func NewCommandUninstall(options *Options) *cobra.Command {
	o := &UninstallOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:                   "uninstall (NAME) [flags]",
		Short:                 "The uninstallation subcommand",
		Long:                  "",
		DisableFlagsInUseLine: true,
		Args:                  bindCommandArgs(&o.ComponentName),
		Run: func(c *cobra.Command, args []string) {
			if o.ComponentName == "" {
				c.Help()
			} else {
				ValidateError(c, func() error {
					if operator, err := manifests.GetProfile(o.ProfileName); err != nil {
						return err
					} else {
						if reconciler, err := helmreconciler.NewHelmReconciler(operator, o.Values, nil); err != nil {
							return err
						} else {
							if err := reconciler.Finalize(); err != nil {
								return err
							}
						}
					}
					return nil
				}())
			}
		},
	}

	return cmds
}
