## Next Release
### Engine

#### Breaking Changes
* redis connectivity has been moved out of the world plugin into the engine
* plugin `Init` has been renamed to `Start`
* `engine.AddCommand` has been renamed to `engine.AddCLICommand`
* consolidated plugins into monorepo

### World

#### New Features
* move command can optionally not broadcast a message to members in a room

#### Bug Fixes

* miss naming of `room.MoveWithMesssageForSession()` to `room.MoveWithMessageForSession()`

#### Breaking Changes
* `world.ParseCommand()` has been removed and moved to the `command_parser` plugin

### Templates

#### New Features
* [template & theming engine plugin]()
  * new templates created for existing commands
