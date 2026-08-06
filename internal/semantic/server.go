package semantic

func (a *Analyzer) checkServer() error {

	if a.program.Server == nil {
		return errorf("missing server block")
	}

	if a.program.Server.Port <= 0 {
		return errorf("invalid server port")
	}

	return nil
}