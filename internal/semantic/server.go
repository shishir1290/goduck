package semantic

func (a *Analyzer) checkServer() {

	if a.program.Server == nil {
		a.addError("missing server block")
		return
	}

	if a.program.Server.Port < 1 ||
		a.program.Server.Port > 65535 {

		a.addError("server port must be between 1 and 65535")
	}
}