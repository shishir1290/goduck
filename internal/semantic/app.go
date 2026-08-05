package semantic

func (a *Analyzer) checkApp() {

	if a.program.App == nil {
		a.addError("missing app declaration")
		return
	}

	if a.program.App.Name == "" {
		a.addError("app name cannot be empty")
	}
}