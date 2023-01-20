package task

type Task interface {
	Name() string
	URL() string
}

type DefaultTask struct {
	name, url string
}

func (d *DefaultTask) Name() string { return d.name }
func (d *DefaultTask) URL() string  { return d.url }

func New(name, url string) Task {
	return &DefaultTask{
		name: name,
		url:  url,
	}
}
