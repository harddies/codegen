package arg

func SetModule(mod string) Options {
	return func(s *Sets) {
		if mod == "" {
			return
		}
		s.Module = mod
	}
}
