package arg

import "strings"

func SetFuncNames(fns []string) Options {
	return func(s *Sets) {
		if len(fns) <= 0 {
			return
		}

		for i := range fns {
			fns[i] = strings.ToTitle(fns[i])
		}
		s.FuncNames = fns
	}
}
