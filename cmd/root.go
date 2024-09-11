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
	target string
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
					arg.SetTarget(target),
				)

				c := module.NewStrategy(cmd, argSets)
				c.Run(cmd, args)
			},
		}
	})
	rootCmd.Flags().StringVarP(&mod, model.FlagNameModule, model.FlagNameModuleShort, "", "generate which code of module. [required]\n - dao\n - bts: generate bts code that is rw between cache and db base on kratosV2 data. e.g.\n\t//go:generate codegen -m bts -t ../../data/\n\t// bts: -null_cache=&do.Demo{} -empty_value=nil -check_null_code=$!=nil&&$.ID==0 -struct_name=thirdPartnerRepo -single_flight_var=g")
	rootCmd.Flags().StringSliceVarP(&fns, model.FlagNameFuncNames, model.FlagNameFuncNamesShort, nil, "specify function names that need to generate.")
	rootCmd.Flags().StringVarP(&output, model.FlagNameOutput, model.FlagNameOutputShort, "", "specify function names that need to generate.")
	rootCmd.Flags().StringVarP(&target, model.FlagNameTarget, model.FlagNameTargetShort, "", "specify target dir. use to bts")
}

func Execute() error {
	return rootCmd.Execute()
}
