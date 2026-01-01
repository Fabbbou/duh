package file_db

// FileHandler allows to handle different file formats without depending on the tomll package directly.
type FileHandler interface {
	Extension() string
	LoadRepositoryFile(path string) (*RepositoryDto, error)
	SaveRepositoryFile(path string, data *RepositoryDto) error
	LoadUserPreferenceFile(path string) (*UserPreferenceDto, error)
	SaveUserPreferenceFile(path string, data *UserPreferenceDto) error
}
