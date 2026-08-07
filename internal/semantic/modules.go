package semantic

func (a *Analyzer) checkModules() error {

	seen := make(map[string]bool)

	for _, module := range a.program.Modules {

		if module.Name == "" {
			return errorf(
				"module name cannot be empty",
			)
		}

		if seen[module.Name] {
			return errorf(
				"duplicate module %s",
				module.Name,
			)
		}

		seen[module.Name] = true

		if module.Controller == "" {
			return errorf(
				"module %s: controller is required",
				module.Name,
			)
		}
	}

	return nil
}