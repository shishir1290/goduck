package ir

type Module struct {
	Name string

	Controller string
	Service    string
	Repository string

	Routes []Route
}