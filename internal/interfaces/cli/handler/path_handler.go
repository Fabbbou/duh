package handler

import (
	"duh/internal/application/usecase"
	"fmt"

	"github.com/spf13/cobra"
)

type PathHandler struct {
	pathUsecase *usecase.PathUsecase
}

func NewPathHandler(pathUsecase *usecase.PathUsecase) *PathHandler {
	return &PathHandler{
		pathUsecase: pathUsecase,
	}
}

func (p *PathHandler) ShowPath(cmd *cobra.Command, args []string) {
	path, err := p.pathUsecase.GetBasePath()
	if err != nil {
		cmd.PrintErrf("Error retrieving repository paths: %v\n", err)
		return
	}
	fmt.Println(path) // Using fmt.Println to match original behavior
}

func (p *PathHandler) ListAllPaths(cmd *cobra.Command, args []string) {
	paths, err := p.pathUsecase.GetAllPaths()
	if err != nil {
		cmd.PrintErrf("Error listing paths: %v\n", err)
		return
	}

	if len(paths) == 0 {
		cmd.Println("No paths configured")
		return
	}

	cmd.Println("Repository paths:")
	for _, path := range paths {
		cmd.Printf("  %s\n", path)
	}
}
