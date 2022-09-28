package directory_source

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectorySource_All(t *testing.T) {
	ds := New("entities_1", "../../testing/fixtures")

	p, err := filepath.Abs("../../testing/fixtures")

	assert.Nil(t, err)

	e, err := ds.All()

	assert.Nil(t, err)

	assert.Equal(t, map[string]map[string]interface{}{
		"entity_1": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_1.yml"),
			},
			"testComponent": "testing",
		},
		"entity_2": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_1.yml"),
			},
			"testComponent": "test2",
		},
		"entity_3": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_2.yml"),
			},
			"testComponent": "testing",
		},
		"entity_4": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_2.yml"),
			},
			"testComponent": "test2",
		},
	}, e)
}

func TestDirectorySource_Find(t *testing.T) {
	ds := New("entities_1", "../../testing/fixtures")

	p, err := filepath.Abs("../../testing/fixtures")

	assert.Nil(t, err)

	e, err := ds.Find(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)

	assert.Equal(t, map[string]map[string]interface{}{
		"entity_1": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_1.yml"),
			},
			"testComponent": "testing",
		},
		"entity_3": {
			"__metadata": map[string]interface{}{
				"entityType": "fake",
				"file":       path.Join(p, "entities_2.yml"),
			},
			"testComponent": "testing",
		},
	}, e)
}

func TestDirectorySource_FindOne(t *testing.T) {
	ds := New("entities_1", "../../testing/fixtures")

	id, _, err := ds.FindOne(map[string]interface{}{
		"testComponent": "testing",
	})

	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestDirectorySource_Save(t *testing.T) {
	ds := New("entities_1", "../../testing/fixtures")

	p, err := filepath.Abs("../../testing/fixtures")

	defer func() {
		resetToEmptyFile(path.Join(p, "entities_3.yml"))
	}()

	assert.Nil(t, err)

	err = ds.Save("entity_5", map[string]interface{}{
		"__metadata": map[string]interface{}{
			"entityType": "fake",
			"file":       path.Join(p, "entities_3.yml"),
		},
		"testComponent": "testing",
	})

	assert.Nil(t, err)
}

func TestDirectorySource_Count(t *testing.T) {
	ds := New("entities_1", "../../testing/fixtures")

	count, err := ds.Count(map[string]interface{}{})

	assert.Nil(t, err)

	assert.Equal(t, int64(4), count)
}

func resetToEmptyFile(file string) {
	f, err := os.Create(file)

	if err != nil {
		panic(err)
	}

	_ = f.Close()
}
