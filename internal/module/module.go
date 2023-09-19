package module

import (
	"codegen/internal/model"

	"github.com/spf13/cobra"
)

type IModule interface {
	Run(c *cobra.Command, args []string)
}

type Strategy struct {
	m IModule
}

func (s *Strategy) Run(command *cobra.Command, args []string) {
	s.m.Run(command, args)
}

func NewStrategy(command *cobra.Command) (s *Strategy) {
	s = new(Strategy)
	switch command.Flag(model.FlagNameModule).Value.String() {
	case model.FlagModuleDao:
		s.m = new(dao)
	}
	return
}
