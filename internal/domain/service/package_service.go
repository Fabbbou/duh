package service

import (
	"duh/internal/domain/entity"
	"duh/internal/domain/errorss"
	"duh/internal/domain/port"
)

type PackageService struct {
	dbPort port.DbPort
}

func NewPackageService(dbPort port.DbPort) *PackageService {
	return &PackageService{
		dbPort: dbPort,
	}
}

// GetPackagesGroupedByStatus returns packages grouped by enabled/disabled status
func (p *PackageService) GetPackagesGroupedByStatus() (map[string][]string, error) {
	packages, err := p.dbPort.GetAllPackages()
	if err != nil {
		return nil, err
	}

	enabledPackages, err := p.dbPort.GetEnabledPackages()
	if err != nil {
		return nil, err
	}

	// Business logic: create a lookup map for enabled packages
	enabledMap := make(map[string]bool, len(enabledPackages))
	for _, pkg := range enabledPackages {
		enabledMap[pkg.Name] = true
	}

	// Business logic: categorize packages
	result := map[string][]string{
		"enabled":  make([]string, 0),
		"disabled": make([]string, 0),
	}

	for _, pkg := range packages {
		if enabledMap[pkg.Name] {
			result["enabled"] = append(result["enabled"], pkg.Name)
		} else {
			result["disabled"] = append(result["disabled"], pkg.Name)
		}
	}

	return result, nil
}

// EnablePackage enables a package and validates business rules
func (p *PackageService) EnablePackage(packageName string) error {
	// Business rule: validate package exists
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}

	return p.dbPort.EnablePackage(packageName)
}

// DisablePackage disables a package and validates business rules
func (p *PackageService) DisablePackage(packageName string) error {
	// Business rule: validate package exists
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}

	// Business rule: cannot disable the last enabled package
	enabled, err := p.dbPort.GetEnabledPackages()
	if err != nil {
		return err
	}

	if len(enabled) <= 1 {
		return &errorss.BusinessRuleError{
			Rule:    "minimum_one_package",
			Message: "cannot disable the last enabled package",
		}
	}

	return p.dbPort.DisablePackage(packageName)
}

// validatePackageExists checks if a package exists
func (p *PackageService) validatePackageExists(packageName string) error {
	packages, err := p.dbPort.GetAllPackages()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.Name == packageName {
			return nil
		}
	}

	return &errorss.NotFoundError{
		Resource: "package",
		ID:       packageName,
	}
}

// Additional methods needed by the use case

func (p *PackageService) DeletePackage(packageName string) error {
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}
	return p.dbPort.DeletePackage(packageName)
}

func (p *PackageService) SetDefaultPackage(packageName string) error {
	if err := p.EnablePackage(packageName); err != nil {
		return err
	}
	return p.dbPort.ChangeDefaultPackage(packageName)
}

func (p *PackageService) GetDefaultPackageName() (string, error) {
	pkg, err := p.dbPort.GetDefaultPackage()
	if err != nil {
		return "", err
	}
	return pkg.Name, nil
}

func (p *PackageService) RenamePackage(oldName, newName string) error {
	if err := p.validatePackageExists(oldName); err != nil {
		return err
	}
	return p.dbPort.RenamePackage(oldName, newName)
}

func (p *PackageService) AddAndEnablePackage(url string, name *string) error {
	pkg, err := p.dbPort.AddPackage(url, name)
	if err != nil {
		return err
	}
	return p.dbPort.EnablePackage(pkg)
}

func (p *PackageService) CreatePackage(name string) error {
	_, err := p.dbPort.CreatePackage(name)
	return err
}

func (p *PackageService) UpdatePackages(strategy string) (entity.PackageUpdateResults, error) {
	results, err := p.dbPort.UpdatePackages(strategy)
	return entity.PackageUpdateResults{
		LocalChangesDetected: results.LocalChangesDetected,
		OtherErrors:          results.OtherErrors,
	}, err
}

func (p *PackageService) EditPackage(packageName string) error {
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}
	return p.dbPort.EditPackage(packageName)
}

func (p *PackageService) PushPackage(packageName string) error {
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}
	return p.dbPort.PushPackage(packageName)
}

func (p *PackageService) EditGitconfig(packageName string) error {
	if err := p.validatePackageExists(packageName); err != nil {
		return err
	}
	return p.dbPort.EditGitconfig(packageName)
}
