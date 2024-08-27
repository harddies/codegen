package module

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/harddies/codegen/internal/arg"
	"github.com/harddies/codegen/internal/model"
	"github.com/spf13/cobra"
)

type bts struct {
	arg.Sets
}

func (d *bts) Name() string {
	return model.FlagModuleBts
}

func (d *bts) Run(cmd *cobra.Command, args []string) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	dir, err := os.Getwd()
	if err != nil {
		return
	}
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		if !strings.HasSuffix(dirEntry.Name(), ".go") || strings.HasSuffix(dirEntry.Name(), ".bts.go") {
			continue
		}

		var info *FileBtsInfo
		if info, err = d.readFile(dirEntry.Name()); err != nil {
			return
		}

		if err = d.generateBtsFile(info, dir, dirEntry.Name()); err != nil {
			return
		}
	}
}

func (d *bts) generateBtsFile(info *FileBtsInfo, dir, orgFilename string) (err error) {
	t, err := template.New(btsTplName).Parse(btsTpl)
	if err != nil {
		return
	}

	orgFilenamePrefix := strings.Split(orgFilename, ".")[0]
	outputFilename := filepath.Join(dir, orgFilenamePrefix+".bts.go")
	f, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	if err = t.Execute(f, info); err != nil {
		return
	}

	fmt.Printf("generate file [%s] successfully\n", outputFilename)
	return
}

func (d *bts) readFile(filename string) (info *FileBtsInfo, err error) {
	var f *os.File
	if f, err = os.OpenFile(filename, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}
	defer f.Close()

	var (
		imports       []string
		isImport      bool
		btsAnnotation string
		isBts         bool
	)
	info = &FileBtsInfo{}
	buf := bufio.NewReader(f)
	for {
		var line []byte
		if line, _, err = buf.ReadLine(); err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				break
			}
			return
		}

		lineS := string(line)

		if strings.HasPrefix(lineS, "package ") {
			info.Package = lineS
			continue
		}

		if strings.HasPrefix(lineS, "import ") {
			isImport = true
			continue
		}
		if isImport {
			if strings.HasPrefix(lineS, ")") {
				isImport = false
			} else {
				imports = append(imports, lineS)
			}
			continue
		}

		if strings.Contains(lineS, "bts:") {
			isBts = true
			btsAnnotation = lineS
			continue
		}
		if isBts {
			funcInfo := new(FileBtsFuncInfo)
			argStrs := strings.Split(btsAnnotation, " ")
			for _, argStr := range argStrs {
				if !strings.HasPrefix(argStr, "-") {
					continue
				}
				btsArg := strings.SplitN(argStr, "=", 2)
				if len(btsArg) != 2 {
					continue
				}
				argName, argValue := strings.TrimLeft(btsArg[0], "-"), btsArg[1]
				switch argName {
				case "null_cache":
					funcInfo.NullCache = argValue
				case "struct_name":
					funcInfo.StructName = argValue
				}
			}

			funcInfo.FuncName = funcNameRegexp.FindAllStringSubmatch(lineS, -1)[0][1]
			funcArgStr := funcArgRegexp.FindAllStringSubmatch(lineS, -1)[0][1]
			funcArgs := strings.Split(funcArgStr, ",")
			var variables []string
			for _, funcArg := range funcArgs {
				funcArg = strings.TrimSpace(funcArg)
				funcArgInfo := strings.Split(funcArg, " ")
				variables = append(variables, funcArgInfo[0])
			}
			funcInfo.Variable = strings.Join(variables, ", ")

			funcInfo.FuncDef = strings.TrimLeft(lineS, "\t")

			returnStr := funcReturnRegexp.FindAllStringSubmatch(lineS, -1)[0][1]
			returnUnits := strings.Split(returnStr, ",")
			for _, returnUnit := range returnUnits {
				returnUnit = strings.TrimSpace(returnUnit)
				returnInfo := strings.Split(returnUnit, " ")
				if len(returnInfo) == 1 {
					if returnInfo[0] == "error" {
						continue
					}
					funcInfo.ReturnResVariable = "res"
					funcInfo.ReturnResType = returnInfo[0]
					continue
				}
				if returnInfo[1] == "error" {
					continue
				}
				funcInfo.ReturnResVariable = returnInfo[0]
				funcInfo.ReturnResType = returnInfo[1]
			}

			fmt.Printf("%+v\n", funcInfo)
			info.FuncInfos = append(info.FuncInfos, funcInfo)

			isBts = false
			btsAnnotation = ""
			continue
		}
	}

	if len(imports) > 0 {
		info.Import = strings.Join(imports, "\n")
	}
	return
}

type FileBtsInfo struct {
	Package   string
	Import    string
	FuncInfos []*FileBtsFuncInfo
}

type FileBtsFuncInfo struct {
	NullCache         string
	StructName        string
	FuncName          string
	Variable          string
	FuncDef           string
	ReturnResVariable string
	ReturnResType     string
}

var (
	funcNameRegexp   = regexp.MustCompile(`^[\t\s]*(\w+\d*)`)
	funcArgRegexp    = regexp.MustCompile(`^[\t\s]*[A-Za-z0-9]+\((.*)\) \(`)
	funcReturnRegexp = regexp.MustCompile(` \((.*)\)`)
)

const (
	btsTplName = "bts"

	btsTpl = `{{ .Package }}

import (
{{ .Import }}
)

{{ range .FuncInfos }}
func (r *{{ .StructName }}) {{ .FuncDef }} {
	addCache := true
	{{ .ReturnResVariable }}, err = r.Cache{{ .FuncName }}({{ .Variable }})
	if err != nil {
		addCache = false
		err = nil
	}
	defer func() {
		if {{ .ReturnResVariable }} == {{ .NullCache }} {
			{{ .ReturnResVariable }} = ""
		}
	}()
	if {{ .ReturnResVariable }} != "" {
		return
	}
	var rr interface{}
	sf := r.cacheSF{{ .FuncName }}({{ .Variable }})
	rr, err, _ = cacheSingleFlight.Do(sf, func() (ri interface{}, e error) {
		ri, e = r.Raw{{ .FuncName }}({{ .Variable }})
		return
	})
	{{ .ReturnResVariable }} = rr.({{ .ReturnResType }})
	if err != nil {
		return
	}
	miss := {{ .ReturnResVariable }}
	if miss == "" {
		miss = {{ .NullCache }}
	}
	if !addCache {
		return
	}
	r.AddCache{{ .FuncName }}({{ .Variable }}, miss)
	return
}

{{ end }}`
)
