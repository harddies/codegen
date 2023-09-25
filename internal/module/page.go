package module

import (
	"codegen/internal/arg"
	"codegen/internal/model"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type page struct {
	arg.Sets
}

func (d *page) Name() string {
	return model.FlagModulePage
}

func (d *page) Run(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	t, err := template.New("page").Parse(pageTpl)
	if err != nil {
		return
	}

	f, err := os.Create(d.Output)
	if err != nil {
		return
	}
	defer f.Close()

	if err = t.Execute(f, d.Sets); err != nil {
		return
	}
}

const (
	pageTpl = `package page

import (
	"io"

	"gms-back/err"
	"gms-back/project"
)

{{ range .FuncNames }}
func {{ . }}(p project.IProject, target string, iReq, iMeta interface{}) (res interface{}, ei err.ErrInfo) {
TAG:
	res = body
	return
}

{{ end }}
func init() {
	{{ range .FuncNames }}
	project.APIRegist(api.{{ . }}, []string{constant.Mist, constant.Dream}, project.Info{
		Auth: false,
		Unmarshal: func(body io.ReadCloser) (interface{}, error) {
			return project.APIUnmarshal(body, &model.{{ . }}Req{})
		},
		Execute: {{ . }},
	})
	{{ end }}
}
`
)