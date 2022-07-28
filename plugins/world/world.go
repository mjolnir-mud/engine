package world

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/command_parser"
	"github.com/mjolnir-mud/engine/plugins/command_parser/pkg/command_set"
	"github.com/mjolnir-mud/engine/plugins/templates"
	"github.com/mjolnir-mud/engine/plugins/world/internal/command_sets"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controllers/create_character"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controllers/game"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controllers/login"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controllers/new_acccount"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controllers/select_controller"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_types/area"
	character2 "github.com/mjolnir-mud/engine/plugins/world/internal/entity_types/character"
	room2 "github.com/mjolnir-mud/engine/plugins/world/internal/entity_types/room"
	session2 "github.com/mjolnir-mud/engine/plugins/world/internal/entity_types/session"
	"github.com/mjolnir-mud/engine/plugins/world/internal/entity_types/zone"
	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"github.com/mjolnir-mud/engine/plugins/world/internal/system_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/character"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/location"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/persisted"
	"github.com/mjolnir-mud/engine/plugins/world/internal/systems/room"
	templates2 "github.com/mjolnir-mud/engine/plugins/world/internal/templates"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/db"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/entity_type"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/system"
	_default "github.com/mjolnir-mud/engine/plugins/world/pkg/themes/default"
	"github.com/spf13/cobra"
)

type world struct {
	stop chan bool
}

var log = logger.Logger

// Init initializes the world
func (w world) Start() error {
	log.Info().Msg("initializing")
	engine.AddCLICommand(&cobra.Command{
		Use:   "world",
		Short: "start the Mjolnir World",
		Long:  "start the Mjolnir World",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("starting server")
			db.Start()
			session.ControllerRegistry.Start()

			// Register systems
			RegisterSystem(location.System)
			RegisterSystem(persisted.System)
			RegisterSystem(character.System)
			RegisterSystem(room.System)

			// Register controllers
			RegisterController(login.Controller)
			RegisterController(new_acccount.Controller)
			RegisterController(create_character.Controller)
			RegisterController(select_controller.Controller)
			RegisterController(game.Controller)

			// Register types
			RegisterEntityType(session2.Type)
			RegisterEntityType(character2.Type)
			RegisterEntityType(room2.Type)
			RegisterEntityType(zone.Type)
			RegisterEntityType(area.Type)

			command_parser.RegisterCommandSet(command_sets.Base)
			command_parser.RegisterCommandSet(command_sets.Movement)

			entity_registry.Start()
			system_registry.Start()

			templates.RegisterTheme(_default.Theme)

			templates.RegisterTemplate(templates2.Say)
			templates.RegisterTemplate(templates2.Walking)

			session.Registry.Start()

			<-w.stop
		},
	})

	return nil
}

// Name returns the name of the plugin
func (w world) Name() string {
	return "world"
}

// GetStringComponent returns the value of a string component.
func GetStringComponent(id string, name string) (string, error) {
	return entity_registry.GetStringComponent(id, name)
}

// GetInt64Component returns the value of an int64 component.
func GetInt64Component(id string, name string) (int64, error) {
	return entity_registry.GetInt64Component(id, name)
}

// GetIntComponent returns the value of an int component.
func GetIntComponent(id string, name string) (int, error) {
	return entity_registry.GetIntComponent(id, name)
}

// GetIntFromHashComponent  returns the value of an int value from a hash component.
func GetIntFromHashComponent(id string, name string, key string) (int, error) {
	return entity_registry.GetIntFromHashComponent(id, name, key)
}

// GetStringFromHashComponent returns the value of a string value from a hash component.
func GetStringFromHashComponent(id string, name string, key string) (string, error) {
	return entity_registry.GetStringFromHashComponent(id, name, key)
}

// RegisterCommandSet registers a command set with the world.
func RegisterCommandSet(set *command_set.CommandSet) {
	//command_registry.RegisterCommandSet(set)
}

// RegisterController registers a controller with the world.
func RegisterController(controller session.Controller) {
	session.ControllerRegistry.Register(controller)
}

// RegisterEntityType registers an entity type with the world.
func RegisterEntityType(entityType entity_type.EntityType) {
	entity_registry.Register(entityType)
}

// RegisterSystem registers a system with the world.
func RegisterSystem(system system.System) {
	system_registry.Register(system)
}

// RegisterLoadableDirectory registers a directory to be loaded.
func RegisterLoadableDirectory(dir string) {
	entity_registry.RegisterLoadableDirectory(dir)
}

func SetStringComponent(id string, name, value string) {
	entity_registry.SetStringComponent(id, name, value)
}

func SetIntComponent(id string, name string, value int) {
	entity_registry.SetIntComponent(id, name, value)
}

func SetInt64Component(id string, name string, value int64) {
	entity_registry.SetInt64Component(id, name, value)
}

func SetStringInHashComponent(id string, name, hash, value string) {
	entity_registry.SetStringInHashComponent(id, name, hash, value)
}

func SetIntInHashComponent(id string, name, key string, value int) {
	entity_registry.SetIntInHashComponent(id, name, key, value)
}

func GetBoolComponent(id string, name string) (bool, error) {
	return entity_registry.GetBoolComponent(id, name)
}

func SetBoolComponent(id string, name string, value bool) {
	entity_registry.SetBoolComponent(id, name, value)
}

func AddEntity(name string, id string, args map[string]interface{}) error {
	return entity_registry.Add(name, id, args)
}

func RemoveEntity(id string) {
	entity_registry.Remove(id)
}

func WriteToConnection(id string, data string) {
	session.Registry.WriteToConnection(id, data)
}

func WriteToConnectionF(id string, format string, args ...interface{}) {
	session.Registry.WriteToConnectionF(id, format, args...)
}

var Plugin = world{
	stop: make(chan bool),
}
