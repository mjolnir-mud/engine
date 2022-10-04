package errors

type TemplateNotFoundError struct {
	Name string
}

func (e TemplateNotFoundError) Error() string {
	return "template not found: " + e.Name
}
