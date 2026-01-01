package tomll

type RepositoryToml struct {
	Aliases  map[string]string `toml:"aliases"`
	Exports  map[string]string `toml:"exports"`
	Metadata MetadataMap       `toml:"metadata"`
}

type MetadataMap struct {
	UrlOrigin  string `toml:"url_origin"`
	NameOrigin string `toml:"name_origin"`
}

type UserPreferenceToml struct {
	Repositories RepositoriesPreference `toml:"repositories"`
}

type RepositoriesPreference struct {
	ActivatedRepositories []string `toml:"activated_repos"`
	DefaultRepositoryName string   `toml:"default_repo_name"`
}
