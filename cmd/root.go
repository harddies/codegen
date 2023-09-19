package cmd

import (
	"sync"

	"codegen/internal/model"
	"codegen/internal/module"

	"github.com/spf13/cobra"
)

var once = sync.Once{}

// rootCmd represents the root command
var rootCmd *cobra.Command

func init() {
	once.Do(func() {
		rootCmd = &cobra.Command{
			Use:   "codegen",
			Short: "A code generation",
			Long:  `It can generate code of dao, grpc and so on`,
			Run: func(cmd *cobra.Command, args []string) {
				c := module.NewStrategy(cmd)
				c.Run(cmd, args)
			},
		}
	})
	rootCmd.Flags().String(model.FlagNameModule, "", "generate which code of module.\n - dao")
}

func Execute() error {
	return rootCmd.Execute()
}
