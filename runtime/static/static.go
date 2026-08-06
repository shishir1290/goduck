package static

import "net/http"

type Static struct {
	Prefix string
	Root   string
}

func New(prefix, root string) *Static {
	return &Static{
		Prefix: prefix,
		Root:   root,
	}
}

func (s *Static) Handler() http.Handler {

	return http.StripPrefix(
		s.Prefix,
		http.FileServer(http.Dir(s.Root)),
	)
}

