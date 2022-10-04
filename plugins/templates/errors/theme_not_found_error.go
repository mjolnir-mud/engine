package errors

type ThemeNotFoundError struct {
	Name string
}

func (e ThemeNotFoundError) Error() string {
	return "theme not found: " + e.Name
}
