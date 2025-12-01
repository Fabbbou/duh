package toml_storage

type RepositoryDb struct {
	Aliases map[string]string `toml:"aliases"`
	Exports map[string]string `toml:"exports"`
}

type UserPreferenceDb struct {
	Repositories map[string]string `toml:"repositories"`
}
