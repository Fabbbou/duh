package handler

import (
	"duh/internal/application/usecase"
	"duh/internal/domain/entity"

	"github.com/spf13/cobra"
)

type PackageHandler struct {
	packageUsecase *usecase.PackageUsecase
}

func NewPackageHandler(packageUsecase *usecase.PackageUsecase) *PackageHandler {
	return &PackageHandler{
		packageUsecase: packageUsecase,
	}
}

func (p *PackageHandler) ListPackages(cmd *cobra.Command, args []string) {
	packages, err := p.packageUsecase.ListPackages()
	if err != nil {
		cmd.PrintErrf("Error listing packages: %v\n", err)
		return
	}

	if len(packages["enabled"]) > 0 {
		cmd.Println("Enabled packages:")
		for _, name := range packages["enabled"] {
			cmd.Printf("  ✓ %s\n", name)
		}
	}

	if len(packages["disabled"]) > 0 {
		cmd.Println("Disabled packages:")
		for _, name := range packages["disabled"] {
			cmd.Printf("  ✗ %s\n", name)
		}
	}

	if len(packages["enabled"]) == 0 && len(packages["disabled"]) == 0 {
		cmd.Println("No packages found")
	}
}

func (p *PackageHandler) EnablePackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.EnablePackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error enabling package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' enabled\n", packageName)
}

func (p *PackageHandler) DisablePackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.DisablePackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error disabling package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' disabled\n", packageName)
}

func (p *PackageHandler) DeletePackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.DeletePackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error deleting package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' deleted\n", packageName)
}

func (p *PackageHandler) SetDefaultPackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.SetDefaultPackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error setting default package: %v\n", err)
		return
	}
	cmd.Printf("Default package set to '%s'\n", packageName)
}

func (p *PackageHandler) GetDefaultPackage(cmd *cobra.Command, args []string) {
	packageName, err := p.packageUsecase.GetDefaultPackage()
	if err != nil {
		cmd.PrintErrf("Error getting default package: %v\n", err)
		return
	}
	cmd.Printf("Default package: %s\n", packageName)
}

func (p *PackageHandler) ShowDefaultPackage(cmd *cobra.Command, args []string) {
	// Show current default and help when no subcommand provided
	currentDefault, err := p.packageUsecase.GetDefaultPackage()
	if err != nil {
		cmd.PrintErrf("Error getting current default package: %v\n", err)
		return
	}
	cmd.Printf("Current default package: %s\n\n", currentDefault)
	cmd.Println("Available commands:")
	cmd.Println("  duh package default set <name>  Set package as default")
}

func (p *PackageHandler) RenamePackage(cmd *cobra.Command, args []string) {
	oldName := args[0]
	newName := args[1]
	err := p.packageUsecase.RenamePackage(oldName, newName)
	if err != nil {
		cmd.PrintErrf("Error renaming package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' renamed to '%s'\n", oldName, newName)
}

func (p *PackageHandler) AddPackage(cmd *cobra.Command, args []string) {
	url := args[0]
	var name *string
	if len(args) == 2 {
		name = &args[1]
	}
	err := p.packageUsecase.AddPackage(url, name)
	if err != nil {
		cmd.PrintErrf("Error adding package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' added and enabled\n", url)
}

func (p *PackageHandler) CreatePackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.CreatePackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error creating package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' created and enabled\n", packageName)
}

func (p *PackageHandler) UpdatePackages(cmd *cobra.Command, args []string) {
	// Determine strategy based on flags
	forceFlag, _ := cmd.Flags().GetBool("force")
	commitFlag, _ := cmd.Flags().GetBool("commit")

	strategy := entity.UpdateSafe // default strategy
	if forceFlag && commitFlag {
		cmd.PrintErrf("Cannot use both --force and --commit flags together\n")
		return
	} else if forceFlag {
		strategy = entity.UpdateForce
	} else if commitFlag {
		strategy = entity.UpdateKeep
	}

	results, err := p.packageUsecase.UpdatePackages(strategy)
	if err != nil {
		cmd.PrintErrf("Error updating packages: %v\n", err)
		return
	}

	// Report results
	if len(results.LocalChangesDetected) > 0 {
		cmd.Println("⚠️  Packages with local changes (not updated):")
		for _, pkg := range results.LocalChangesDetected {
			cmd.Printf("  • %s\n", pkg)
		}
		cmd.Println("\nUse --commit to auto-commit changes or --force to discard them.")
	}

	if len(results.OtherErrors) > 0 {
		cmd.Println("❌ Packages with errors:")
		for _, err := range results.OtherErrors {
			cmd.Printf("  • %v\n", err)
		}
	}

	totalPackages := len(results.LocalChangesDetected) + len(results.OtherErrors)
	if totalPackages == 0 {
		cmd.Println("✅ All packages updated successfully")
	} else {
		cmd.Printf("\n%d packages had issues during update\n", totalPackages)
	}
}

func (p *PackageHandler) EditPackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.EditPackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error editing package: %v\n", err)
		return
	}
}

func (p *PackageHandler) PushPackage(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.PushPackage(packageName)
	if err != nil {
		cmd.PrintErrf("Error pushing package: %v\n", err)
		return
	}
	cmd.Printf("Package '%s' pushed successfully\n", packageName)
}

func (p *PackageHandler) EditGitconfig(cmd *cobra.Command, args []string) {
	packageName := args[0]
	err := p.packageUsecase.EditGitconfig(packageName)
	if err != nil {
		cmd.PrintErrf("Error editing gitconfig: %v\n", err)
		return
	}
}
