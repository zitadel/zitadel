package database

var Config = make(map[string]any)

func AddDatabaseConfig(name string, configure func(map[string]any) error) {
	Config[name] = configure
}
