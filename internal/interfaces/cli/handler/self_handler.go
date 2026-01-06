package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/interfaces/cli/std"

	"github.com/spf13/cobra"
)

type SelfHandler struct {
	selfUsecase *usecase.SelfUsecase
}

func NewSelfHandler(selfUsecase *usecase.SelfUsecase) *SelfHandler {
	return &SelfHandler{
		selfUsecase: selfUsecase,
	}
}

func (s *SelfHandler) ShowConfigPath(cmd *cobra.Command, args []string) {
	configPath, err := s.selfUsecase.GetBasePath()
	if err != nil {
		std.Errf("Error retrieving config path: %v\n", err)
		return
	}
	std.Ln(configPath)
}

func (s *SelfHandler) ShowRepositoriesPath(cmd *cobra.Command, args []string) {
	reposPath, err := s.selfUsecase.RepositoriesPath()
	if err != nil {
		std.Errf("Error retrieving repositories path: %v\n", err)
		return
	}
	std.Ln(reposPath)
}

func (s *SelfHandler) GetVersion(cmd *cobra.Command, args []string) {
	std.Ln(s.selfUsecase.GetVersion())
}

func (s *SelfHandler) Update(cmd *cobra.Command, args []string) {
	std.Ln("Checking for updates...")
	err := s.selfUsecase.UpdateSelf()
	if err != nil {
		std.Errf("Update failed: %v\n", err)
		return
	}
	std.Ln("Successfully updated to the latest version!")
	std.Ln("Please restart duh or reload your shell to use the new version.")
}
