package usecase

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/service"
)

type PackageUsecase struct {
	packageService *service.PackageService
}

func NewPackageUsecase(packageService *service.PackageService) *PackageUsecase {
	return &PackageUsecase{
		packageService: packageService,
	}
}

func (p *PackageUsecase) ListPackages() (map[string][]string, error) {
	// Delegate to domain service for business logic
	return p.packageService.GetPackagesGroupedByStatus()
}

func (p *PackageUsecase) EnablePackage(packageName string) error {
	// Delegate to domain service for business logic
	return p.packageService.EnablePackage(packageName)
}

func (p *PackageUsecase) DisablePackage(packageName string) error {
	// Delegate to domain service for business logic
	return p.packageService.DisablePackage(packageName)
}

func (p *PackageUsecase) DeletePackage(packageName string) error {
	// Application layer: orchestrate disable then delete
	if err := p.packageService.DisablePackage(packageName); err != nil {
		// If disable fails due to business rules, still allow deletion
	}
	return p.packageService.DeletePackage(packageName)
}

func (p *PackageUsecase) SetDefaultPackage(packageName string) error {
	// Application layer: orchestrate enable then set default
	return p.packageService.SetDefaultPackage(packageName)
}

func (p *PackageUsecase) GetDefaultPackage() (string, error) {
	// Simple delegation
	return p.packageService.GetDefaultPackageName()
}

func (p *PackageUsecase) RenamePackage(oldName, newName string) error {
	// Delegate to domain service
	return p.packageService.RenamePackage(oldName, newName)
}

func (p *PackageUsecase) AddPackage(url string, name *string) error {
	// Application layer: orchestrate add then enable
	return p.packageService.AddAndEnablePackage(url, name)
}

func (p *PackageUsecase) CreatePackage(name string) error {
	// Delegate to domain service
	return p.packageService.CreatePackage(name)
}

func (p *PackageUsecase) UpdatePackages(strategy string) (entity.PackageUpdateResults, error) {
	// Delegate to domain service
	return p.packageService.UpdatePackages(strategy)
}

func (p *PackageUsecase) EditPackage(packageName string) error {
	// Delegate to domain service
	return p.packageService.EditPackage(packageName)
}

func (p *PackageUsecase) PushPackage(packageName string) error {
	// Delegate to domain service
	return p.packageService.PushPackage(packageName)
}

func (p *PackageUsecase) EditGitconfig(packageName string) error {
	// Delegate to domain service
	return p.packageService.EditGitconfig(packageName)
}
