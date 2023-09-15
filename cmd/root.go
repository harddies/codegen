package cmd

import (
	"codegen/internal/module"

	"github.com/spf13/cobra"
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "codegen",
	Short: "A code generation",
	Long:  `It can generate code of dao, grpc and so on`,
	Run: func(cmd *cobra.Command, args []string) {
		m := module.New(cmd, args)
		m.Run(cmd, args)
	},
}

func init() {
	rootCmd.Flags().String("module", "", "generate which code of module.\n - dao")
}

func Execute() {
	_ = rootCmd.Execute()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
