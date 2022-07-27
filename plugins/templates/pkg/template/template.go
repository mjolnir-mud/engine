package template

type Template interface {
	Name() string
	Style() string
	Render(ctx interface{}) (string, error)
}
