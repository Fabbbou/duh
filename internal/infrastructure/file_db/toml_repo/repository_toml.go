package toml_repo

type RepositoryToml struct {
	Aliases  map[string]string `toml:"aliases"`
	Exports  map[string]string `toml:"exports"`
	Metadata MetadataMap       `toml:"metadata"`
}

type MetadataMap struct {
	UrlOrigin  string `toml:"url_origin"`
	NameOrigin string `toml:"name_origin"`
}
