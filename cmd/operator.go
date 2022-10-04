package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"kore3lab.io/kore/manifests"
	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
	"kore3lab.io/kore/pkg/helmreconciler"
)

// a struct to support command
type OperatorOptions struct {
	*Options
}

// returns a cobra command
func NewCommandOperator(options *Options) *cobra.Command {

	o := &OperatorOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:                   "operator",
		Short:                 "The management of kore controller",
		Long:                  "",
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// init
	cmds.AddCommand(&cobra.Command{
		Use:                   "init",
		Short:                 "remove controller",
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			ValidateError(c, func() error {
				if operator, err := manifests.GetProfile("base"); err != nil {
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
		},
	})

	// remove
	cmds.AddCommand(&cobra.Command{
		Use:                   "remove",
		Short:                 "Remove controller",
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			ValidateError(c, func() error {
				if operator, err := manifests.GetProfile("base"); err != nil {
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
		},
	})

	return cmds
}
