package session

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionType_Name(t *testing.T) {
	assert.Equal(t, "session", Type.Name())
}

func TestSessionType_Create(t *testing.T) {
	assert.Equal(t, map[string]interface{}{
		"expireIn": 900,
		"store": map[string]interface{}{
			"controller": "login",
		},
		"commandSets": []interface{}{
			"base",
			"movement",
		},
	}, Type.Create("", map[string]interface{}{}))
}
