package module

import (
	"codegen/internal/arg"
	"codegen/internal/model"
	"os"
	"strings"

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

func NewStrategy(command *cobra.Command, args *arg.Sets) (s *Strategy) {
	s = new(Strategy)
	switch strings.ToLower(args.Module) {
	case model.FlagModuleDao:
		s.m = &dao{Sets: *args}
	case model.FlagModulePage:
		s.m = &page{Sets: *args}
	default:
		_ = command.Usage()
		os.Exit(1)
	}
	return
}
