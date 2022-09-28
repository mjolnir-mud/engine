package errors

import "fmt"

type MetadataRequiredError struct {
	ID string
}

func (e MetadataRequiredError) Error() string {
	return fmt.Sprintf("entity with id %s did not return an entity with the type set in the metadata", e.ID)
}
