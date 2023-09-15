package module

import (
	"codegen/internal/model"

	"github.com/spf13/cobra"
)

type IModule interface {
	Run(cmd *cobra.Command, args []string)
}

func New(cmd *cobra.Command, args []string) (m IModule) {
	switch cmd.Flag(model.FlagNameModule).Value.String() {
	case model.FlagModuleDao:
		m = new(dao)
	}
	return
}
