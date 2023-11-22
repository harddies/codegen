package module

import (
	"fmt"
	"os"
	"text/template"

	"github.com/harddies/codegen/internal/arg"
	"github.com/harddies/codegen/internal/model"

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

	t, err := template.New(pageTplName).Parse(pageTpl)
	if err != nil {
		return
	}

	f, err := os.OpenFile(d.Output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	if err = t.Execute(f, d.Sets); err != nil {
		return
	}

	fmt.Printf("generate file [%s] successfully\n", d.Output)
}

const (
	pageTplName = "page"
	pageTpl     = `package page

import (
	"context"
	"io"

	"gms-back/api"
	"gms-back/constant"
	"gms-back/err"
	"gms-back/project"
	"gms-back/proto/request"
)
{{ range .FuncNames }}
func {{ . }}(p project.IProject, target string, iReq, iMeta interface{}) (res interface{}, ei err.ErrInfo) {
	var (
		routeMeta project.RouteMeta
		body	  *request.{{ . }}Res
		ctx		  = context.TODO()
	)
	req, ok := iReq.(*request.{{ . }}Req)
	if !ok {
		ei.Code = err.ParamErr
		goto TAG
	}
	if routeMeta, ok = iMeta.(project.RouteMeta); !ok {
		ei.Code = err.ParamErr
		goto TAG
	}

	ei.Code = err.Succ
	body = &request.{{ . }}Res{}

TAG:
	res = body
	return
}
{{ end }}
func init() { {{ range .FuncNames }}
	project.APIRegist(api.{{ . }}, constant.GetProjectGroupGames(), project.Info{
		Auth: true,
		Unmarshal: func(body io.ReadCloser) (interface{}, error) {
			return project.APIUnmarshal(body, &request.{{ . }}Req{})
		},
		Execute: {{ . }},
	}){{ end }}
}
`
)
