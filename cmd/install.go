package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"kore3lab.io/kore/manifests"
	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
	"kore3lab.io/kore/pkg/helmreconciler"
)

// a struct to support command
type InstallOptions struct {
	*Options
}

// returns a cobra command
func NewCommandInstall(options *Options) *cobra.Command {
	o := &InstallOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:                   "install (NAME) [flags]",
		Short:                 "The installation subcommand",
		Long:                  "",
		DisableFlagsInUseLine: true,
		Args:                  bindCommandArgs(&o.ProfileName),
		Run: func(c *cobra.Command, args []string) {
			if o.ProfileName == "" {
				c.Help()
			} else {
				ValidateError(c, func() error {
					if operator, err := manifests.GetProfile(o.ProfileName); err != nil {
						return err
					} else {
						if reconciler, err := helmreconciler.NewHelmReconciler(operator, o.Values, nil); err != nil {
							return err
						} else {
							if status := reconciler.Reconcile(); status.Status == installv1alpha1.STATUS_ERROR {
								return fmt.Errorf(status.Message)
							}
						}
					}
					return nil
				}())
			}
		},
	}
	cmds.Flags().StringVarP(&o.ProfileName, "name", "", "", "Name of profile")

	return cmds
}
