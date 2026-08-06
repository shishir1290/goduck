package generator

type Project struct {
	Name string

	Files []File

	Directories []string
}

func (p *Project) AddFile(
	path,
	content string,
) {

	p.Files = append(
		p.Files,
		File{
			Path:    path,
			Content: content,
		},
	)
}

func (p *Project) AddDirectory(path string) {

	p.Directories = append(
		p.Directories,
		path,
	)
}