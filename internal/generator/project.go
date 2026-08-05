package generator

type Project struct {
	Name  string
	Files []File
}

func (p *Project) AddFile(path string, content string) {
	p.Files = append(p.Files, File{
		Path:    path,
		Content: content,
	})
}