package arg

func SetTarget(target string) Options {
	return func(s *Sets) {
		if target == "" {
			return
		}
		s.Target = target
	}
}
