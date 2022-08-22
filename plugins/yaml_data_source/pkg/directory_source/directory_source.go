package directory_source

import (
	"fmt"
	"github.com/mjolnir-mud/engine/plugins/yaml_data_source/internal/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/constants"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/errors"
	constants2 "github.com/mjolnir-mud/engine/plugins/yaml_data_source/pkg/constants"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type MetadataFileRequiredError struct {
	ID string
}

func (e MetadataFileRequiredError) Error() string {
	return fmt.Sprintf("metadata file required for entity %s", e.ID)
}

type DirectorySource struct {
	name   string
	path   string
	logger zerolog.Logger
}

var log = logger.Instance

func New(name string, path string) DirectorySource {
	path, err := filepath.Abs(path)

	if err != nil {
		log.Error().Err(err).Msgf("failed to get absolute path for %s", path)
		panic(err)
	}

	return DirectorySource{
		name: name,
		path: path,
		logger: log.With().
			Str("directory", path).
			Str("name", name).
			Logger(),
	}
}

func (d DirectorySource) Name() string {
	return d.name
}

func (d DirectorySource) Start() error {
	return nil
}

func (d DirectorySource) Stop() error {
	return nil
}

func (d DirectorySource) All() (map[string]map[string]interface{}, error) {
	entities, err := d.loadDirectory(d.path)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to load directory")
		return nil, err
	}

	return entities, nil
}

func (d DirectorySource) Find(search map[string]interface{}) (map[string]map[string]interface{}, error) {
	entities, err := d.loadDirectory(d.path)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to load directory")
		return nil, err
	}

	searchResults := make(map[string]map[string]interface{})

	// if the search map is empty, return all entities
	if len(search) == 0 {
		return entities, nil
	}

	for id, entity := range entities {
		for k, v := range search {
			if entity[k] != v {
				continue
			}
			searchResults[id] = entity
		}
	}

	return searchResults, nil
}

func (d DirectorySource) FindOne(search map[string]interface{}) (string, map[string]interface{}, error) {
	entities, err := d.Find(search)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to find")
		return "", nil, err
	}

	for id, entity := range entities {
		return id, entity, nil
	}

	return "", nil, errors.EntityNotFoundError{ID: ""}
}

func (d DirectorySource) Count(search map[string]interface{}) (int64, error) {
	recs, err := d.Find(search)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to find")
		return 0, err
	}

	return int64(len(recs)), nil
}

func (d DirectorySource) Save(entityId string, entity map[string]interface{}) error {
	metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

	if !ok {
		d.logger.Error().Msg("failed to find metadata")
		return errors.MetadataRequiredError{ID: entityId}
	}

	file, ok := metadata[constants2.MetadataFileKey].(string)

	if !ok {
		return MetadataFileRequiredError{ID: entityId}
	}

	// load from a yml file and process
	d.logger.Debug().Str("file", file).Msg("saving entity to file")
	entities, err := d.loadFromFile(file)

	if err != nil {
		d.logger.Error().Err(err).Msgf("failed to load file %s", file)
		return err
	}

	entities[entityId] = entity

	content, err := yaml.Marshal(entities)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to marshal yaml")
		return err
	}

	err = ioutil.WriteFile(file, content, 0644)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to write file")
		return err
	}

	return nil
}

func (d DirectorySource) loadDirectory(dir string) (map[string]map[string]interface{}, error) {
	// load all yml files in the directory
	files, err := os.Open(dir)

	if err != nil {
		d.logger.Error().Err(err).Msgf("failed to open directory %s", dir)
		return nil, err
	}

	defer func() {
		_ = files.Close()
	}()

	names, err := files.Readdirnames(0)

	if err != nil {
		d.logger.Error().Err(err).Msgf("failed to read directory %s", dir)
		return nil, err
	}

	entities := make(map[string]map[string]interface{})

	for _, name := range names {
		if strings.HasSuffix(name, ".yml") {
			ents, err := d.loadFromFile(fmt.Sprintf("%s/%s", dir, name))

			for k, v := range ents {
				entities[k] = v
			}

			if err != nil {
				d.logger.Error().Err(err).Msgf("failed to load file %s", name)
				return nil, err
			}
		}
	}

	return entities, nil
}

func (d DirectorySource) loadFromFile(file string) (map[string]map[string]interface{}, error) {

	// load from a yml file and process
	d.logger.Debug().Str("file", file).Msg("loading entities from file")
	content, err := os.ReadFile(file)

	if err != nil {
		d.logger.Error().Str("file", file).Err(err).Msgf("failed to read file %s", file)
		return nil, err
	}

	ymlEntries := make(map[string]map[string]interface{})

	err = yaml.Unmarshal(content, ymlEntries)

	if err != nil {
		d.logger.Error().Err(err).Msg("failed to parse yaml")
		return nil, err
	}

	// add the __file metadata to the entity
	for id, entity := range ymlEntries {
		metadata, ok := entity[constants.MetadataKey].(map[string]interface{})

		if !ok {
			d.logger.Error().Str("file", file).Msg("failed to find metadata")
			return nil, errors.MetadataRequiredError{ID: id}
		}

		metadata[constants2.MetadataFileKey] = file
		entity[constants.MetadataKey] = metadata

		ymlEntries[id] = entity
	}

	return ymlEntries, nil
}
