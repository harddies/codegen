package module

import (
	"fmt"

	"codegen/internal/arg"
	"codegen/internal/model"

	"github.com/spf13/cobra"
)

type dao struct {
	arg.Sets
}

func (d *dao) Name() string {
	return model.FlagModuleDao
}

func (d *dao) Run(cmd *cobra.Command, args []string) {
	fmt.Println(d.Name())
}

const (
	daoTpl = `package dao

import (

)

func (d *{ .data.Name }) Get
`
)
