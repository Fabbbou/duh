package contexts

import (
	"duh/internal/domain/service"
	"duh/internal/infrastructure/toml_storage"
	"path/filepath"
)

// BuildCliService builds and returns a CliService with all its dependencies properly initialized.
// If a basePath is provided, it uses that as root path for Duh configs
// If not, it defaults to the standard path in the user's home directory : ~/.local/share/duh
func BuildCliService(basePath *string) (service.CliService, error) {
	var pathProvider service.PathProvider
	if basePath != nil {
		pathProvider = service.NewCustomPathProvider(*basePath)
	} else {
		pathProvider = &service.BasePathProvider{}
	}
	directoryService := service.NewDirectoryService(pathProvider)
	_, err := directoryService.CreateRepository("local")
	if err != nil {
		return service.CliService{}, err
	}

	userPreferencePath := filepath.Join(*basePath, "user_preferences.toml")
	userPrefRepo, err := toml_storage.NewTomlDbRepository(userPreferencePath)
	if err != nil {
		return service.CliService{}, err
	}
	userPrefService := service.NewUserPreferenceService(userPrefRepo)
	err = userPrefService.InitUserPreference()
	if err != nil {
		return service.CliService{}, err
	}

	dbRepoFactory := toml_storage.NewTomlDbRepositoryFactory()
	repositoryService := service.NewRepositoriesService(directoryService, dbRepoFactory)

	return service.NewCliService(repositoryService, userPrefService), nil
}
