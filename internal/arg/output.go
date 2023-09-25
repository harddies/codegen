package arg

func SetOutput(output string) Options {
	return func(s *Sets) {
		if output == "" {
			return
		}
		s.Output = output
	}
}
