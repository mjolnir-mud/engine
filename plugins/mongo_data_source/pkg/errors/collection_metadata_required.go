package errors

import (
	"fmt"

	"github.com/mjolnir-mud/engine/plugins/mongo_data_source/pkg/constants"
)

type CollectionMetadataRequiredError struct {
	ID string
}

func (e CollectionMetadataRequiredError) Error() string {
	return fmt.Sprintf("data source %s does not return an entity with the %s set in the metadata", e.ID, constants.MetadataCollectionKey)
}
