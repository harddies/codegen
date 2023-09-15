package module

import (
	"fmt"

	"github.com/spf13/cobra"
)

type dao struct {
}

func (d *dao) Run(cmd *cobra.Command, args []string) {
	fmt.Println(cmd.Flag("module").Value)
}
