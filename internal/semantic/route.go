package semantic

import "fmt"

func (a *Analyzer) checkRoutes() error {

	seen := make(map[string]bool)

	for _, route := range a.program.Routes {

		if route.Path == "" {
			return errorf("route path cannot be empty")
		}

		if route.Controller == "" {
			return errorf("missing controller")
		}

		if route.Action == "" {
			return errorf("missing action")
		}

		key := fmt.Sprintf(
			"%s:%s",
			route.Method,
			route.Path,
		)

		if seen[key] {
			return errorf(
				"duplicate route %s %s",
				route.Method,
				route.Path,
			)
		}

		seen[key] = true
	}

	return nil
}