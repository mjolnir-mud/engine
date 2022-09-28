package errors

import "fmt"

type CollectionMismatchError struct {
	SourceCollection string
	TargetCollection string
}

func (e CollectionMismatchError) Error() string {
	return fmt.Sprintf("data source %s does not return an entity with the %s set in the metadata", e.SourceCollection, e.TargetCollection)
}
