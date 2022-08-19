package errors

import "fmt"

type MetadataRequiredError struct {
	ID string
}

func (e MetadataRequiredError) Error() string {
	return fmt.Sprintf("data source %s does not return an entity with the type set in the metadata", e.ID)
}
