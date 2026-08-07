package semantic

import "fmt"

func (a *Analyzer) checkRoutes() error {

	for _, module := range a.program.Modules {

		seen := make(map[string]bool)

		for _, route := range module.Routes {

			if route.Path == "" {
				return errorf(
					"module %s: route path cannot be empty",
					module.Name,
				)
			}

			if route.Method == "" {
				return errorf(
					"module %s: route method cannot be empty",
					module.Name,
				)
			}

			if route.Action == "" {
				return errorf(
					"module %s: missing action",
					module.Name,
				)
			}

			key := fmt.Sprintf(
				"%s:%s",
				route.Method,
				route.Path,
			)

			if seen[key] {
				return errorf(
					"module %s: duplicate route %s %s",
					module.Name,
					route.Method,
					route.Path,
				)
			}

			seen[key] = true
		}
	}

	return nil
}