## Next Release

* Added `ecs` plugin
* Added `data_sources` plugin
* Added `yaml_data_source` plugin
* Added `templates` plugin
* Added `command_parser` plugin
* Added `mongo_data_source` plugin

### Engine

#### Breaking Changes
* redis connectivity has been moved out of the world plugin into the engine
* plugin `Init` has been renamed to `Start`
* `engine.AddCommand` has been renamed to `engine.AddCLICommand`
* consolidated plugins into monorepo

### World

#### New Features
* move command can optionally not broadcast a message to members in a room
* character starting locationc an be configured by setting the `MJOLNIR_CHARACTER_STARTING_LOCATION` environment 
  variable

#### Bug Fixes

* miss naming of `room.MoveWithMesssageForSession()` to `room.MoveWithMessageForSession()`
* characters are not created with a starting location

#### Breaking Changes
* `world.ParseCommand()` has been removed and moved to the `command_parser` plugin
* ECS components have been broken out into their own plugin

### Templates

#### New Features
* [template & theming engine plugin]()
  * new templates created for existing commands
