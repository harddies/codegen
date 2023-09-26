package arg

func SetFuncNames(fns []string) Options {
	return func(s *Sets) {
		if len(fns) <= 0 {
			return
		}
		s.FuncNames = fns
	}
}
