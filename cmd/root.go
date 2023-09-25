package cmd

import (
	"codegen/internal/arg"
	"codegen/internal/model"
	"codegen/internal/module"
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

var once = sync.Once{}

// rootCmd represents the root command
var (
	rootCmd *cobra.Command

	moduleValue    string
	funcNamesValue []string
)

func init() {
	once.Do(func() {
		rootCmd = &cobra.Command{
			Use:   "codegen",
			Short: "A code generation",
			Long:  `It can generate code of dao, grpc and so on`,
			Run: func(cmd *cobra.Command, args []string) {
				argSets := arg.New(arg.SetModule(moduleValue), arg.SetFuncNames(funcNamesValue), arg.SetOutput(""))

				c := module.NewStrategy(cmd, argSets)
				c.Run(cmd, args)
				fmt.Println(moduleValue, funcNamesValue)
			},
		}
	})
	rootCmd.Flags().StringVarP(&moduleValue, model.FlagNameModule, "m", "", "generate which code of module. [required]\n - dao\n - page")
	rootCmd.Flags().StringSliceVarP(&funcNamesValue, model.FlagNameFuncNames, "f", nil, "specify function names that need to generate.")
}

func Execute() error {
	return rootCmd.Execute()
}
