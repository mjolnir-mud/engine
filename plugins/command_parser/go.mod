module github.com/mjolnir-mud/engine/plugins/command_parser

go 1.17

replace (
	github.com/mjolnir-mud/engine => ../../
	github.com/mjolnir-mud/engine/plugins/controllers => ../controllers
	github.com/mjolnir-mud/engine/plugins/ecs => ../ecs
	github.com/mjolnir-mud/engine/plugins/sessions => ../sessions
	github.com/mjolnir-mud/engine/plugins/templates => ../templates
)

require (
	github.com/alecthomas/kong v0.6.1
	github.com/mjolnir-mud/engine/plugins/sessions v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/charmbracelet/lipgloss v0.5.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-redis/redis/v9 v9.0.0-beta.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mjolnir-mud/engine v0.1.0 // indirect
	github.com/mjolnir-mud/engine/plugins/controllers v0.0.0-00010101000000-000000000000 // indirect
	github.com/mjolnir-mud/engine/plugins/ecs v0.0.0-00010101000000-000000000000 // indirect
	github.com/mjolnir-mud/engine/plugins/templates v0.0.0-00010101000000-000000000000 // indirect
	github.com/muesli/reflow v0.2.1-0.20210115123740-9e1d0d53df68 // indirect
	github.com/muesli/termenv v0.11.1-0.20220204035834-5ac8409525e0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/zerolog v1.27.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
