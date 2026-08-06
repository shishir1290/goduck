package semantic

func (a *Analyzer) checkApp() error {

	if a.program.App == nil {
		return errorf("missing app declaration")
	}

	if a.program.App.Name == "" {
		return errorf("app name cannot be empty")
	}

	return nil
}