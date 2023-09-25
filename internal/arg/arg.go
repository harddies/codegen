package arg

import (
	"fmt"
	"os"
	"path"

	"codegen/internal/model"
)

var (
	dir string
)

type Options func(s *Sets)

func New(opts ...Options) *Sets {
	var err error

	if dir, err = os.Getwd(); err != nil {
		fmt.Println(err)
	}

	s := &Sets{
		FuncNames: []string{model.FuncNamesDefaultName},
		Output:    path.Join(dir, model.OutputFilenameDefaultName),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Sets struct {
	Module    string
	FuncNames []string
	Output    string
}
