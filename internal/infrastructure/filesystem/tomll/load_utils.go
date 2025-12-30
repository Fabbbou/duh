package tomll

func LoadRepository(filePath string) (*RepositoryToml, error) {
	return LoadToml[RepositoryToml](filePath)
}

func LoadUserPreferences(filePath string) (*UserPreferenceToml, error) {
	return LoadToml[UserPreferenceToml](filePath)
}
