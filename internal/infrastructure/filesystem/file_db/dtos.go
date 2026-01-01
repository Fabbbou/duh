package file_db

type Key = string
type Value = string

type AliasesMap = map[Key]Value
type ExportsMap = map[Key]Value

type RepositoryDto struct {
	Aliases  AliasesMap
	Exports  ExportsMap
	Metadata MetadataDto
}

type MetadataDto struct {
	UrlOrigin  string
	NameOrigin string
}

type UserPreferenceDto struct {
	Repositories RepositoriesPreferenceDto
}

type RepositoriesPreferenceDto struct {
	ActivatedRepositories []string
	DefaultRepositoryName string
}
