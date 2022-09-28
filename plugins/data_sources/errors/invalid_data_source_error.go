package errors

import "fmt"

type InvalidDataSourceError struct {
	Source string
}

func (e InvalidDataSourceError) Error() string {
	return fmt.Sprintf("data source %s does not exist", e.Source)
}
