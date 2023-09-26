package cmd

import (
	"sync"

	"github.com/harddies/codegen/internal/arg"
	"github.com/harddies/codegen/internal/model"
	"github.com/harddies/codegen/internal/module"

	"github.com/spf13/cobra"
)

var once = sync.Once{}

// rootCmd represents the root command
var (
	rootCmd *cobra.Command

	mod    string
	fns    []string
	output string
)

func init() {
	once.Do(func() {
		rootCmd = &cobra.Command{
			Use:   "codegen",
			Short: "A code generation",
			Long:  `It can generate code of page, dao, grpc and so on`,
			Run: func(cmd *cobra.Command, args []string) {
				argSets := arg.New(
					arg.SetModule(mod),
					arg.SetFuncNames(fns),
					arg.SetOutput(output),
				)

				c := module.NewStrategy(cmd, argSets)
				c.Run(cmd, args)
			},
		}
	})
	rootCmd.Flags().StringVarP(&mod, model.FlagNameModule, model.FlagNameModuleShort, "", "generate which code of module. [required]\n - dao\n - page")
	rootCmd.Flags().StringSliceVarP(&fns, model.FlagNameFuncNames, model.FlagNameFuncNamesShort, nil, "specify function names that need to generate.")
	rootCmd.Flags().StringVarP(&output, model.FlagNameOutput, model.FlagNameOutputShort, "", "specify function names that need to generate.")
}

func Execute() error {
	return rootCmd.Execute()
}
