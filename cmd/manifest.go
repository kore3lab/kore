package cmd

import (
	"github.com/spf13/cobra"

	"kore3lab.io/kore/manifests"
	"kore3lab.io/kore/pkg/helm"
)

// a struct to support command
type ManifestOptions struct {
	*Options
}

// returns a cobra command
func NewCommandManifest(options *Options) *cobra.Command {
	o := &ManifestOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:                   "manifest",
		Short:                 "The management of manifests",
		Long:                  "",
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
			//ValidateError(c, func() error {
			//	return nil
			//}())
		},
	}

	// generate
	cmdGenerate := &cobra.Command{
		Use:                   "generate",
		Short:                 "Print manifesats",
		DisableFlagsInUseLine: true,
		Args:                  bindCommandArgs(&o.ProfileName),
		Run: func(c *cobra.Command, args []string) {
			ValidateError(c, func() error {
				if operator, err := manifests.GetProfile(o.ProfileName); err != nil {
					return err
				} else {
					if manifeset, err := helm.RenderToYaml(operator, o.Values); err != nil {
						return err
					} else {
						o.Println(manifeset)
					}
				}
				return nil
			}())
		},
	}
	cmdGenerate.Flags().StringVarP(&o.ProfileName, "name", "", "", "Name of profile")

	cmds.AddCommand(cmdGenerate)

	return cmds
}
