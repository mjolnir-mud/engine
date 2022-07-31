package entity_registry

import (
	"context"
	"testing"

	"github.com/mjolnir-mud/engine"
	"github.com/stretchr/testify/assert"
)

type testEntityType struct{}

func (t *testEntityType) Name() string {
	return "test"
}

func (t *testEntityType) Create(_ string, args map[string]interface{}) map[string]interface{} {
	return args
}

func setup() {
	engine.Start("test")
	engine.Redis.FlushDB(context.Background())
	Register(&testEntityType{})
	Start()
}

func teardown() {
	Stop()
	engine.Redis.FlushDB(context.Background())
	engine.Stop()
}

func TestAllComponents(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "__type:test:testStringComponent", "string", 0)
	engine.Redis.Set(context.Background(), "test:testStringComponent", "test", 0)

	engine.Redis.Set(context.Background(), "__type:test:testHashComponent", "map", 0)
	engine.Redis.HSet(context.Background(), "test:testHashComponent", "test", "test")

	engine.Redis.Set(context.Background(), "__type:test:testSetComponent", "set", 0)
	engine.Redis.SAdd(context.Background(), "test:testSetComponent", "test")

	components := AllComponents("test")

	assert.Equal(t, map[string]interface{}{
		"testStringComponent": "test",
		"testHashComponent":   map[string]interface{}{"test": "test"},
		"testSetComponent":    []interface{}{"test"},
	}, components)

	teardown()
}

func TestSetIntComponent(t *testing.T) {
	setup()

	SetIntComponent("test", "testIntComponent", 1)

	i, err := engine.Redis.Get(context.Background(), "test:testIntComponent").Int()

	assert.Nil(t, err)

	assert.Equal(t, 1, i)

	teardown()
}

func TestSetBoolComponent(t *testing.T) {
	setup()

	SetBoolComponent("test", "testBoolComponent", true)

	b, err := engine.Redis.Get(context.Background(), "test:testBoolComponent").Bool()

	assert.Nil(t, err)

	assert.Equal(t, true, b)

	SetBoolComponent("test", "testBoolComponent", false)

	b, err = engine.Redis.Get(context.Background(), "test:testBoolComponent").Bool()

	assert.Nil(t, err)

	assert.Equal(t, false, b)

	teardown()
}

func TestSetInt64Component(t *testing.T) {
	setup()

	AddInt64ComponentToEntity("test", "testInt64Component", 1)

	i, err := engine.Redis.Get(context.Background(), "test:testInt64Component").Int64()

	assert.Nil(t, err)

	assert.Equal(t, int64(1), i)

	teardown()
}

func TestSetStringComponent(t *testing.T) {
	setup()

	AddStringComponentToEntity("test", "testStringComponent", "test")

	s := engine.Redis.Get(context.Background(), "test:testStringComponent").Val()

	assert.Equal(t, "test", s)

	teardown()
}

func TestSetIntInHashComponent(t *testing.T) {
	setup()

	SetIntInHashComponent("test", "testHashComponent", "test", 1)

	i, err := engine.Redis.HGet(context.Background(), "test:testHashComponent", "test").Int()

	assert.Nil(t, err)

	assert.Equal(t, 1, i)

	teardown()
}

func TestSetStringInHashComponent(t *testing.T) {
	setup()

	SetStringInHashComponent("test", "testHashComponent", "test", "test")

	s := engine.Redis.HGet(context.Background(), "test:testHashComponent", "test").Val()

	assert.Equal(t, "test", s)

	teardown()
}

func TestAddStringToSetComponent(t *testing.T) {
	setup()

	AddStringToSetComponent("test", "testSetComponent", "test")

	s, err := engine.Redis.SMembers(context.Background(), "test:testSetComponent").Result()

	assert.Nil(t, err)

	assert.Equal(t, []string{"test"}, s)

	teardown()
}

func TestGetBoolComponent(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testBoolComponent", 1, 0)

	b, err := GetBoolComponent("test", "testBoolComponent")

	assert.Nil(t, err)

	assert.Equal(t, true, b)

	teardown()
}

func TestGetIntComponent(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testIntComponent", 1, 0)

	i, err := GetIntComponent("test", "testIntComponent")

	assert.Nil(t, err)

	assert.Equal(t, 1, i)

	teardown()
}

func TestGetInt64Component(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testInt64Component", 1, 0)

	i, err := GetInt64Component("test", "testInt64Component")

	assert.Nil(t, err)

	assert.Equal(t, int64(1), i)

	teardown()
}

func TestGetStringComponent(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testStringComponent", "test", 0)

	s, err := GetStringComponent("test", "testStringComponent")

	assert.Nil(t, err)

	assert.Equal(t, "test", s)

	teardown()
}

func TestGetIntFromHashComponent(t *testing.T) {
	setup()

	engine.Redis.HSet(context.Background(), "test:testHashComponent", "test", 1)

	i, err := GetIntFromHashComponent("test", "testHashComponent", "test")

	assert.Nil(t, err)

	assert.Equal(t, 1, i)

	teardown()
}

func TestGetStringFromHashComponent(t *testing.T) {
	setup()

	engine.Redis.HSet(context.Background(), "test:testHashComponent", "test", "test")

	s, err := GetStringFromHashComponent("test", "testHashComponent", "test")

	assert.Nil(t, err)

	assert.Equal(t, "test", s)

	teardown()
}

func TestGetStringsFromSetComponent(t *testing.T) {
	setup()

	engine.Redis.SAdd(context.Background(), "test:testSetComponent", "test")

	s, err := GetStringsFromSetComponent("test", "testSetComponent")

	assert.Nil(t, err)

	assert.Equal(t, []string{"test"}, s)

	teardown()
}

func TestHasComponent(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testBoolComponent", 1, 0)

	b := HasComponent("test", "testBoolComponent")

	assert.Equal(t, true, b)

	teardown()
}

func TestHasComponentValue(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "test:testBoolComponent", 1, 0)
	engine.Redis.Set(context.Background(), "test:testIntComponent", 1, 0)
	engine.Redis.Set(context.Background(), "test:testInt64Component", 1, 0)
	engine.Redis.Set(context.Background(), "test:testStringComponent", "test", 0)
	engine.Redis.HSet(context.Background(), "test:testHashComponent", "test", 1)
	engine.Redis.SAdd(context.Background(), "test:testSetComponent", "test")

	b := HasComponentValue("test", "testBoolComponent", true)

	assert.Equal(t, true, b)

	b = HasComponentValue("test", "testIntComponent", 1)

	assert.Equal(t, true, b)

	b = HasComponentValue("test", "testInt64Component", int64(1))

	assert.Equal(t, true, b)

	b = HasComponentValue("test", "testStringComponent", "test")

	assert.Equal(t, true, b)

	b = HasComponentValue("test", "testHashComponent", map[string]interface{}{"test": "1"})

	assert.Equal(t, true, b)

	b = HasComponentValue("test", "testSetComponent", []string{"test"})

	assert.Equal(t, true, b)

	teardown()
}

func TestEntityExists(t *testing.T) {
	setup()

	engine.Redis.Set(context.Background(), "__type:test:testBoolComponent", true, 0)

	b := EntityExists("test")

	assert.Equal(t, true, b)

	teardown()
}

func TestAdd(t *testing.T) {
	setup()

	err := AddWithID("test", "testEntity", map[string]interface{}{
		"testString": "string",
		"testInt":    1,
		"testInt64":  int64(1),
		"testBool":   true,
		"testHash":   map[string]interface{}{"test": "1"},
		"testSet":    []string{"test"},
	})
	assert.Nil(t, err)

	s, err := GetStringComponent("testEntity", "testString")
	assert.Nil(t, err)
	assert.Equal(t, "string", s)

	i, err := GetIntComponent("testEntity", "testInt")
	assert.Nil(t, err)
	assert.Equal(t, 1, i)

	i64, err := GetInt64Component("testEntity", "testInt64")
	assert.Nil(t, err)
	assert.Equal(t, int64(1), i64)

	b, err := GetBoolComponent("testEntity", "testBool")
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	h, err := GetHashComponent("testEntity", "testHash")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"test": "1"}, h)

	teardown()
}

func TestSetStringInSetComponent(t *testing.T) {
	setup()

	SetStringInSetComponent("test", "testSetComponent", "test")

	s, err := GetStringsFromSetComponent("test", "testSetComponent")
	assert.Nil(t, err)
	assert.Equal(t, []string{"test"}, s)

	teardown()
}

func TestAllEntitiesByType(t *testing.T) {
	setup()

	_ = AddWithID("test", "testEntityA", map[string]interface{}{})
	_ = AddWithID("test", "testEntityB", map[string]interface{}{})

	entities := AllEntitiesByType("test")
	assert.Equal(t, 2, len(entities))

	teardown()
}

func TestAllEntitiesByTypeWithComponent(t *testing.T) {
	setup()

	_ = AddWithID("test", "testEntityA", map[string]interface{}{
		"testEntityA": "test",
	})
	_ = AddWithID("test", "testEntityB", map[string]interface{}{})

	entities := AllEntitiesByTypeWithComponent("test", "testEntityA")
	assert.Equal(t, 1, len(entities))

	teardown()
}

func TestAllEntitiesByTypeWithComponentValue(t *testing.T) {
	setup()

	_ = AddWithID("test", "testEntityA", map[string]interface{}{
		"testEntityA": "test",
	})
	_ = AddWithID("test", "testEntityB", map[string]interface{}{})

	entities := AllEntitiesByTypeWithComponentValue("test", "testEntityA", "test")
	assert.Equal(t, 1, len(entities))

	teardown()
}
