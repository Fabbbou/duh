package toml_repo

type RepositoryToml struct {
	Aliases map[string]string `toml:"aliases"`
	Exports map[string]string `toml:"exports"`
}
